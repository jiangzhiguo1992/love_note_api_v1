package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"libs/utils"
	"models/entity"
	"services"
)

const (
	// context传递的key
	CONTEXT_KEY_USER = "token_user"
	CONTEXT_KEY_BODY = "request_body"
	// head-params
	HEAD_PARAMS_ADMIN_KEY = "adminKey"
	HEAD_PARAMS_APP_KEY1  = "appKey"
	HEAD_PARAMS_APP_KEY2  = "app_key"
	HEAD_PARAMS_APP_KEY3  = "App_key"
	HEAD_PARAMS_TOKEN1    = "accessToken"
	HEAD_PARAMS_TOKEN2    = "access_token"
	HEAD_PARAMS_TOKEN3    = "Access_token"
	HEAD_PARAMS_PLATFORM1 = "platform"
	HEAD_PARAMS_PLATFORM2 = "Platform"
	HEAD_PARAMS_LANGUAGE1 = "language"
	HEAD_PARAMS_LANGUAGE2 = "Language"
)

// checkRequestHeader 检查head
func checkRequestHeader(w http.ResponseWriter, r *http.Request) {
	adminKey := utils.GetConfigStr("conf", "app.conf", "phone", "admin_key")
	appKey := utils.GetConfigStr("conf", "app.conf", "phone", "app_key")
	// 验证admin-key
	if strings.TrimSpace(r.Header.Get(HEAD_PARAMS_ADMIN_KEY)) == adminKey {
		return // admin不验证后面的
	}
	// 验证app-key 简单防止api攻击
	var keyPush1 = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_APP_KEY1))
	var keyPush2 = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_APP_KEY2))
	var keyPush3 = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_APP_KEY3))
	if appKey != keyPush1 && appKey != keyPush2 && appKey != keyPush3 {
		// utils.LogWarn("checkRequestHeader", r.Header)
		response403(w, r)
	}
}

// 模块开关
func checkModelStatus(w http.ResponseWriter, r *http.Request, model string) {
	if !getModelStatus(model) {
		responseModelClose(w, r)
	}
}

func getModelStatus(model string) bool {
	return utils.GetConfigBool("conf", "model.conf", "model", model)
}

/*****************************************************************************************************/
/*****************************************************************************************************/

// checkTokenUser 只有在router中和router中没有检查user的才会用
func checkTokenUser(w http.ResponseWriter, r *http.Request) *entity.User {
	token := strings.TrimSpace(r.Header.Get(HEAD_PARAMS_TOKEN1))
	if len(token) <= 0 {
		token = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_TOKEN2))
	}
	if len(token) <= 0 {
		token = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_TOKEN3))
	}
	if len(token) <= 0 {
		response401(w, r)
	}
	user, _ := services.GetUserByToken(token)
	if user == nil || user.Id <= 0 {
		response401(w, r)
	} else if services.IsUserBlack(user) {
		response406(w, r)
	} else if !services.IsUserInfoComplete(user) {
		response417NoUserInfo(w, r, user)
	}
	return user
}

// checkTokenCouple 检验cp
func checkTokenCouple(w http.ResponseWriter, r *http.Request) *entity.User {
	user := checkTokenUser(w, r)
	if user.Couple != nil && user.Couple.Id > 0 {
		return user
	}
	// couple
	couple, _ := services.GetCoupleVisibleByUser(user.Id)
	if couple == nil || couple.Id <= 0 {
		response417NoCP(w, r)
	}
	user.Couple = couple
	return user
}

// getTokenUser 从router的检验结果中获取user
func getTokenUser(r *http.Request) *entity.User {
	c := r.Context()
	value := c.Value(CONTEXT_KEY_USER)
	if value == nil {
		return &entity.User{}
	} else {
		return value.(*entity.User)
	}
}

// getTokenCouple
func getTokenCouple(r *http.Request) (*entity.User, error) {
	user := getTokenUser(r)
	if user.Couple != nil && user.Couple.Id > 0 {
		return user, nil
	}
	// couple
	couple, _ := services.GetCoupleVisibleByUser(user.Id)
	if couple == nil || couple.Id <= 0 {
		return user, errors.New("couple_need")
	}
	user.Couple = couple
	return user, nil
}

// checkRequestBody 检查传输来的数据
func checkRequestBody(w http.ResponseWriter, r *http.Request, v interface{}) {
	bytes := getRequestBody(r)
	if bytes == nil || len(bytes) <= 0 {
		response417Toast(w, r, "request_body_nil")
	}
	err := json.Unmarshal(bytes, v)
	if v == nil || err != nil {
		response417Toast(w, r, "data_decode_err")
	}
}

// getRequestBody 获取传输来的数据
func getRequestBody(r *http.Request) []byte {
	c := r.Context()
	value := c.Value(CONTEXT_KEY_BODY)
	if value == nil {
		return make([]byte, 0)
	} else {
		return value.([]byte)
	}
}

// getUrlSuffixInt
func getUrlSuffixInt(r *http.Request, prefix string, index int) (int, error) {
	suffixes, err := getUrlSuffixes(r, prefix)
	if err != nil {
		return 0, err
	}
	suf := suffixes[index]
	suf32, err := strconv.Atoi(suf)
	return suf32, err
}

// getUrlSuffixInt64
func getUrlSuffixInt64(r *http.Request, prefix string, index int) (int64, error) {
	suffixes, err := getUrlSuffixes(r, prefix)
	if err != nil {
		return 0, err
	}
	suf := suffixes[index]
	suf64, err := strconv.ParseInt(suf, 10, 64)
	return suf64, err
}

// getUrlSuffixes 获取路径中最后的/后面的内容
func getUrlSuffixes(r *http.Request, prefix string) ([]string, error) {
	path := r.URL.Path
	trimPrefix := strings.TrimPrefix(strings.TrimSpace(path), strings.TrimSpace(prefix))
	split := strings.Split(strings.TrimSpace(trimPrefix), "/")
	suffix := make([]string, 0)
	for _, v := range split {
		s := strings.TrimSpace(v)
		if len(s) > 0 {
			suffix = append(suffix, s)
		}
	}
	//if len(suffix) <= 0 {
	//	return suffix, errors.New("request_url_nil")
	//}
	return suffix, nil
}
