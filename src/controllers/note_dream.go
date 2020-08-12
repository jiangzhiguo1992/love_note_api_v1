package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerDream(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_dream")
	if r.Method == http.MethodPost {
		PostDream(w, r)
	} else if r.Method == http.MethodDelete {
		DelDream(w, r)
	} else if r.Method == http.MethodPut {
		PutDream(w, r)
	} else if r.Method == http.MethodGet {
		GetDream(w, r)
	} else {
		response405(w, r)
	}
}

// PostDream
func PostDream(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	dream := &entity.Dream{}
	checkRequestBody(w, r, dream)
	// 开始插入
	dream, err := services.AddDream(user.Id, couple.Id, dream)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Dream *entity.Dream `json:"dream"`
	}{dream})
}

// DelDream
func DelDream(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	did, _ := strconv.ParseInt(values.Get("did"), 10, 64)
	// 开始删除
	err := services.DelDream(user.Id, couple.Id, did)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutDream
func PutDream(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	dream := &entity.Dream{}
	checkRequestBody(w, r, dream)
	// 开始插入
	dream, err := services.UpdateDream(user.Id, couple.Id, dream)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Dream *entity.Dream `json:"dream"`
	}{dream})
}

// GetDream
func GetDream(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	did, _ := strconv.ParseInt(values.Get("did"), 10, 64)
	if list {
		who, _ := strconv.Atoi(values.Get("who"))
		page, _ := strconv.Atoi(values.Get("page"))
		var suid int64
		if who == services.LIST_WHO_BY_ME {
			suid = user.Id
		} else if who == services.LIST_WHO_BY_TA {
			suid = services.GetTaId(user)
		} else {
			suid = 0
		}
		dreamList, err := services.GetDreamListByUserCouple(user.Id, suid, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			DreamList []*entity.Dream `json:"dreamList"`
		}{dreamList})
	} else if did > 0 {
		dream, err := services.GetDreamById(user.Id, couple.Id, did)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Dream *entity.Dream `json:"dream"`
		}{dream})
	} else {
		response405(w, r)
	}
}
