package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerFood(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_food")
	if r.Method == http.MethodPost {
		PostFood(w, r)
	} else if r.Method == http.MethodDelete {
		DelFood(w, r)
	} else if r.Method == http.MethodPut {
		PutFood(w, r)
	} else if r.Method == http.MethodGet {
		GetFood(w, r)
	} else {
		response405(w, r)
	}
}

// PostFood
func PostFood(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	food := &entity.Food{}
	checkRequestBody(w, r, food)
	// 开始插入
	food, err := services.AddFood(user.Id, couple.Id, food)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Food *entity.Food `json:"food"`
	}{food})
}

// DelFood
func DelFood(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	fid, _ := strconv.ParseInt(values.Get("fid"), 10, 64)
	// 开始删除
	err := services.DelFood(user.Id, couple.Id, fid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutFood
func PutFood(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	food := &entity.Food{}
	checkRequestBody(w, r, food)
	// 开始插入
	food, err := services.UpdateFood(user.Id, couple.Id, food)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Food *entity.Food `json:"food"`
	}{food})
}

// GetFood
func GetFood(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		foodList, err := services.GetFoodListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			FoodList []*entity.Food `json:"foodList"`
		}{foodList})
	} else {
		response405(w, r)
	}
}
