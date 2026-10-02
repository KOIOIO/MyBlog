// Package user 提供用户用例编排（注册/登录/找回密码/资料/冻结/图表等）。
package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/config"
	"server/internal/common/crypto"
	"server/internal/common/errs"
	"server/internal/domain/auth"
	"server/internal/domain/user"

	"github.com/gofrs/uuid"
)

// whiteImageList 允许上传的图片扩展名白名单（与 utils/upload.WhiteImageList 一致）。
var whiteImageList = map[string]struct{}{
	".jpg":  {},
	".png":  {},
	".jpeg": {},
	".ico":  {},
	".tiff": {},
	".gif":  {},
	".svg":  {},
	".webp": {},
}

// ChartResult 用户图表数据。
type ChartResult struct {
	DateList     []string
	LoginData    []int
	RegisterData []int
}

// UserService 用户用例。
type UserService struct {
	users  user.UserRepository
	logins user.LoginRecordRepository
	cache  user.Cache
	geo    user.GeoProvider
	auth   auth.AuthPort
	cfg    *config.Config
}

// NewUserService 构造用户用例。
func NewUserService(users user.UserRepository, logins user.LoginRecordRepository, cache user.Cache, geo user.GeoProvider, authPort auth.AuthPort, cfg *config.Config) *UserService {
	return &UserService{users: users, logins: logins, cache: cache, geo: geo, auth: authPort, cfg: cfg}
}

// Register 注册（邮箱重复校验 + 密码散列 + 默认信息）。
func (s *UserService) Register(ctx context.Context, username, password, email string) (*user.User, error) {
	if _, err := s.users.FindByEmail(ctx, email); !errors.Is(err, errs.ErrNotFound) {
		return nil, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}

	u := &user.User{Username: username, Password: password, Email: email}
	u.ResetRegisterInfo(uuid.Must(uuid.NewV4()), crypto.BcryptHash(password), "/image/avatar.jpg")

	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// EmailLogin 邮箱登录。
func (s *UserService) EmailLogin(ctx context.Context, email, password string) (*user.User, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !crypto.BcryptCheck(password, u.Password) {
		return nil, errors.New("incorrect email or password")
	}
	return u, nil
}

// ForgotPassword 找回密码。
func (s *UserService) ForgotPassword(ctx context.Context, email, newPassword string) error {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	u.SetPasswordHash(crypto.BcryptHash(newPassword))
	return s.users.Update(ctx, u)
}

// Card 用户卡片信息。
func (s *UserService) Card(ctx context.Context, uuidStr string) (*user.User, error) {
	uid, err := uuid.FromString(uuidStr)
	if err != nil {
		return nil, err
	}
	return s.users.FindByUUID(ctx, uid)
}

// Logout 登出（删除会话 + 刷新令牌加入黑名单，错误忽略保持旧行为）。
func (s *UserService) Logout(ctx context.Context, uuidVal uuid.UUID, refreshToken string) error {
	_ = s.auth.DelSession(ctx, uuidVal.String())
	if refreshToken != "" {
		_ = s.auth.Blacklist(ctx, refreshToken)
	}
	return nil
}

// ResetPassword 修改密码。
func (s *UserService) ResetPassword(ctx context.Context, userID uint, password, newPassword string) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !crypto.BcryptCheck(password, u.Password) {
		return errors.New("original password does not match the current account")
	}
	u.SetPasswordHash(crypto.BcryptHash(newPassword))
	return s.users.Update(ctx, u)
}

// Info 个人信息。
func (s *UserService) Info(ctx context.Context, userID uint) (*user.User, error) {
	return s.users.FindByID(ctx, userID)
}

// ChangeInfo 修改个人信息。
func (s *UserService) ChangeInfo(ctx context.Context, userID uint, username, address, signature string) error {
	return s.users.UpdateInfo(ctx, userID, username, address, signature)
}

// UploadAvatar 上传头像（本地存储，行为与旧实现一致）。
func (s *UserService) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error) {
	// 大小校验（单位 MB）
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(s.cfg.Upload.Size) {
		return "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, s.cfg.Upload.Size)
	}

	// 类型校验
	ext := filepath.Ext(file.Filename)
	if _, exists := whiteImageList[ext]; !exists {
		return "", errors.New("don't upload files that aren't image types")
	}

	// 生成随机文件名并保存
	name := strings.TrimSuffix(file.Filename, ext)
	filename := crypto.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	dir := s.cfg.Upload.Path + "/avatar/"

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

	// 更新头像字段（返回更新前用户，供旧头像清理）
	old, err := s.users.UpdateAvatar(ctx, userID, "/"+dir+filename)
	if err != nil {
		return "", err
	}

	// 清理旧头像文件（仅当位于头像上传目录时）
	if strings.HasPrefix(old.Avatar, "/"+dir) {
		_ = os.Remove(strings.TrimPrefix(old.Avatar, "/"))
	}

	return "/" + dir + filename, nil
}

// Weather 获取天气（Redis 缓存 + 高德查询）。
func (s *UserService) Weather(ctx context.Context, ip string) (string, error) {
	result, err := s.cache.Get(ctx, "weather-"+ip)
	if err != nil {
		loc, err := s.geo.LocationByIP(ctx, ip)
		if err != nil {
			return "", err
		}
		live, err := s.geo.WeatherByAdcode(ctx, loc.Adcode)
		if err != nil {
			return "", err
		}

		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"

		if err := s.cache.Set(ctx, "weather-"+ip, weather, time.Hour); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

// Chart 用户图表数据（登录/注册人数）。
func (s *UserService) Chart(ctx context.Context, days int) (*ChartResult, error) {
	loginCounts, _ := s.logins.CountByDate(ctx, days)
	registerCounts, _ := s.users.CountByDate(ctx, days)

	res := &ChartResult{}
	startDate := time.Now().AddDate(0, 0, -days)
	for i := 1; i <= days; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	for _, date := range res.DateList {
		res.LoginData = append(res.LoginData, loginCounts[date])
		res.RegisterData = append(res.RegisterData, registerCounts[date])
	}
	return res, nil
}

// List 用户列表。
func (s *UserService) List(ctx context.Context, cond user.ListCond) ([]*user.User, int64, error) {
	return s.users.Page(ctx, cond)
}

// Freeze 冻结用户（会话中的刷新令牌加入黑名单）。
func (s *UserService) Freeze(ctx context.Context, id uint) error {
	u, err := s.users.UpdateFreeze(ctx, id, true)
	if err != nil {
		return err
	}
	jwtStr, _ := s.auth.GetSession(ctx, u.UUID.String())
	if jwtStr != "" {
		_ = s.auth.Blacklist(ctx, jwtStr)
	}
	return nil
}

// Unfreeze 解冻用户。
func (s *UserService) Unfreeze(ctx context.Context, id uint) error {
	_, err := s.users.UpdateFreeze(ctx, id, false)
	return err
}

// LoginList 登录日志列表。
func (s *UserService) LoginList(ctx context.Context, cond user.LoginListCond) ([]*user.LoginRecord, int64, error) {
	return s.logins.Page(ctx, cond)
}
