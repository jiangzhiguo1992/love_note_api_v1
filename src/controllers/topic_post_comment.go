package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerPostComment(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_comment")
	if r.Method == http.MethodPost {
		PostPostComment(w, r)
	} else if r.Method == http.MethodDelete {
		DelPostComment(w, r)
	} else if r.Method == http.MethodGet {
		GetPostComment(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostComment
func PostPostComment(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	couple := user.Couple
	// 接受参数
	comment := &entity.PostComment{}
	checkRequestBody(w, r, comment)
	// 开始插入
	taId := services.GetTaId(user)
	comment, err := services.AddPostComment(user.Id, taId, couple.Id, comment)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		PostComment *entity.PostComment `json:"postComment"`
	}{comment})
}

// DelPostComment
func DelPostComment(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	pcid, _ := strconv.ParseInt(values.Get("pcid"), 10, 64)
	// 开始删除
	err := services.DelPostComment(user.Id, couple.Id, pcid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetPostComment
func GetPostComment(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	cid := services.GetCoupleIdByUser(user)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	subList, _ := strconv.ParseBool(values.Get("sub_list"))
	admin, _ := strconv.ParseBool(values.Get("admin"))
	report, _ := strconv.ParseBool(values.Get("report"))
	uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
	total, _ := strconv.ParseBool(values.Get("total"))
	pcid, _ := strconv.ParseInt(values.Get("pcid"), 10, 64)
	if list {
		pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
		order, _ := strconv.Atoi(values.Get("order"))
		page, _ := strconv.Atoi(values.Get("page"))
		commentList, err := services.GetPostCommentListByPost(user.Id, cid, pid, order, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostCommentList []*entity.PostComment `json:"postCommentList"`
		}{commentList})
	} else if subList {
		pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
		tcid, _ := strconv.ParseInt(values.Get("tcid"), 10, 64)
		order, _ := strconv.Atoi(values.Get("order"))
		page, _ := strconv.Atoi(values.Get("page"))
		commentList, err := services.GetPostToCommentList(user.Id, cid, pid, tcid, order, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostCommentList []*entity.PostComment `json:"postCommentList"`
		}{commentList})
	} else if admin {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
		tcid, _ := strconv.ParseInt(values.Get("tcid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		commentList, err := services.GetPostCommentList(uid, pid, tcid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostCommentList []*entity.PostComment `json:"postCommentList"`
		}{commentList})
	} else if report {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		commentList, err := services.GetPostCommentReportList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostCommentList []*entity.PostComment `json:"postCommentList"`
		}{commentList})
	} else if uid > 0 { // 要在admin后面
		pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
		order, _ := strconv.Atoi(values.Get("order"))
		page, _ := strconv.Atoi(values.Get("page"))
		commentList, err := services.GetPostCommentListByUserPost(user.Id, cid, pid, uid, order, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostCommentList []*entity.PostComment `json:"postCommentList"`
		}{commentList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
		total := services.GetPostCommentTotalByCreateWithDel(create)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else if pcid > 0 {
		comment, err := services.GetPostCommentByIdWithAll(user.Id, cid, pcid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PostComment *entity.PostComment `json:"postComment"`
		}{comment})
	} else {
		response405(w, r)
	}
}
