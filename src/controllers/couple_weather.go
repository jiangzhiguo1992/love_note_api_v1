package controllers

//import (
//	"net/http"
//	"strconv"
//
//	"libs/utils"
//	"models/mysql"
//	"services"
//)
//
//func HandlerWeather(w http.ResponseWriter, r *http.Request) {
//	checkModelStatus(w, r, "couple_weather")
//	if r.Method == http.MethodGet {
//		GetWeather(w, r)
//	} else {
//		response405(w, r)
//	}
//}
//
//// GetWeather
//func GetWeather(w http.ResponseWriter, r *http.Request) {
//	user, _ := getTokenCouple(r)
//	couple := user.Couple
//	// 接受参数
//	values := r.URL.Query()
//	forecast, _ := strconv.ParseBool(values.Get("forecast"))
//	if forecast {
//		var weatherForecastListMe, weatherForecastListTa []*services.WeatherForecast
//		var showMe, showTa string
//		// 我的地址，没有上传我的地址，则从后台获取上次上传的地址
//		myPlace, err := services.GetPlaceLatestByUserCouple(user.Id, couple.Id)
//		if err != nil {
//			showMe = err.Error()
//		} else if myPlace == nil || (myPlace.Longitude == 0 && myPlace.Latitude == 0) {
//			showMe = "limit_lat_lon_nil"
//		} else {
//			weatherForecastListMe, err = services.GetWeatherForecast(myPlace.Longitude, myPlace.Latitude)
//			if err != nil {
//				showMe = err.Error()
//			}
//		}
//		// ta的地址，不用上传，直接后台获取
//		taId := services.GetTaId(user)
//		taPlace, err := services.GetPlaceLatestByUserCouple(taId, couple.Id)
//		if err != nil {
//			showTa = err.Error()
//		} else if taPlace == nil || (taPlace.Longitude == 0 && taPlace.Latitude == 0) {
//			showTa = "limit_lat_lon_nil"
//		} else {
//			weatherForecastListTa, err = services.GetWeatherForecast(taPlace.Longitude, taPlace.Latitude)
//			if err != nil {
//				showTa = err.Error()
//			}
//		}
//		// show
//		language := "zh-cn"
//		entry, err := mysql.GetEntryLatestByUser(user.Id)
//		if err == nil && entry != nil {
//			language = entry.Language
//		}
//		if len(showMe) > 0 {
//			showMe = utils.GetLanguage(language, showMe)
//		}
//		if len(showTa) > 0 {
//			showTa = utils.GetLanguage(language, showTa)
//		}
//		// 封装
//		weatherForecastMe := &services.WeatherForecastInfo{
//			Show:                showMe,
//			WeatherForecastList: weatherForecastListMe,
//		}
//		weatherForecastTa := &services.WeatherForecastInfo{
//			Show:                showTa,
//			WeatherForecastList: weatherForecastListTa,
//		}
//		// 返回
//		response200Data(w, r, struct {
//			WeatherForecastMe *services.WeatherForecastInfo `json:"weatherForecastMe"`
//			WeatherForecastTa *services.WeatherForecastInfo `json:"weatherForecastTa"`
//		}{weatherForecastMe, weatherForecastTa})
//	} else {
//		response405(w, r)
//	}
//}
