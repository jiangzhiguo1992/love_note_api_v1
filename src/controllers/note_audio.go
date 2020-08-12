package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerAudio(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_audio")
	if r.Method == http.MethodPost {
		PostAudio(w, r)
	} else if r.Method == http.MethodDelete {
		DelAudio(w, r)
	} else if r.Method == http.MethodGet {
		GetAudio(w, r)
	} else {
		response405(w, r)
	}
}

// PostAudio
func PostAudio(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	audio := &entity.Audio{}
	checkRequestBody(w, r, audio)
	// 开始插入
	audio, err := services.AddAudio(user.Id, couple.Id, audio)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Audio *entity.Audio `json:"audio"`
	}{audio})
}

// DelAudio
func DelAudio(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
	// 开始删除
	err := services.DelAudio(user.Id, couple.Id, aid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetAudio
func GetAudio(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		audioList, err := services.GetAudioListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			AudioList []*entity.Audio `json:"audioList"`
		}{audioList})
	} else {
		response405(w, r)
	}
}
