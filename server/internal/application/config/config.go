// Package config 提供系统配置应用服务。
package config

import (
	"context"

	cconfig "server/config"
	domainconfig "server/internal/domain/config"
	"server/internal/domain/website"
)

// Service 配置应用服务（对注入的配置指针做原地修改 + YAML 持久化）。
type Service struct {
	cfg    *cconfig.Config
	store  domainconfig.ConfigStore
	images website.WebsiteImagePort
}

// NewService 构造配置应用服务。
func NewService(cfg *cconfig.Config, store domainconfig.ConfigStore, images website.WebsiteImagePort) *Service {
	return &Service{cfg: cfg, store: store, images: images}
}

// UpdateWebsite 更新网站配置（图片类别联动：logo/fullLogo/qqImage/wechatImage 差异处理）。
func (s *Service) UpdateWebsite(ctx context.Context, website cconfig.Website) error {
	oldArray := []string{
		s.cfg.Website.Logo,
		s.cfg.Website.FullLogo,
		s.cfg.Website.QQImage,
		s.cfg.Website.WechatImage,
	}
	newArray := []string{
		website.Logo,
		website.FullLogo,
		website.QQImage,
		website.WechatImage,
	}
	added, removed := domainconfig.DiffArrays(oldArray, newArray)
	if len(removed) > 0 {
		if err := s.images.InitCategory(ctx, removed); err != nil {
			return err
		}
	}
	if len(added) > 0 {
		if err := s.images.ChangeCategory(ctx, added, "系统"); err != nil {
			return err
		}
	}
	s.cfg.Website = website
	return s.store.Save(ctx)
}

// UpdateSystem 更新系统配置（仅 useMultipoint/sessionsSecret/ossType）。
func (s *Service) UpdateSystem(ctx context.Context, system cconfig.System) error {
	s.cfg.System.UseMultipoint = system.UseMultipoint
	s.cfg.System.SessionsSecret = system.SessionsSecret
	s.cfg.System.OssType = system.OssType
	return s.store.Save(ctx)
}

// UpdateEmail 更新邮箱配置。
func (s *Service) UpdateEmail(ctx context.Context, email cconfig.Email) error {
	s.cfg.Email = email
	return s.store.Save(ctx)
}

// UpdateQiniu 更新七牛云配置。
func (s *Service) UpdateQiniu(ctx context.Context, qiniu cconfig.Qiniu) error {
	s.cfg.Qiniu = qiniu
	return s.store.Save(ctx)
}

// UpdateJwt 更新 JWT 配置。
func (s *Service) UpdateJwt(ctx context.Context, jwt cconfig.Jwt) error {
	s.cfg.Jwt = jwt
	return s.store.Save(ctx)
}

// UpdateGaode 更新高德配置。
func (s *Service) UpdateGaode(ctx context.Context, gaode cconfig.Gaode) error {
	s.cfg.Gaode = gaode
	return s.store.Save(ctx)
}
