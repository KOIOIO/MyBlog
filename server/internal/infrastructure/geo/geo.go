// Package geo 提供高德地图 API 实现（IP 定位 / 实时天气）。
package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"server/internal/common/httpclient"
	"server/internal/domain/user"

	"go.uber.org/zap"
)

// Client 高德 API 客户端。
type Client struct {
	key    string
	logger *zap.Logger
}

// NewClient 构造高德客户端。
func NewClient(key string, logger *zap.Logger) *Client {
	return &Client{key: key, logger: logger}
}

// ipResponse 高德 IP 定位响应。
type ipResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

// weatherResponse 高德天气响应。
type weatherResponse struct {
	Status   string      `json:"status"`
	Count    string      `json:"count"`
	Info     string      `json:"info"`
	InfoCode string      `json:"infocode"`
	Lives    []live      `json:"lives"`
	Forecast interface{} `json:"forecast"`
}

type live struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTime    string `json:"reporttime"`
}

// LocationByIP 根据 IP 获取地理位置信息（内网 IP 不传参，由高德按出口 IP 定位）。
func (c *Client) LocationByIP(ctx context.Context, ip string) (user.GeoInfo, error) {
	data := ipResponse{}
	params := map[string]string{"key": c.key}
	if !isPrivateIP(ip) {
		params["ip"] = ip
	}
	res, err := httpclient.Request("https://restapi.amap.com/v3/ip", "GET", nil, params, nil)
	if err != nil {
		return user.GeoInfo{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return user.GeoInfo{}, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return user.GeoInfo{}, err
	}
	if err := json.Unmarshal(byteData, &data); err != nil {
		return user.GeoInfo{}, err
	}
	return user.GeoInfo{Province: data.Province, City: data.City, Adcode: data.Adcode}, nil
}

// WeatherByAdcode 根据城市编码获取实时天气。
func (c *Client) WeatherByAdcode(ctx context.Context, adcode string) (user.WeatherInfo, error) {
	data := weatherResponse{}
	params := map[string]string{"city": adcode, "key": c.key}
	res, err := httpclient.Request("https://restapi.amap.com/v3/weather/weatherInfo", "GET", nil, params, nil)
	if err != nil {
		return user.WeatherInfo{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return user.WeatherInfo{}, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return user.WeatherInfo{}, err
	}
	if err := json.Unmarshal(byteData, &data); err != nil {
		return user.WeatherInfo{}, err
	}
	if len(data.Lives) == 0 {
		return user.WeatherInfo{}, fmt.Errorf("no live weather data available")
	}
	l := data.Lives[0]
	return user.WeatherInfo{
		Province:      l.Province,
		City:          l.City,
		Weather:       l.Weather,
		Temperature:   l.Temperature,
		WindDirection: l.WindDirection,
		WindPower:     l.WindPower,
		Humidity:      l.Humidity,
	}, nil
}

// isPrivateIP 判断 IP 是否为内网/回环/保留地址；解析失败按内网处理。
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
}

// 接口编译期校验。
var _ user.GeoProvider = (*Client)(nil)
