package task

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"server/global"
)

func RegisterScheduledTasks(c *cron.Cron) error {
	// 文章浏览量同步已迁移至 internal/application/article（bootstrap.InitCron 注册），此处不再注册旧任务。
	if _, err := c.AddFunc("@hourly", func() {
		if err := GetHotListSyncTask(); err != nil {
			global.Log.Error("Failed to get hot list:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	if _, err := c.AddFunc("@daily", func() {
		if err := GetCalendarSyncTask(); err != nil {
			global.Log.Error("Failed to get calendar:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	return nil
}
