package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerWhisper(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_whisper")
	if r.Method == http.MethodPost {
		PostWhisper(w, r)
	} else if r.Method == http.MethodGet {
		GetWhisper(w, r)
	} else {
		response405(w, r)
	}
}

// PostWhisper
func PostWhisper(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	whisper := &entity.Whisper{}
	checkRequestBody(w, r, whisper)
	// 开始插入
	whisper, err := services.AddWhisper(user.Id, couple.Id, whisper)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Whisper *entity.Whisper `json:"whisper"`
	}{whisper})
}

// GetWhisper
func GetWhisper(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		channel := values.Get("channel")
		page, _ := strconv.Atoi(values.Get("page"))
		whisperList, err := services.GetWhisperListByCoupleChannel(user.Id, couple.Id, channel, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			WhisperList []*entity.Whisper `json:"whisperList"`
		}{whisperList})
	} else {
		response405(w, r)
	}
}
