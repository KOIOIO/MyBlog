// Package calendar 提供农历日历防腐层（6tail lunar-go 本地计算 + Redis 缓存 24 小时）。
// 计算逻辑为 utils/calendar 纯搬运，产物类型为 domain/website.Calendar。
package calendar

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"server/internal/domain/website"

	"github.com/6tail/lunar-go/calendar"
	"github.com/go-redis/redis"
)

var solarTerm = []string{
	"立春", "雨水", "惊蛰", "春分", "清明", "谷雨",
	"立夏", "小满", "芒种", "夏至", "小暑", "大暑",
	"立秋", "处暑", "白露", "秋分", "寒露", "霜降",
	"立冬", "小雪", "大雪", "冬至", "小寒", "大寒",
}

// westernZodiac 西方星座分界日（每月对应日期之前归入上一个星座）。
var zodiacStartDay = [12]int{20, 19, 21, 20, 21, 22, 23, 23, 23, 23, 23, 22}
var zodiacNames = [12]string{"摩羯", "水瓶", "双鱼", "白羊", "金牛", "双子", "巨蟹", "狮子", "处女", "天秤", "天蝎", "射手"}

// listToStr 将 container.List 中的宜/忌词条拼接为空格分隔字符串。
func listToStr(l *list.List) string {
	items := make([]string, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(string))
	}
	return strings.Join(items, " ")
}

// westernZodiac 根据公历月日计算西方星座。
func westernZodiac(month, day int) string {
	idx := month - 1
	if day < zodiacStartDay[month-1] {
		idx = (month + 10) % 12
	}
	return zodiacNames[idx] + "座"
}

// compute 本地计算公历对应的农历信息（dateStr 格式 "2006/0102"）。
func compute(dateStr string) (website.Calendar, error) {
	if len(dateStr) != 9 || dateStr[4] != '/' {
		return website.Calendar{}, fmt.Errorf("invalid dateStr: %s", dateStr)
	}
	year, err := strconv.Atoi(dateStr[0:4])
	if err != nil {
		return website.Calendar{}, err
	}
	month, err := strconv.Atoi(dateStr[5:7])
	if err != nil {
		return website.Calendar{}, err
	}
	day, err := strconv.Atoi(dateStr[7:9])
	if err != nil {
		return website.Calendar{}, err
	}

	solar := calendar.NewSolarFromYmd(year, month, day)
	l := solar.GetLunar()

	// 公历日期 + 星期
	date := fmt.Sprintf("%d年%d月%d日 ", year, month, day) + solar.GetWeekInChinese()

	// 节气：上一个节气到今天的天数，下一个节气距今天数
	solarTermStr := ""
	if prevJq := l.GetPrevJieQi(); prevJq != nil {
		daysSince := solar.Subtract(prevJq.GetSolar())
		daysToNext := ""
		if nextJq := l.GetNextJieQi(); nextJq != nil {
			daysToNext = "　距离" + nextJq.GetName() + "还有" +
				strconv.Itoa(nextJq.GetSolar().Subtract(solar)) + "天"
		}
		solarTermStr = fmt.Sprintf("%s第%d天%s", prevJq.GetName(), daysSince, daysToNext)
	}

	// 年第几天
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	target := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dayOfYear := int(target.Sub(yearStart).Hours()/24) + 1

	return website.Calendar{
		Date:         date,
		LunarDate:    l.GetMonthInChinese() + "月" + l.GetDayInChinese(),
		Ganzhi:       l.GetYearInGanZhi() + "年 " + l.GetMonthInGanZhi() + "月 " + l.GetDayInGanZhi() + "日",
		Zodiac:       westernZodiac(month, day),
		DayOfYear:    "今年第" + strconv.Itoa(dayOfYear) + "天",
		SolarTerm:    solarTermStr,
		Auspicious:   listToStr(l.GetDayYi()),
		Inauspicious: listToStr(l.GetDayJi()),
	}, nil
}

// Provider 日历提供者（Redis 缓存 24 小时）。
type Provider struct {
	cache *redis.Client
}

// NewProvider 构造日历提供者。
func NewProvider(cache *redis.Client) *Provider {
	return &Provider{cache: cache}
}

var _ website.CalendarProvider = (*Provider)(nil)

// GetCalendarByDate 获取日历信息（缓存优先）。
func (p *Provider) GetCalendarByDate(_ context.Context, dateStr string) (website.Calendar, error) {
	result, err := p.cache.Get("calendar-" + dateStr).Result()
	if err != nil {
		cal, err := compute(dateStr)
		if err != nil {
			return website.Calendar{}, err
		}
		data, err := json.Marshal(cal)
		if err != nil {
			return website.Calendar{}, err
		}
		if err := p.cache.Set("calendar-"+dateStr, data, time.Hour*24).Err(); err != nil {
			return website.Calendar{}, err
		}
		return cal, nil
	}
	var cal website.Calendar
	if err := json.Unmarshal([]byte(result), &cal); err != nil {
		return website.Calendar{}, err
	}
	return cal, nil
}
