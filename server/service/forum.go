package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"server/utils/upload"
)

// ForumDetailData 帖子详情返回结构（帖子 + 评论树）
type ForumDetailData struct {
	database.ForumPost
	Comments []database.ForumComment `json:"comments"`
}

type ForumService struct {
}

var (
	forumAllowedCategory = map[string]struct{}{"技术": {}, "生活": {}}
)

// ForumTags 返回固定标签库全部标签
func (forumService *ForumService) ForumTags() ([]database.BlogTag, error) {
	var tags []database.BlogTag
	if err := global.DB.Order("`group` ASC, number DESC, tag ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// checkTagsExist 校验标签都存在于固定标签库
func (forumService *ForumService) checkTagsExist(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	var count int64
	if err := global.DB.Model(&database.BlogTag{}).Where("tag IN ?", tags).Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(tags) {
		return errors.New("标签不存在")
	}
	return nil
}

// incrTagCount 标签引用数 +1
func (forumService *ForumService) incrTagCount(tx *gorm.DB, tags []string) error {
	for _, tag := range tags {
		if err := tx.Model(&database.BlogTag{}).Where("tag = ?", tag).
			Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
			return err
		}
	}
	return nil
}

// decrTagCount 标签引用数 -1
func (forumService *ForumService) decrTagCount(tx *gorm.DB, tags []string) error {
	for _, tag := range tags {
		if err := tx.Model(&database.BlogTag{}).Where("tag = ?", tag).
			Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
	}
	return nil
}

// ForumPublish 发布帖子，返回新帖子 ID
func (forumService *ForumService) ForumPublish(req request.ForumPublish) (uint, error) {
	// 校验分类
	if _, ok := forumAllowedCategory[req.Category]; !ok {
		return 0, errors.New("分类只能是 技术 或 生活")
	}
	// 校验标签
	if err := forumService.checkTagsExist(req.Tags); err != nil {
		return 0, err
	}

	var postID uint
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		post := database.ForumPost{
			UserID:   req.UserID,
			Title:    req.Title,
			Content:  req.Content,
			Category: req.Category,
			Tags:     database.JSONStringArray(req.Tags),
			Images:   database.JSONStringArray(req.Images),
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		postID = post.ID
		// 更新标签引用数
		return forumService.incrTagCount(tx, req.Tags)
	})
	return postID, err
}

// ForumList 帖子分页列表
func (forumService *ForumService) ForumList(req request.ForumList) (interface{}, int64, error) {
	db := global.DB.Model(&database.ForumPost{})

	if req.Category != nil && *req.Category != "" {
		db = db.Where("category = ?", *req.Category)
	}
	if req.Tag != nil && *req.Tag != "" {
		// tags 以 JSON 数组字符串存储，精确匹配带引号的标签
		db = db.Where("tags LIKE ?", "%\""+*req.Tag+"\"%")
	}

	option := other.MySQLOption{
		PageInfo: req.PageInfo,
		Order:    "created_at DESC",
		Where:    db,
		Preload:  []string{"User"},
	}
	return utils.MySQLPagination(&database.ForumPost{}, option)
}

// ForumDetail 帖子详情 + 评论树，浏览数 +1
func (forumService *ForumService) ForumDetail(id uint) (post ForumDetailData, err error) {
	if err = global.DB.Preload("User").First(&post.ForumPost, id).Error; err != nil {
		return post, err
	}

	// 浏览数 +1
	global.DB.Model(&database.ForumPost{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1))

	// 查询一级评论，预加载二级回复及其用户
	var comments []database.ForumComment
	if err = global.DB.Where("post_id = ? AND parent_id = 0", id).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, uuid, username, avatar, signature")
		}).
		Preload("Children.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, uuid, username, avatar, signature")
		}).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return post, err
	}
	post.Comments = comments
	return post, nil
}

// ForumLike 切换点赞，返回是否已点赞与点赞数
func (forumService *ForumService) ForumLike(req request.ForumLike) (liked bool, likeCount int, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var like database.ForumLike
		dlErr := tx.Where("post_id = ? AND user_id = ?", req.PostID, req.UserID).First(&like).Error
		if errors.Is(dlErr, gorm.ErrRecordNotFound) {
			// 未点赞 -> 创建
			if cErr := tx.Create(&database.ForumLike{PostID: req.PostID, UserID: req.UserID}).Error; cErr != nil {
				return cErr
			}
			liked = true
			if uErr := tx.Model(&database.ForumPost{}).Where("id = ?", req.PostID).
				Update("like_count", gorm.Expr("like_count + ?", 1)).Error; uErr != nil {
				return uErr
			}
		} else if dlErr != nil {
			return dlErr
		} else {
			// 已点赞 -> 取消
			if dErr := tx.Unscoped().Delete(&like).Error; dErr != nil {
				return dErr
			}
			liked = false
			if uErr := tx.Model(&database.ForumPost{}).Where("id = ?", req.PostID).
				Update("like_count", gorm.Expr("like_count - ?", 1)).Error; uErr != nil {
				return uErr
			}
		}
		// 读取最新点赞数
		var post database.ForumPost
		if pErr := tx.Select("like_count").First(&post, req.PostID).Error; pErr != nil {
			return pErr
		}
		likeCount = post.LikeCount
		return nil
	})
	return liked, likeCount, err
}

