package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerDiary(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_diary")
	if r.Method == http.MethodPost {
		PostDiary(w, r)
	} else if r.Method == http.MethodDelete {
		DelDiary(w, r)
	} else if r.Method == http.MethodPut {
		PutDiary(w, r)
	} else if r.Method == http.MethodGet {
		GetDiary(w, r)
	} else {
		response405(w, r)
	}
}

// PostDiary
func PostDiary(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	diary := &entity.Diary{}
	checkRequestBody(w, r, diary)
	// 开始插入
	diary, err := services.AddDiary(user.Id, couple.Id, diary)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Diary *entity.Diary `json:"diary"`
	}{diary})
}

// DelDiary
func DelDiary(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	did, _ := strconv.ParseInt(values.Get("did"), 10, 64)
	// 开始删除
	err := services.DelDiary(user.Id, couple.Id, did)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutDiary
func PutDiary(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	diary := &entity.Diary{}
	checkRequestBody(w, r, diary)
	// 开始插入
	diary, err := services.UpdateDiary(user.Id, couple.Id, diary)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Diary *entity.Diary `json:"diary"`
	}{diary})
}

// GetDiary
func GetDiary(w http.ResponseWriter, r *http.Request) {
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
		diaryList, err := services.GetDiaryListByUserCouple(user.Id, suid, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			DiaryList []*entity.Diary `json:"diaryList"`
		}{diaryList})
	} else if did > 0 {
		diary, err := services.GetDiaryById(user.Id, couple.Id, did)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Diary *entity.Diary `json:"diary"`
		}{diary})
	} else {
		response405(w, r)
	}
}
