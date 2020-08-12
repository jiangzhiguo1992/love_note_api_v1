package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerMatchCoin(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_match_coin")
	if r.Method == http.MethodPost {
		PostMatchCoin(w, r)
	} else {
		response405(w, r)
	}
}

// PostMatchCoin
func PostMatchCoin(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	coin := &entity.MatchCoin{}
	checkRequestBody(w, r, coin)
	// 开始插入
	coin, err := services.AddMatchCoin(user.Id, couple.Id, coin)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		MatchCoin *entity.MatchCoin `json:"matchCoin"`
	}{coin})
}
