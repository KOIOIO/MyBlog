package task

import (
	"encoding/json"
	"server/global"
	"server/utils"
	"time"
)

func GetCalendarSyncTask() error {
	dataStr := time.Now().Format("2006/0102")
	calendar, err := utils.GetCalendar(dataStr)
	if err != nil {
		return err
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return err
	}
	if err := global.Redis.Set("calendar-"+dataStr, data, time.Hour*24).Err(); err != nil {
		return err
	}
	return nil
}
