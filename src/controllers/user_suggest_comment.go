package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerSuggestComment(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "suggest_comment")
	if r.Method == http.MethodPost {
		PostSuggestComment(w, r)
	} else if r.Method == http.MethodDelete {
		DelSuggestComment(w, r)
	} else if r.Method == http.MethodGet {
		GetSuggestComment(w, r)
	} else {
		response405(w, r)
	}
}

// PostSuggestComment
func PostSuggestComment(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	comment := &entity.SuggestComment{}
	checkRequestBody(w, r, comment)
	// 开始上传
	comment, err := services.AddSuggestComment(user.Id, comment)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Comment *entity.SuggestComment `json:"comment"`
	}{comment})
}

// DelSuggestComment
func DelSuggestComment(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接收数据
	values := r.URL.Query()
	scid, _ := strconv.ParseInt(values.Get("scid"), 10, 64)
	// 开始删除
	err := services.DelSuggestComment(user.Id, scid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetSuggestComment
func GetSuggestComment(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	admin, _ := strconv.ParseBool(values.Get("admin"))
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	total, _ := strconv.ParseBool(values.Get("total"))
	if admin {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		commentList, err := services.GetSuggestCommentList(uid, sid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SuggestCommentList []*entity.SuggestComment `json:"suggestCommentList"`
		}{commentList})
	} else if sid > 0 { // 要在admin后面
		page, _ := strconv.Atoi(values.Get("page"))
		// 获取列表
		commentList, err := services.GetSuggestCommentListWithAll(user.Id, sid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SuggestCommentList []*entity.SuggestComment `json:"suggestCommentList"`
		}{commentList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
		total := services.GetSuggestCommentTotalByCreate(create)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
