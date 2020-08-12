package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerCoupleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetCoupleHome(w, r)
	} else {
		response405(w, r)
	}
}

// GetCoupleHome
func GetCoupleHome(w http.ResponseWriter, r *http.Request) {
	// me
	me, _ := getTokenCouple(r)
	couple := me.Couple
	if couple == nil || couple.Id <= 0 {
		response200Data(w, r, struct {
			User *entity.User `json:"user"`
		}{me})
	}
	// ta
	taId := services.GetTaId(me)
	ta, _ := services.GetUserById(taId)
	if ta != nil {
		ta.Password = ""
		ta.UserToken = ""
	}
	// togetherDay
	//togetherDay := services.GetCoupleTogetherDay(couple.Id)
	// wallPaper
	var wallPaper *entity.WallPaper
	if getModelStatus("couple_wall") {
		wallPaper, _ = services.GetWallPaperByCouple(couple.Id)
	}
	// place
	//var placeMe, placeTa *entity.Place
	//if getModelStatus("couple_place") {
	//	placeMe, _ = services.GetPlaceLatestByUser(me.Id)
	//	placeTa, _ = services.GetPlaceLatestByUser(taId)
	//}
	// weather
	//var weatherTodayMe, weatherTodayTa *services.WeatherToday
	//if getModelStatus("couple_weather") {
	//	if placeMe != nil {
	//		weatherTodayMe, _ = services.GetWeatherToday(placeMe.Longitude, placeMe.Latitude)
	//	}
	//	if placeTa != nil {
	//		weatherTodayTa, _ = services.GetWeatherToday(placeTa.Longitude, placeTa.Latitude)
	//	}
	//}
	// 返回
	response200Data(w, r, struct {
		User *entity.User `json:"user"`
		Ta   *entity.User `json:"ta"`
		//TogetherDay int               `json:"togetherDay"`
		WallPaper *entity.WallPaper `json:"wallPaper"`
		//PlaceMe        *entity.Place          `json:"placeMe"`
		//PlaceTa        *entity.Place          `json:"placeTa"`
		//WeatherTodayMe *services.WeatherToday `json:"weatherTodayMe"`
		//WeatherTodayTa *services.WeatherToday `json:"weatherTodayTa"`
	}{
		me,
		ta,
		//togetherDay,
		wallPaper,
		//placeMe,
		//placeTa,
		//weatherTodayMe,
		//weatherTodayTa,
	})
}
