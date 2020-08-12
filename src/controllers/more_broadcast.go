package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerBroadcast(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_broadcast")
	if r.Method == http.MethodPost {
		PostBroadcast(w, r)
	} else if r.Method == http.MethodDelete {
		DelBroadcast(w, r)
	} else if r.Method == http.MethodGet {
		GetBroadcast(w, r)
	} else {
		response405(w, r)
	}
}

// PostBroadcast
func PostBroadcast(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接受参数
	broadcast := &entity.Broadcast{}
	checkRequestBody(w, r, broadcast)
	// 开始插入
	broadcast, err := services.AddBroadcast(broadcast)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Broadcast *entity.Broadcast `json:"broadcast"`
	}{broadcast})
}

// DelBroadcast
func DelBroadcast(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接收数据
	values := r.URL.Query()
	bid, _ := strconv.ParseInt(values.Get("bid"), 10, 64)
	// 开始删除
	err := services.DelBroadcast(bid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetBroadcast
func GetBroadcast(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		broadcastList, err := services.GetBroadcastList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			BroadcastList []*entity.Broadcast `json:"broadcastList"`
		}{broadcastList})
	} else {
		response405(w, r)
	}
}
