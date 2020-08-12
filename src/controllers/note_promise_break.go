package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerPromiseBreak(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_promise_break")
	if r.Method == http.MethodPost {
		PostPromiseBreak(w, r)
	} else if r.Method == http.MethodDelete {
		DelPromiseBreak(w, r)
	} else if r.Method == http.MethodGet {
		GetPromiseBreak(w, r)
	} else {
		response405(w, r)
	}
}

// PostPromiseBreak
func PostPromiseBreak(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	promiseBreak := &entity.PromiseBreak{}
	checkRequestBody(w, r, promiseBreak)
	// 开始插入
	promiseBreak, err := services.AddPromiseBreak(user.Id, couple.Id, promiseBreak)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		PromiseBreak *entity.PromiseBreak `json:"promiseBreak"`
	}{promiseBreak})
}

// DelPromiseBreak
func DelPromiseBreak(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	pbid, _ := strconv.ParseInt(values.Get("pbid"), 10, 64)
	// 开始删除
	err := services.DelPromiseBreak(user.Id, couple.Id, pbid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetPromiseBreak
func GetPromiseBreak(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	if pid > 0 {
		page, _ := strconv.Atoi(values.Get("page"))
		promiseBreakList, err := services.GetPromiseBreakListByCouplePromise(user.Id, couple.Id, pid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PromiseBreakList []*entity.PromiseBreak `json:"promiseBreakList"`
		}{promiseBreakList})
	} else {
		response405(w, r)
	}
}
