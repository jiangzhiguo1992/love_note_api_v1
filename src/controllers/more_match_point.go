package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerMatchPoint(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_match_point")
	if r.Method == http.MethodPost {
		PostMatchPoint(w, r)
	} else {
		response405(w, r)
	}
}

// PostMatchPoint
func PostMatchPoint(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	point := &entity.MatchPoint{}
	checkRequestBody(w, r, point)
	// 开始插入
	point, err := services.ToggleMatchPoint(user.Id, couple.Id, point)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		MatchPoint *entity.MatchPoint `json:"matchPoint"`
	}{point})
}