// ForumComment 发表评论
func (forumService *ForumService) ForumComment(req request.ForumComment) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 校验帖子存在
		var post database.ForumPost
		if err := tx.Select("id").First(&post, req.PostID).Error; err != nil {
			return errors.New("帖子不存在")
		}
		// 若为二级回复，校验父评论存在
		if req.ParentID != 0 {
			var parent database.ForumComment
			if err := tx.Select("id").First(&parent, req.ParentID).Error; err != nil {
				return errors.New("父评论不存在")
			}
		}
		if err := tx.Create(&database.ForumComment{
			PostID:   req.PostID,
			ParentID: req.ParentID,
			UserID:   req.UserID,
			Content:  req.Content,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&database.ForumPost{}).Where("id = ?", req.PostID).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	})
}

// ForumManageList 管理员帖子列表
func (forumService *ForumService) ForumManageList(req request.ForumManageList) (interface{}, int64, error) {
	db := global.DB.Model(&database.ForumPost{})
	if req.RoleID != appTypes.Admin {
		// 普通用户只能管理自己发布的帖子
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.Title != nil && *req.Title != "" {
		db = db.Where("title LIKE ?", "%"+*req.Title+"%")
	}
	option := other.MySQLOption{
		PageInfo: req.PageInfo,
		Order:    "created_at DESC",
		Where:    db,
		Preload:  []string{"User"},
	}
	return utils.MySQLPagination(&database.ForumPost{}, option)
}

// ForumDelete 软删除帖子（同时删除相关评论和点赞）
func (forumService *ForumService) ForumDelete(ids []uint, userID uint, roleID appTypes.RoleID) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 收集帖子标签，用于标签计数回退
		var posts []database.ForumPost
		postQuery := tx.Where("id IN ?", ids)
		if roleID != appTypes.Admin {
			// 普通用户只能删除自己发布的帖子
			postQuery = postQuery.Where("user_id = ?", userID)
		}
		if err := postQuery.Find(&posts).Error; err != nil {
			return err
		}
		if len(posts) != len(ids) {
			return errors.New("无权删除非本人发布的帖子")
		}
		// 删除点赞
		if err := tx.Unscoped().Where("post_id IN ?", ids).Delete(&database.ForumLike{}).Error; err != nil {
			return err
		}
		// 软删除评论
		if err := tx.Where("post_id IN ?", ids).Delete(&database.ForumComment{}).Error; err != nil {
			return err
		}
		// 软删除帖子
		if err := tx.Where("id IN ?", ids).Delete(&database.ForumPost{}).Error; err != nil {
			return err
		}
		// 标签计数回退
		for _, post := range posts {
			tags := []string(post.Tags)
			if err := forumService.decrTagCount(tx, tags); err != nil {
				return err
			}
		}
		return nil
	})
}

// ForumManageComments 管理员评论列表（含帖子标题、用户信息）
func (forumService *ForumService) ForumManageComments(req request.ForumManageComments) (interface{}, int64, error) {
	type commentWithPost struct {
		database.ForumComment
		PostTitle string `json:"post_title"`
	}
	db := global.DB.Model(&database.ForumComment{})
	if req.PostID != nil {
		db = db.Where("post_id = ?", *req.PostID)
	}
	if req.RoleID != appTypes.Admin {
		// 普通用户只能管理自己帖子下的评论
		db = db.Where("post_id IN (SELECT id FROM forum_posts WHERE user_id = ?)", req.UserID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	var comments []database.ForumComment
	if err := db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("id, uuid, username, avatar")
	}).Order("created_at DESC").
		Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 补充帖子标题
	postIDSet := map[uint]struct{}{}
	for i := range comments {
		postIDSet[comments[i].PostID] = struct{}{}
	}
	postIDs := make([]uint, 0, len(postIDSet))
	for id := range postIDSet {
		postIDs = append(postIDs, id)
	}
	var posts []database.ForumPost
	global.DB.Select("id, title").Where("id IN ?", postIDs).Find(&posts)
	postTitleMap := map[uint]string{}
	for _, p := range posts {
		postTitleMap[p.ID] = p.Title
	}

	list := make([]commentWithPost, 0, len(comments))
	for i := range comments {
		list = append(list, commentWithPost{
			ForumComment: comments[i],
			PostTitle:    postTitleMap[comments[i].PostID],
		})
	}
	return list, total, nil
}

// ForumCommentDelete 软删除评论，更新帖子评论数
func (forumService *ForumService) ForumCommentDelete(ids []uint, userID uint, roleID appTypes.RoleID) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var comments []database.ForumComment
		commentQuery := tx.Where("id IN ?", ids)
		if roleID != appTypes.Admin {
			// 普通用户只能删除自己帖子下的评论
			commentQuery = commentQuery.Where("post_id IN (SELECT id FROM forum_posts WHERE user_id = ?)", userID)
		}
		if err := commentQuery.Find(&comments).Error; err != nil {
			return err
		}
		if len(comments) != len(ids) {
			return errors.New("无权删除非本人帖子下的评论")
		}
		// 统计每个帖子被删的评论数
		postCount := map[uint]int{}
		for _, c := range comments {
			postCount[c.PostID]++
		}
		if err := tx.Where("id IN ?", ids).Delete(&database.ForumComment{}).Error; err != nil {
			return err
		}
		for postID, n := range postCount {
			if err := tx.Model(&database.ForumPost{}).Where("id = ?", postID).
				Update("comment_count", gorm.Expr("comment_count - ?", n)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ForumUpload 论坛图片上传，存到 uploads/forum/ 目录
func (forumService *ForumService) ForumUpload(file *multipart.FileHeader) (string, error) {
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Config.Upload.Size) {
		return "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)
	}

	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := upload.WhiteImageList[ext]; !exists {
		return "", errors.New("don't upload files that aren't image types")
	}

	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	dir := global.Config.Upload.Path + "/forum/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	dst := dir + filename
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	return "/" + dst, nil
}
