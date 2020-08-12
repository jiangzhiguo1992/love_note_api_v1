package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "user_login")
	if r.Method == http.MethodPost {
		PostUserLogin(w, r)
	} else if r.Method == http.MethodGet {
		GetUserLogin(w, r)
	} else {
		response405(w, r)
	}
}

// PostUserLogin 密码/验证码登录
func PostUserLogin(w http.ResponseWriter, r *http.Request) {
	// 参数
	values := r.URL.Query()
	paramsType, _ := strconv.Atoi(values.Get("type"))
	code := values.Get("code")
	paramsUser := &entity.User{}
	checkRequestBody(w, r, paramsUser)
	// 查找user，这里的err不检查，后续会用到
	dbUser, err := services.GetUserByPhone(paramsUser.Phone)
	response417ErrDialog(w, r, err)
	switch paramsType {
	case services.USER_TYPE_LOG_PWD: // 密码登录，需要先注册，才能登录
		if dbUser == nil {
			response417Dialog(w, r, "user_phone_no_exist")
		}
		// 密码校验
		if len([]rune(dbUser.Password)) <= 0 {
			response417Dialog(w, r, "user_pwd_no_set")
		} else if paramsUser.Password != dbUser.Password {
			response417Dialog(w, r, "user_pwd_wrong")
		}
	case services.USER_TYPE_LOG_VER: // 验证码登录，可直接跳过注册
		err = services.SmsCheckCode(paramsUser.Phone, entity.SMS_TYPE_LOGIN, code)
		response417ErrDialog(w, r, err)
		if dbUser == nil {
			// 没注册过，则顺便注册
			dbUser, err = services.AddUserByVer(paramsUser)
			response417ErrDialog(w, r, err)
		}
	default:
		response405(w, r)
	}
	// 开始登录
	dbUser, err = services.UpdateUserOnLogin(dbUser)
	response417ErrDialog(w, r, err)
	// 装载信息
	dbUser.Couple, _ = services.GetCoupleVisibleByUser(dbUser.Id)
	dbUser.Password = ""
	// 返回
	if dbUser.Status <= entity.STATUS_DELETE {
		response406(w, r)
	} else if !services.IsUserInfoComplete(dbUser) {
		response417NoUserInfo(w, r, dbUser)
	} else {
		response200DataToast(w, r, "user_login_success", struct {
			User *entity.User `json:"user"`
		}{dbUser})
	}
}

// GetUserLogin
func GetUserLogin(w http.ResponseWriter, r *http.Request) {
	// 接受参数
	values := r.URL.Query()
	phone := values.Get("phone")
	pwd := values.Get("pwd")
	// 查找user，这里的err不检查，后续会用到
	dbUser, err := services.GetUserByPhone(phone)
	response417ErrDialog(w, r, err)
	if dbUser == nil {
		response417Dialog(w, r, "user_phone_no_exist")
	} else if len([]rune(dbUser.Password)) <= 0 {
		response417Dialog(w, r, "user_pwd_no_set")
	} else if pwd != dbUser.Password {
		response417Dialog(w, r, "user_pwd_wrong")
	}
	dbUser.Password = ""
	// admin
	if !services.IsAdminister(dbUser) {
		response417Toast(w, r, "")
	}
	// 返回
	response200Data(w, r, struct {
		User *entity.User `json:"user"`
	}{dbUser})
}
