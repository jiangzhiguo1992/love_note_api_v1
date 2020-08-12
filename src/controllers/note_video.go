package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerVideo(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_video")
	if r.Method == http.MethodPost {
		PostVideo(w, r)
	} else if r.Method == http.MethodDelete {
		DelVideo(w, r)
	} else if r.Method == http.MethodPut {
		PutVideo(w, r)
	} else if r.Method == http.MethodGet {
		GetVideo(w, r)
	} else {
		response405(w, r)
	}
}

// PostVideo
func PostVideo(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	video := &entity.Video{}
	checkRequestBody(w, r, video)
	// 开始插入
	video, err := services.AddVideo(user.Id, couple.Id, video)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Video *entity.Video `json:"video"`
	}{video})
}

// DelVideo
func DelVideo(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	vid, _ := strconv.ParseInt(values.Get("vid"), 10, 64)
	// 开始删除
	err := services.DelVideo(user.Id, couple.Id, vid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutVideo
func PutVideo(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	video := &entity.Video{}
	checkRequestBody(w, r, video)
	// 开始插入
	video, err := services.UpdateVideo(user.Id, couple.Id, video)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Video *entity.Video `json:"video"`
	}{video})
}

// GetVideo
func GetVideo(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		videoList, err := services.GetVideoListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			VideoList []*entity.Video `json:"videoList"`
		}{videoList})
	} else {
		response405(w, r)
	}
}
