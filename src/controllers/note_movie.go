package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerMovie(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_movie")
	if r.Method == http.MethodPost {
		PostMovie(w, r)
	} else if r.Method == http.MethodDelete {
		DelMovie(w, r)
	} else if r.Method == http.MethodPut {
		PutMovie(w, r)
	} else if r.Method == http.MethodGet {
		GetMovie(w, r)
	} else {
		response405(w, r)
	}
}

// PostMovie
func PostMovie(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	movie := &entity.Movie{}
	checkRequestBody(w, r, movie)
	// 开始插入
	movie, err := services.AddMovie(user.Id, couple.Id, movie)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Movie *entity.Movie `json:"movie"`
	}{movie})
}

// DelMovie
func DelMovie(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	mid, _ := strconv.ParseInt(values.Get("mid"), 10, 64)
	// 开始删除
	err := services.DelMovie(user.Id, couple.Id, mid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutMovie
func PutMovie(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	movie := &entity.Movie{}
	checkRequestBody(w, r, movie)
	// 开始插入
	movie, err := services.UpdateMovie(user.Id, couple.Id, movie)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Movie *entity.Movie `json:"movie"`
	}{movie})
}

// GetMovie
func GetMovie(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		movieList, err := services.GetMovieListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			MovieList []*entity.Movie `json:"movieList"`
		}{movieList})
	} else {
		response405(w, r)
	}
}
