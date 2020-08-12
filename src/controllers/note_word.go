package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerWord(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_word")
	if r.Method == http.MethodPost {
		PostWord(w, r)
	} else if r.Method == http.MethodDelete {
		DelWord(w, r)
	} else if r.Method == http.MethodGet {
		GetWord(w, r)
	} else {
		response405(w, r)
	}
}

// PostWord
func PostWord(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	word := &entity.Word{}
	checkRequestBody(w, r, word)
	// 开始插入
	word, err := services.AddWord(user.Id, couple.Id, word)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Word *entity.Word `json:"word"`
	}{word})
}

// DelWord
func DelWord(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	wid, _ := strconv.ParseInt(values.Get("wid"), 10, 64)
	// 开始删除
	err := services.DelWord(user.Id, couple.Id, wid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "")
}

// GetWord
func GetWord(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		// 查询
		wordList, err := services.GetWordListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			WordList []*entity.Word `json:"wordList"`
		}{wordList})
	} else {
		response405(w, r)
	}
}
