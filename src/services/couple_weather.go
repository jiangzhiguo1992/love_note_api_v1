package services

import (
	"libs/aliyun"
	"libs/utils"
)

// GetWeatherToday 今日天气
func GetWeatherToday(lon, lat float64) (*WeatherToday, error) {
	result, err := aliyun.GetWeatherToday(lon, lat)
	if result == nil || result.Condition == nil || err != nil {
		return nil, err
	}
	condition := result.Condition
	// 天气信息封装
	weather := &WeatherToday{
		Condition: condition.Condition,
		Icon:      condition.Icon,
		Temp:      condition.Temp,
		//Humidity:  condition.Humidity + "%RH",
		Humidity:  "",
		WindLevel: condition.WindLevel,
		WindDir:   condition.WindDir,
		UpdateAt:  utils.GetUnixByTimeFormat(condition.Updatetime, WEATHER_TIME_FORMAT_ALL),
	}
	return weather, nil
}

// GetWeatherForecast 天气预报6天
func GetWeatherForecast(lon, lat float64) ([]*WeatherForecast, error) {
	result, err := aliyun.GetWeatherForecast(lon, lat)
	if result == nil || result.Forecast == nil || len(result.Forecast) <= 0 || err != nil {
		return nil, err
	}
	forecastList := result.Forecast
	// 天气预报列表封装
	weatherList := make([]*WeatherForecast, 0)
	for i := 0; i < len(forecastList); i++ {
		forecast := forecastList[i]
		if forecast == nil {
			continue
		}
		weather := &WeatherForecast{
			TimeAt:         utils.GetUnixByTimeFormat(forecast.PredictDate, WEATHER_TIME_FORMAT_DAY),
			ConditionDay:   forecast.ConditionDay,
			ConditionNight: forecast.ConditionNight,
			IconDay:        forecast.ConditionIdDay,
			IconNight:      forecast.ConditionIdNight,
			TempDay:        forecast.TempDay,
			TempNight:      forecast.TempNight,
			WindDay:        forecast.WindLevelDay + "级" + forecast.WindDirDay,
			WindNight:      forecast.WindLevelNight + "级" + forecast.WindDirNight,
			UpdateAt:       utils.GetUnixByTimeFormat(forecast.Updatetime, WEATHER_TIME_FORMAT_ALL),
		}
		weatherList = append(weatherList, weather)
	}
	return weatherList, nil
}
