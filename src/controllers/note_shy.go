package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerShy(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_shy")
	if r.Method == http.MethodPost {
		PostShy(w, r)
	} else if r.Method == http.MethodDelete {
		DelShy(w, r)
	} else if r.Method == http.MethodGet {
		GetShy(w, r)
	} else {
		response405(w, r)
	}
}

// PostShy
func PostShy(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	shy := &entity.Shy{}
	checkRequestBody(w, r, shy)
	// 开始插入
	shy, err := services.AddShy(user.Id, couple.Id, shy)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Shy *entity.Shy `json:"shy"`
	}{shy})
}

// DelShy
func DelShy(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	// 开始删除
	err := services.DelShy(user.Id, couple.Id, sid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetShy
func GetShy(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	date, _ := strconv.ParseBool(values.Get("date"))
	if date {
		year, _ := strconv.Atoi(values.Get("year"))
		month, _ := strconv.Atoi(values.Get("month"))
		// 数据
		shyList, err := services.GetShyListByCoupleYearMonth(user.Id, couple.Id, year, month)
		response417ErrDialog(w, r, err) // dialog
		// 返回
		response200Data(w, r, struct {
			ShyList []*entity.Shy `json:"shyList"`
		}{shyList})
	} else {
		response405(w, r)
	}
}
