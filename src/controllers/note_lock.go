package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerLock(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_lock")
	if r.Method == http.MethodPost {
		PostLock(w, r)
	} else if r.Method == http.MethodPut {
		PutLock(w, r)
	} else if r.Method == http.MethodGet {
		GetLock(w, r)
	} else {
		response405(w, r)
	}
}

// PostLock
func PostLock(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	lock := &entity.Lock{}
	checkRequestBody(w, r, lock)
	// 开始插入
	lock, err := services.AddLock(user.Id, couple.Id, lock)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Lock *entity.Lock `json:"lock"`
	}{lock})
}

// PutLock
func PutLock(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	modify, _ := strconv.ParseBool(values.Get("modify"))
	toggle, _ := strconv.ParseBool(values.Get("toggle"))
	body := &entity.Lock{}
	checkRequestBody(w, r, body)
	// 旧数据检查
	lock, err := services.GetLockByUserCouple(user.Id, couple.Id)
	response417ErrDialog(w, r, err)
	if lock == nil {
		response417Dialog(w, r, "nil_lock")
	}
	toast := ""
	if modify {
		// 修改密码
		code := values.Get("code")
		// sms检查
		taId := services.GetTaId(user)
		ta, _ := services.GetUserById(taId)
		if ta == nil {
			response417Dialog(w, r, "nil_user")
		}
		err := services.SmsCheckCode(ta.Phone, entity.SMS_TYPE_LOCK, code)
		response417ErrDialog(w, r, err)
		// 开始插入
		lock.Password = body.Password
		lock, err = services.UpdateLock(lock)
		response417ErrDialog(w, r, err)
		toast = "db_update_success"
	} else if toggle {
		// 开关锁
		if lock.IsLock && body.Password != lock.Password {
			response417Dialog(w, r, "user_pwd_wrong")
		}
		lock.IsLock = !lock.IsLock
		lock, err = services.UpdateLock(lock)
		response417ErrDialog(w, r, err)
	} else {
		response405(w, r)
	}
	// 返回
	response200DataToast(w, r, toast, struct {
		Lock *entity.Lock `json:"lock"`
	}{lock})
}

// GetLock
func GetLock(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	lock, err := services.GetLockByUserCouple(user.Id, couple.Id)
	response417ErrDialog(w, r, err) // dialog
	if lock != nil {
		lock.Password = ""
	}
	// 返回
	response200Data(w, r, struct {
		Lock *entity.Lock `json:"lock"`
	}{lock})
}
