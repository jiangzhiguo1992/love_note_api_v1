package aliyun

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"libs/utils"
)

type (
	// 天气结果
	WeatherResult struct {
		Code int      `json:"code"`
		Data *Weather `json:"data"`
		Msg  string   `json:"msg"`
		Rc   *struct {
			C int    `json:"c"`
			P string `json:"p"`
		} `json:"rc"`
	}
	// 天气信息
	Weather struct {
		City *struct {
			CityId   int    `json:"cityId"`
			Counname string `json:"counname"`
			Name     string `json:"name"`
			Pname    string `json:"pname"`
			Timezone string `json:"timezone"`
		} `json:"city"`
		Condition *struct {
			Condition  string `json:"condition"`
			Humidity   string `json:"humidity"`
			Icon       string `json:"icon"`
			Temp       string `json:"temp"`
			Updatetime string `json:"updatetime"`
			WindDir    string `json:"windDir"`
			WindLevel  string `json:"windLevel"`
		} `json:"condition"`
		Forecast []*struct {
			ConditionDay     string `json:"conditionDay"`
			ConditionIdDay   string `json:"conditionIdDay"`
			ConditionIdNight string `json:"conditionIdNight"`
			ConditionNight   string `json:"conditionNight"`
			PredictDate      string `json:"predictDate"`
			TempDay          string `json:"tempDay"`
			TempNight        string `json:"tempNight"`
			Updatetime       string `json:"updatetime"`
			WindDirDay       string `json:"windDirDay"`
			WindDirNight     string `json:"windDirNight"`
			WindLevelDay     string `json:"windLevelDay"`
			WindLevelNight   string `json:"windLevelNight"`
		} `json:"forecast"`
	}
)

const (
	WEATHER_URL_TODAY = "http://mojibasic.market.alicloudapi.com/whapi/json/aliweather/briefcondition"
	WEATHER_URL_6DAYS = "http://mojibasic.market.alicloudapi.com/whapi/json/aliweather/briefforecast6days"
)

func GetWeatherToday(lon, lat float64) (*Weather, error) {
	return getWeather(WEATHER_URL_TODAY, lon, lat)
}

func GetWeatherForecast(lon, lat float64) (*Weather, error) {
	return getWeather(WEATHER_URL_6DAYS, lon, lat)
}

func getWeather(addressUrl string, lon, lat float64) (*Weather, error) {
	if lon == 0 && lat == 0 {
		return nil, errors.New("limit_lat_lon_nil")
	}
	weatherAppCode := utils.GetConfigStr("conf", "third.conf", "weather", "app_code")
	// form
	longitude := strconv.FormatFloat(lon, 'f', -1, 64)
	latitude := strconv.FormatFloat(lat, 'f', -1, 64)
	form := url.Values{}
	form.Set("lon", longitude)
	form.Set("lat", latitude)
	formReader := strings.NewReader(form.Encode())
	// request
	request, err := http.NewRequest(http.MethodPost, addressUrl, formReader)
	if err != nil {
		//utils.LogErr("weather", err)
		return nil, errors.New("weather_request_make_fail")
	}
	// header
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "APPCODE "+weatherAppCode)
	// http
	response, err := http.DefaultClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		//utils.LogErr("weather", err)
		return nil, errors.New("weather_info_request_fail")
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//utils.LogErr("weather", err)
		return nil, errors.New("weather_info_read_fail")
	}
	result := &WeatherResult{}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		//utils.LogErr("weather", err)
		return nil, errors.New("weather_info_decode_fail")
	}
	if result != nil && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New("weather_info_no_exist")
}
