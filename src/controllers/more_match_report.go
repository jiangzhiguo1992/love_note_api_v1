package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerMatchReport(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_match_report")
	if r.Method == http.MethodPost {
		PostMatchReport(w, r)
	} else {
		response405(w, r)
	}
}

// PostMatchReport
func PostMatchReport(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	report := &entity.MatchReport{}
	checkRequestBody(w, r, report)
	// 开始插入
	report, err := services.AddMatchReport(user.Id, couple.Id, report)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		MatchReport *entity.MatchReport `json:"matchReport"`
	}{report})
}
