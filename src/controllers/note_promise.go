package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerPromise(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_promise")
	if r.Method == http.MethodPost {
		PostPromise(w, r)
	} else if r.Method == http.MethodDelete {
		DelPromise(w, r)
	} else if r.Method == http.MethodPut {
		PutPromise(w, r)
	} else if r.Method == http.MethodGet {
		GetPromise(w, r)
	} else {
		response405(w, r)
	}
}

// PostPromise
func PostPromise(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	promise := &entity.Promise{}
	checkRequestBody(w, r, promise)
	// 数据检查
	if promise.HappenId != user.Id && promise.HappenId != services.GetTaId(user) {
		response417Toast(w, r, "promise_nil_happen_id")
	}
	// 开始插入
	promise, err := services.AddPromise(user.Id, couple.Id, promise)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Promise *entity.Promise `json:"promise"`
	}{promise})
}

// DelPromise
func DelPromise(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	// 开始删除
	err := services.DelPromise(user.Id, couple.Id, pid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutPromise
func PutPromise(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	promise := &entity.Promise{}
	checkRequestBody(w, r, promise)
	// 开始插入
	promise, err := services.UpdatePromise(user.Id, couple.Id, promise, true)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Promise *entity.Promise `json:"promise"`
	}{promise})
}

// GetPromise
func GetPromise(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	if list {
		who, _ := strconv.Atoi(values.Get("who"))
		page, _ := strconv.Atoi(values.Get("page"))
		var suid int64
		if who == services.LIST_WHO_BY_ME {
			suid = user.Id
		} else if who == services.LIST_WHO_BY_TA {
			suid = services.GetTaId(user)
		} else {
			suid = 0
		}
		promiseList, err := services.GetPromiseListByUserCouple(user.Id, suid, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PromiseList []*entity.Promise `json:"promiseList"`
		}{promiseList})
	} else if pid > 0 {
		promise, err := services.GetPromiseById(user.Id, couple.Id, pid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Promise *entity.Promise `json:"promise"`
		}{promise})
	} else {
		response405(w, r)
	}
}
