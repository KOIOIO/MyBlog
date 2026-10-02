package hotsearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/domain/website"
	"strconv"
	"strings"
	"time"
)

type Bilibili struct {
}

type bilibiliResponse struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Pic   string `json:"pic"`
			Bvid  string `json:"bvid"`
			Stat  struct {
				View int64 `json:"view"`
			} `json:"stat"`
		} `json:"list"`
	} `json:"data"`
}

// GetHotSearchData 获取哔哩哔哩热门视频榜（官方公开接口）
func (*Bilibili) GetHotSearchData(maxNum int) (website.HotSearchData, error) {
	// 设置超时，避免上游无响应导致请求挂起
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/popular?ps=30&pn=1", nil)
	if err != nil {
		return website.HotSearchData{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := client.Do(req)
	if err != nil {
		return website.HotSearchData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return website.HotSearchData{}, err
	}

	var result bilibiliResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return website.HotSearchData{}, err
	}
	if result.Code != 0 {
		return website.HotSearchData{}, fmt.Errorf("bilibili api error: code=%d", result.Code)
	}

	updateTime := time.Now().Format("2006-01-02 15:04:05")

	var hotList []website.HotItem
	for i, item := range result.Data.List {
		if i >= maxNum {
			break
		}
		hotList = append(hotList, website.HotItem{
			Index:       i + 1,
			Title:       item.Title,
			Description: cleanDescription(item.Desc),
			Image:       strings.Replace(item.Pic, "http://", "https://", 1),
			Popularity:  formatViewCount(item.Stat.View),
			URL:         "https://www.bilibili.com/video/" + item.Bvid,
		})
	}

	return website.HotSearchData{Source: "哔哩哔哩热榜", UpdateTime: updateTime, HotList: hotList}, nil
}

// cleanDescription 去除简介中的换行并截断，避免表格展示过长
func cleanDescription(desc string) string {
	desc = strings.Join(strings.Fields(desc), " ")
	runes := []rune(desc)
	if len(runes) > 150 {
		return string(runes[:150]) + "..."
	}
	return desc
}

// formatViewCount 将播放量格式化为 万/亿 可读形式
func formatViewCount(n int64) string {
	switch {
	case n >= 100000000:
		return fmt.Sprintf("%.1f亿", float64(n)/100000000)
	case n >= 10000:
		return fmt.Sprintf("%.1f万", float64(n)/10000)
	default:
		return strconv.FormatInt(n, 10)
	}
}
