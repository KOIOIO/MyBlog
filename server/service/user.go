package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"server/utils/upload"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserService struct {
}

func (userService *UserService) Register(u database.User) (database.User, error) {
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}

	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/avatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email

	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}

	return u, nil
}

func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err
}

func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

func (userService *UserService) Logout(c *gin.Context) {
	uuid := utils.GetUUID(c)
	jwtStr := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	global.Redis.Del(uuid.String())
	_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("original password does not match the current account")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}

// UploadAvatar 上传用户头像，文件保存在本地 uploads/avatar/ 目录并更新用户头像字段
func (userService *UserService) UploadAvatar(userID uint, file *multipart.FileHeader) (string, error) {
	// 大小校验（与图片上传保持一致，单位 MB）
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Config.Upload.Size) {
		return "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)
	}

	// 类型校验，只允许图片类型
	ext := filepath.Ext(file.Filename)
	if _, exists := upload.WhiteImageList[ext]; !exists {
		return "", errors.New("don't upload files that aren't image types")
	}

	// 生成随机文件名并保存到 uploads/avatar/ 目录
	name := strings.TrimSuffix(file.Filename, ext)
	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	dir := global.Config.Upload.Path + "/avatar/"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	dst, err := os.Create(dir + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// 更新用户头像字段
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return "", err
	}
	oldAvatar := user.Avatar
	avatarURL := "/" + dir + filename
	if err := global.DB.Model(&user).Update("avatar", avatarURL).Error; err != nil {
		return "", err
	}

	// 清理旧头像文件（仅当旧头像位于头像上传目录时，避免误删默认头像或第三方头像）
	if strings.HasPrefix(oldAvatar, "/"+dir) {
		_ = os.Remove(strings.TrimPrefix(oldAvatar, "/"))
	}

	return avatarURL, nil
}

func (userService *UserService) UserWeather(ip string) (string, error) {
	// 从redis中获取天气数据，如果没有数据，则调用高德api进行查询
	result, err := global.Redis.Get("weather-" + ip).Result()
	if err != nil {
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}

		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"

		// 将天气数据存入redis
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	// 构建查询条件
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))

	var res response.UserChart

	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	// 获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	// 获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)

	for _, date := range res.DateList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}

	return res, nil
}

func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		db = db.Where("uuid = ?", info.UUID)
	}

	if info.Register != nil {
		db = db.Where("register = ?", info.Register)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.User{}, option)
}

func (userService *UserService) UserFreeze(req request.UserOperation) error {
	var user database.User
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}

	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
	if jwtStr != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
	}

	return nil
}

func (userService *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(&database.User{}, req.ID).Update("freeze", false).Error
}

func (userService *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		var userID uint
		if err := global.DB.Model(database.User{}).Where("uuid = ?", *info.UUID).Pluck("id", &userID); err != nil {
			return nil, 0, nil
		}
		db = db.Where("user_id = ?", userID)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}

	return utils.MySQLPagination(&database.Login{}, option)
}
