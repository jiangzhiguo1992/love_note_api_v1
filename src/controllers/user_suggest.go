package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerSuggest(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "suggest")
	if r.Method == http.MethodPost {
		PostSuggest(w, r)
	} else if r.Method == http.MethodDelete {
		DelSuggest(w, r)
	} else if r.Method == http.MethodPut {
		PutSuggest(w, r)
	} else if r.Method == http.MethodGet {
		GetSuggest(w, r)
	} else {
		response405(w, r)
	}
}

// PostSuggest
func PostSuggest(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接收数据
	suggest := &entity.Suggest{}
	checkRequestBody(w, r, suggest)
	// 开始插入
	suggest, err := services.AddSuggest(user.Id, suggest)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Suggest *entity.Suggest `json:"suggest"`
	}{suggest})
}

// DelSuggest
func DelSuggest(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接收数据
	values := r.URL.Query()
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	// 开始删除
	err := services.DelSuggest(user.Id, sid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutSuggest
func PutSuggest(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接收数据
	suggest := &entity.Suggest{}
	checkRequestBody(w, r, suggest)
	// 开始删除
	suggest, err := services.UpdateSuggest(suggest)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Suggest *entity.Suggest `json:"suggest"`
	}{suggest})
}

// GetSuggest
func GetSuggest(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接收数据
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	follow, _ := strconv.ParseBool(values.Get("follow"))
	mine, _ := strconv.ParseBool(values.Get("mine"))
	uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
	total, _ := strconv.ParseBool(values.Get("total"))
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	if list {
		status, _ := strconv.Atoi(values.Get("status"))
		kind, _ := strconv.Atoi(values.Get("kind"))
		page, _ := strconv.Atoi(values.Get("page"))
		suggestList, err := services.GetSuggestListByStatusKind(user.Id, status, kind, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SuggestList []*entity.Suggest `json:"suggestList"`
		}{suggestList})
	} else if follow {
		page, _ := strconv.Atoi(values.Get("page"))
		suggestList, err := services.GetSuggestListByUserFollow(user.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SuggestList []*entity.Suggest `json:"suggestList"`
		}{suggestList})
	} else if mine {
		page, _ := strconv.Atoi(values.Get("page"))
		suggestList, err := services.GetSuggestListByUser(user.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SuggestList []*entity.Suggest `json:"suggestList"`
		}{suggestList})
	} else if uid > 0 {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		suggestList, err := services.GetSuggestListByUser(uid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SuggestList []*entity.Suggest `json:"suggestList"`
		}{suggestList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
		total := services.GetSuggestTotalByCreateWithDel(create)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else if sid > 0 {
		suggest, err := services.GetSuggestByIdWithAll(user.Id, sid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Suggest *entity.Suggest `json:"suggest"`
		}{suggest})
	} else {
		response405(w, r)
	}
}
