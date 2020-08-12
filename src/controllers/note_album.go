package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerAlbum(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_album")
	if r.Method == http.MethodPost {
		PostAlbum(w, r)
	} else if r.Method == http.MethodDelete {
		DelAlbum(w, r)
	} else if r.Method == http.MethodPut {
		PutAlbum(w, r)
	} else if r.Method == http.MethodGet {
		GetAlbum(w, r)
	} else {
		response405(w, r)
	}
}

// PostAlbum
func PostAlbum(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	album := &entity.Album{}
	checkRequestBody(w, r, album)
	// 开始插入
	album, err := services.AddAlbum(user.Id, couple.Id, album)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Album *entity.Album `json:"album"`
	}{album})
}

// DelAlbum
func DelAlbum(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
	// 开始删除
	err := services.DelAlbum(user.Id, couple.Id, aid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutAlbum
func PutAlbum(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	album := &entity.Album{}
	checkRequestBody(w, r, album)
	// 开始插入
	album, err := services.UpdateAlbum(user.Id, couple.Id, album, true)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Album *entity.Album `json:"album"`
	}{album})
}

// GetAlbum
func GetAlbum(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		albumList, err := services.GetAlbumListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			AlbumList []*entity.Album `json:"albumList"`
		}{albumList})
	} else if aid > 0 {
		album, err := services.GetAlbumById(user.Id, couple.Id, aid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Album *entity.Album `json:"album"`
		}{album})
	} else {
		response405(w, r)
	}
}
