package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerMensesDay(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_menses")
	if r.Method == http.MethodPost {
		PostMensesDay(w, r)
	} else {
		response405(w, r)
	}
}

// PostMensesDay 添加
func PostMensesDay(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	mensesDay := &entity.MensesDay{}
	checkRequestBody(w, r, mensesDay)
	// 开始插入
	taId := services.GetTaId(user)
	mensesDay, err := services.AddMensesDay(user.Id, taId, couple.Id, mensesDay)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		MensesDay *entity.MensesDay `json:"mensesDay"`
	}{mensesDay})
}
