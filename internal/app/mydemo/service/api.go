package service

import (
	"context"
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type WeatherResp struct {
	Reason     string
	Error_code int
	Result     WeatherResult
}

type WeatherResult struct {
	City     string
	Realtime RealTimeWeather
	Future   []FutureWeather
}

type RealTimeWeather struct {
	Temperature string
	Humidity    string
	Info        string
	Wid         string
	Direct      string
	Power       string
	Aqi         string
}

type FutureWeather struct {
	Date        string
	Temperature string
	Weather     string
	Wid         Wid
	Direct      string
}

type Wid struct {
	Day   string
	Night string
}

type DataStruct struct {
	Weather model.WeatherItem
	Slices  []model.ResponseCheckinItem
}

type Weather struct {
	Temperature string
	Weather     string
}

func GetWeather(ctx context.Context, city string) (todayWeather Weather, err error) {
	apiUrl := "http://apis.juhe.cn/simpleWeather/query"
	apiKey := Cfg.WeatherAppKey
	data := url.Values{}
	data.Set("key", apiKey)
	data.Set("city", city)
	str := fmt.Sprintf("key=%s&city=%s&sign_salt=%s", apiKey, city, Cfg.WeatherSignSalt)
	md5Str := utils.MakeMd5(str)
	data.Set("sign", md5Str) //实际不需要这个参数，目前为练习

	//resp, err := http.Get(apiUrl + "?" + data.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", apiUrl+"?"+data.Encode(), nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	weatherdate, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var weatherResp WeatherResp
	err = json.Unmarshal(weatherdate, &weatherResp)
	if err != nil {
		return
	}

	errCode := weatherResp.Error_code
	if errCode != 0 {
		err = errors.New("errCode != 0")
		return
	}

	todayWeatherInfo := weatherResp.Result.Realtime.Info
	todayTemperature := weatherResp.Result.Realtime.Temperature
	todayWeather = Weather{
		Weather:     todayWeatherInfo,
		Temperature: todayTemperature,
	}
	return
}
