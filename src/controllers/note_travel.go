package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerTravel(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_travel")
	if r.Method == http.MethodPost {
		PostTravel(w, r)
	} else if r.Method == http.MethodDelete {
		DelTravel(w, r)
	} else if r.Method == http.MethodPut {
		PutTravel(w, r)
	} else if r.Method == http.MethodGet {
		GetTravel(w, r)
	} else {
		response405(w, r)
	}
}

// PostTravel
func PostTravel(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	travel := &entity.Travel{}
	checkRequestBody(w, r, travel)
	// 开始插入
	travel, err := services.AddTravel(user.Id, couple.Id, travel)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Travel *entity.Travel `json:"travel"`
	}{travel})
}

// DelTravel
func DelTravel(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	tid, _ := strconv.ParseInt(values.Get("tid"), 10, 64)
	// 开始删除
	err := services.DelTravel(user.Id, couple.Id, tid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutTravel
func PutTravel(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	travel := &entity.Travel{}
	checkRequestBody(w, r, travel)
	// 开始插入
	travel, err := services.UpdateTravel(user.Id, couple.Id, travel)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Travel *entity.Travel `json:"travel"`
	}{travel})
}

// GetTravel
func GetTravel(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	tid, _ := strconv.ParseInt(values.Get("tid"), 10, 64)
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		travelList, err := services.GetTravelListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			TravelList []*entity.Travel `json:"travelList"`
		}{travelList})
	} else if tid > 0 {
		travel, err := services.GetTravelByIdWithForeign(user.Id, couple.Id, tid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Travel *entity.Travel `json:"travel"`
		}{travel})
	} else {
		response405(w, r)
	}
}
