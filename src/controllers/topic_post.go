package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post")
	if r.Method == http.MethodPost {
		PostPost(w, r)
	} else if r.Method == http.MethodDelete {
		DelPost(w, r)
	} else if r.Method == http.MethodPut {
		PutPost(w, r)
	} else if r.Method == http.MethodGet {
		GetPost(w, r)
	} else {
		response405(w, r)
	}
}

// PostPost
func PostPost(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	couple := user.Couple
	// 接受参数
	post := &entity.Post{}
	checkRequestBody(w, r, post)
	// 开始插入
	post, err := services.AddPost(user.Id, couple.Id, post)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Post *entity.Post `json:"post"`
	}{post})
}

// DelPost
func DelPost(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	// 开始删除
	err := services.DelPost(user.Id, couple.Id, pid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutPost
func PutPost(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接收数据
	post := &entity.Post{}
	checkRequestBody(w, r, post)
	// 开始删除
	post, err := services.UpdatePost(post)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Post *entity.Post `json:"post"`
	}{post})
}

// GetPost
func GetPost(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	cid := services.GetCoupleIdByUser(user)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	collect, _ := strconv.ParseBool(values.Get("collect"))
	mine, _ := strconv.ParseBool(values.Get("mine"))
	report, _ := strconv.ParseBool(values.Get("report"))
	admin, _ := strconv.ParseBool(values.Get("admin"))
	total, _ := strconv.ParseBool(values.Get("total"))
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	if list {
		search := values.Get("search")
		create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
		kind, _ := strconv.Atoi(values.Get("kind"))
		subKind, _ := strconv.Atoi(values.Get("sub_kind"))
		official, _ := strconv.ParseBool(values.Get("official"))
		well, _ := strconv.ParseBool(values.Get("well"))
		page, _ := strconv.Atoi(values.Get("page"))
		var postList []*entity.Post
		var err error
		if len(search) > 0 {
			postList, err = services.GetPostListBySearch(user.Id, cid, search, page)
		} else if kind == 0 {
			postList, err = services.GetPostListByCreate(user.Id, cid, create, page)
		} else {
			postList, err = services.GetPostListByCreateKindOfficialWell(user.Id, cid, create, kind, subKind, official, well, page)
		}
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostList []*entity.Post `json:"postList"`
		}{postList})
	} else if collect {
		if cid <= 0 {
			response417NoCP(w, r)
		}
		me, _ := strconv.ParseBool(values.Get("me"))
		page, _ := strconv.Atoi(values.Get("page"))
		searchId := user.Id
		if !me {
			searchId = services.GetTaId(user)
		}
		postList, err := services.GetPostListByUserCoupleCollect(user.Id, searchId, cid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostList []*entity.Post `json:"postList"`
		}{postList})
	} else if mine {
		if cid <= 0 {
			response417NoCP(w, r)
		}
		page, _ := strconv.Atoi(values.Get("page"))
		postList, err := services.GetPostListByUserCouple(user.Id, cid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostList []*entity.Post `json:"postList"`
		}{postList})
	} else if admin {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		postList, err := services.GetPostList(uid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostList []*entity.Post `json:"postList"`
		}{postList})
	} else if report {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		postList, err := services.GetPostReportList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostList []*entity.Post `json:"postList"`
		}{postList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
		total := services.GetPostTotalByCreateWithDel(create)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else if pid > 0 {
		post, err := services.GetPostByIdWithAll(user.Id, cid, pid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Post *entity.Post `json:"post"`
		}{post})
	} else {
		response405(w, r)
	}
}
