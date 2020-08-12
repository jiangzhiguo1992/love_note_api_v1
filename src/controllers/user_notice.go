package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerNotice(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "notice")
	if r.Method == http.MethodPost {
		PostNotice(w, r)
	} else if r.Method == http.MethodDelete {
		DelNotice(w, r)
	} else if r.Method == http.MethodPut {
		PutNotice(w, r)
	} else if r.Method == http.MethodGet {
		GetNotice(w, r)
	} else {
		response405(w, r)
	}
}

// PostNotice
func PostNotice(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接受参数
	notice := &entity.Notice{}
	checkRequestBody(w, r, notice)
	// 开始插入
	notice, err := services.AddNotice(notice)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_add_success")
}

// DelNotice
func DelNotice(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接收数据
	values := r.URL.Query()
	nid, _ := strconv.ParseInt(values.Get("nid"), 10, 64)
	// 开始删除
	err := services.DelNotice(nid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutNotice
func PutNotice(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	nid, _ := strconv.ParseInt(values.Get("nid"), 10, 64)
	services.AddNoticeRead(user.Id, nid)
	response200Toast(w, r, "")
}

// GetNotice
func GetNotice(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	all, _ := strconv.ParseBool(values.Get("all"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		// count
		commonCount := &services.CommonCount{
			NoticeNewCount: services.GetNoticeCountByNoRead(user.Id),
		}
		// list
		noticeList, err := services.GetNoticeListByUser(user.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			CommonCount *services.CommonCount `json:"commonCount"`
			NoticeList  []*entity.Notice      `json:"noticeList"`
		}{commonCount, noticeList})
	} else if all {
		page, _ := strconv.Atoi(values.Get("page"))
		noticeList, err := services.GetNoticeList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			NoticeList []*entity.Notice `json:"noticeList"`
		}{noticeList})
	} else {
		response405(w, r)
	}
}
