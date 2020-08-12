package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerSleep(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_sleep")
	if r.Method == http.MethodPost {
		PostSleep(w, r)
	} else if r.Method == http.MethodDelete {
		DelSleep(w, r)
	} else if r.Method == http.MethodGet {
		GetSleep(w, r)
	} else {
		response405(w, r)
	}
}

// PostSleep
func PostSleep(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	sleep := &entity.Sleep{}
	checkRequestBody(w, r, sleep)
	// 开始插入
	sleep, err := services.AddSleep(user.Id, couple.Id, sleep)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Sleep *entity.Sleep `json:"sleep"`
	}{sleep})
}

// DelSleep
func DelSleep(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	// 开始删除
	err := services.DelSleep(user.Id, couple.Id, sid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetSleep
func GetSleep(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	latest, _ := strconv.ParseBool(values.Get("latest"))
	date, _ := strconv.ParseBool(values.Get("date"))
	if latest {
		taId := services.GetTaId(user)
		sleepMe, _ := services.GetSleepLatestByUserCouple(user.Id, couple.Id)
		sleepTa, _ := services.GetSleepLatestByUserCouple(taId, couple.Id)
		// 返回
		response200Data(w, r, struct {
			SleepMe *entity.Sleep `json:"sleepMe"`
			SleepTa *entity.Sleep `json:"sleepTa"`
		}{sleepMe, sleepTa})
	} else if date {
		year, _ := strconv.Atoi(values.Get("year"))
		month, _ := strconv.Atoi(values.Get("month"))
		// 数据
		sleepList, err := services.GetSleepListByCoupleYearMonth(user.Id, couple.Id, year, month)
		response417ErrDialog(w, r, err) // dialog
		// 返回
		response200Data(w, r, struct {
			SleepList []*entity.Sleep `json:"sleepList"`
		}{sleepList})
	} else {
		response405(w, r)
	}
}
