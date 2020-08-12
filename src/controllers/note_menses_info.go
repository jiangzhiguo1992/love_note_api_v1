package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerMensesInfo(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_menses")
	if r.Method == http.MethodPut {
		PutMensesInfo(w, r)
	} else if r.Method == http.MethodGet {
		GetMensesInfo(w, r)
	} else {
		response405(w, r)
	}
}

// PutMensesInfo 发表
func PutMensesInfo(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	mensesInfo := &services.MensesInfo{}
	checkRequestBody(w, r, mensesInfo)
	// 我的
	if user.Sex == entity.USER_SEX_GIRL {
		lengthMe := mensesInfo.MensesLengthMe
		mensesLength, err := services.AddMensesLength(user.Id, couple.Id, lengthMe)
		response417ErrToast(w, r, err)
		mensesInfo.MensesLengthMe = mensesLength
	}
	// ta的
	taId := services.GetTaId(user)
	ta, _ := services.GetUserById(taId)
	if ta == nil {
		response200Show(w, r, "nil_user")
	} else if ta.Sex == entity.USER_SEX_GIRL {
		lengthTa := mensesInfo.MensesLengthTa
		mensesLength, err := services.AddMensesLength(ta.Id, couple.Id, lengthTa)
		response417ErrToast(w, r, err)
		mensesInfo.MensesLengthTa = mensesLength
	}
	// 返回
	response200Data(w, r, struct {
		MensesInfo *services.MensesInfo `json:"mensesInfo"`
	}{mensesInfo})
}

// GetMensesInfo 获取
func GetMensesInfo(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 默认
	mensesInfo := &services.MensesInfo{
		CanMe:          false,
		CanTa:          false,
		MensesLengthMe: nil,
		MensesLengthTa: nil,
	}
	// 我的
	if user.Sex == entity.USER_SEX_GIRL {
		mensesInfo.CanMe = true
		mensesInfo.MensesLengthMe, _ = services.GetMensesLengthByUserCouple(user.Id, couple.Id)
	}
	// ta的
	taId := services.GetTaId(user)
	ta, _ := services.GetUserById(taId)
	if ta == nil {
		response200Show(w, r, "nil_user")
	} else if ta.Sex == entity.USER_SEX_GIRL {
		mensesInfo.CanTa = true
		mensesInfo.MensesLengthTa, _ = services.GetMensesLengthByUserCouple(ta.Id, couple.Id)
	}
	// 返回
	response200Data(w, r, struct {
		MensesInfo *services.MensesInfo `json:"mensesInfo"`
	}{mensesInfo})

}
