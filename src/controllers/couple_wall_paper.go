package controllers

import (
	"net/http"

	"models/entity"
	"services"
)

func HandlerWallPaper(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "couple_wall")
	if r.Method == http.MethodPost {
		PostWallPaper(w, r)
	} else if r.Method == http.MethodGet {
		GetWallPaper(w, r)
	} else {
		response405(w, r)
	}
}

// PostWallPaper
func PostWallPaper(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接收数据
	wallPaper := &entity.WallPaper{}
	checkRequestBody(w, r, wallPaper)
	wallPaper.CoupleId = couple.Id
	// 开始增加/修改
	wallPaper, err := services.UpdateWallPaper(wallPaper)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		WallPaper *entity.WallPaper `json:"wallPaper"`
	}{wallPaper})
}

// GetWallPaper
func GetWallPaper(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 获取bg
	wallPaper, err := services.GetWallPaperByCouple(couple.Id)
	response200ErrShow(w, r, err)
	// 返回
	response200Data(w, r, struct {
		WallPaper *entity.WallPaper `json:"wallPaper"`
	}{wallPaper})
}
