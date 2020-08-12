package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerAngry(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_angry")
	if r.Method == http.MethodPost {
		PostAngry(w, r)
	} else if r.Method == http.MethodDelete {
		DelAngry(w, r)
	} else if r.Method == http.MethodPut {
		PutAngry(w, r)
	} else if r.Method == http.MethodGet {
		GetAngry(w, r)
	} else {
		response405(w, r)
	}
}

// PostAngry
func PostAngry(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	angry := &entity.Angry{}
	checkRequestBody(w, r, angry)
	// 数据检查
	if angry.HappenId != user.Id && angry.HappenId != services.GetTaId(user) {
		response417Toast(w, r, "angry_nil_happen_id")
	}
	// 开始插入
	angry, err := services.AddAngry(user.Id, couple.Id, angry)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Angry *entity.Angry `json:"angry"`
	}{angry})
}

// DelAngry
func DelAngry(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
	// 开始删除
	err := services.DelAngry(user.Id, couple.Id, aid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutAngry
func PutAngry(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	angry := &entity.Angry{}
	checkRequestBody(w, r, angry)
	// 开始插入
	angry, err := services.UpdateAngry(user.Id, couple.Id, angry)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Angry *entity.Angry `json:"angry"`
	}{angry})
}

// GetAngry
func GetAngry(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
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
		angryList, err := services.GetAngryListByUserCouple(user.Id, suid, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			AngryList []*entity.Angry `json:"angryList"`
		}{angryList})
	} else if aid > 0 {
		angry, err := services.GetAngryByIdWithGiftPromise(user.Id, couple.Id, aid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Angry *entity.Angry `json:"angry"`
		}{angry})
	} else {
		response405(w, r)
	}
}
