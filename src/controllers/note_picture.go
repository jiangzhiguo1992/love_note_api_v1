package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerPicture(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_picture")
	if r.Method == http.MethodPost {
		PostPicture(w, r)
	} else if r.Method == http.MethodDelete {
		DelPicture(w, r)
	} else if r.Method == http.MethodPut {
		PutPicture(w, r)
	} else if r.Method == http.MethodGet {
		GetPicture(w, r)
	} else {
		response405(w, r)
	}
}

// PostPicture
func PostPicture(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	body := &struct {
		PictureList []*entity.Picture `json:"pictureList"`
	}{}
	checkRequestBody(w, r, body)
	// 开始插入
	pictureList, err := services.AddPictureList(user.Id, couple.Id, body.PictureList)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		PictureList []*entity.Picture `json:"pictureList"`
	}{pictureList})
}

// DelPicture
func DelPicture(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	// 开始删除
	err := services.DelPicture(user.Id, couple.Id, pid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutPicture
func PutPicture(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	picture := &entity.Picture{}
	checkRequestBody(w, r, picture)
	// 开始插入
	picture, err := services.UpdatePicture(user.Id, couple.Id, picture)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Picture *entity.Picture `json:"picture"`
	}{picture})
}

// GetPicture
func GetPicture(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
	if aid > 0 {
		page, _ := strconv.Atoi(values.Get("page"))
		pictureList, err := services.GetPictureListByCoupleAlbum(user.Id, couple.Id, aid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			PictureList []*entity.Picture `json:"pictureList"`
		}{pictureList})
	} else {
		response405(w, r)
	}
}
