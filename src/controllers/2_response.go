package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"libs/utils"
	"models/entity"
)

type (
	// 返回body
	HttpResult struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	// resultCode 一般为417时使用
	RESULT_CODE_OK           = 0
	RESULT_CODE_TOAST        = 10
	RESULT_CODE_DIALOG       = 11
	RESULT_CODE_NO_USER_INFO = 20
	RESULT_CODE_NO_CP        = 30
)

func responseModelClose(w http.ResponseWriter, r *http.Request) {
	response417Dialog(w, r, "request_model_close")
}

// response200 正常
func response200ErrShow(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		response200Show(w, r, err.Error())
	}
}

func response200Show(w http.ResponseWriter, r *http.Request, msg string) {
	language := strings.ToLower(strings.TrimSpace(r.Header.Get(HEAD_PARAMS_LANGUAGE1)))
	if len(language) <= 0 {
		language = strings.ToLower(strings.TrimSpace(r.Header.Get(HEAD_PARAMS_LANGUAGE2)))
	}
	message := utils.GetLanguage(language, msg)
	result := &HttpResult{Code: RESULT_CODE_OK, Data: struct {
		Show string `json:"show"`
	}{message}}
	response200(w, r, result)
}

func response200DataToast(w http.ResponseWriter, r *http.Request, msg string, data interface{}) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: msg, Data: data}
	response200(w, r, result)
}

func response200Toast(w http.ResponseWriter, r *http.Request, msg string) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: msg}
	response200(w, r, result)
}

func response200Data(w http.ResponseWriter, r *http.Request, data interface{}) {
	result := &HttpResult{Code: RESULT_CODE_OK, Data: data}
	response200(w, r, result)
}

func response200(w http.ResponseWriter, r *http.Request, result *HttpResult) {
	responseJson(w, r, http.StatusOK, result)
}

// response401 无用户令牌
func response401(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: "http_status_401"}
	responseJson(w, r, http.StatusUnauthorized, result)
}

// response403 key不对
func response403(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: "http_status_403"}
	responseJson(w, r, http.StatusForbidden, result)
}

// Response404 路径错误
func Response404(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: "http_status_404"}
	responseJson(w, r, http.StatusNotFound, result)
}

// response405 方法错误
func response405(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: "http_status_405"}
	responseJson(w, r, http.StatusMethodNotAllowed, result)
}

// response406 不允许进入
func response406(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: "http_status_406"}
	responseJson(w, r, http.StatusNotAcceptable, result)
}

// response409 app新版本
func response409(w http.ResponseWriter, r *http.Request, msg string, data interface{}) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: msg, Data: data}
	responseJson(w, r, http.StatusConflict, result)
}

// response417 逻辑错误
func response417ErrToast(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		response417Toast(w, r, err.Error())
	}
}

func response417ErrDialog(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		response417Dialog(w, r, err.Error())
	}
}

func response417Toast(w http.ResponseWriter, r *http.Request, msg string) {
	result := &HttpResult{Code: RESULT_CODE_TOAST, Message: msg}
	response417(w, r, result)
}

func response417Dialog(w http.ResponseWriter, r *http.Request, msg string) {
	result := &HttpResult{Code: RESULT_CODE_DIALOG, Message: msg}
	response417(w, r, result)
}

func response417NoUserInfo(w http.ResponseWriter, r *http.Request, user *entity.User) {
	result := &HttpResult{Code: RESULT_CODE_NO_USER_INFO, Message: "user_info_nil", Data: struct {
		User *entity.User `json:"user"`
	}{user}}
	response417(w, r, result)
}

func response417NoCP(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_NO_CP, Message: "couple_need"}
	response417(w, r, result)
}

func response417(w http.ResponseWriter, r *http.Request, result *HttpResult) {
	responseJson(w, r, http.StatusExpectationFailed, result)
}

// response503 服务器维护
func response503(w http.ResponseWriter, r *http.Request) {
	result := &HttpResult{Code: RESULT_CODE_OK, Message: "http_status_503"}
	responseJson(w, r, http.StatusServiceUnavailable, result)
}

// responseJson
func responseJson(w http.ResponseWriter, r *http.Request, status int, body *HttpResult) {
	if body == nil {
		body = &HttpResult{Message: "", Data: struct {
		}{}}
	}
	if body.Data == nil {
		body.Data = struct {
		}{}
	}
	language := strings.ToLower(strings.TrimSpace(r.Header.Get(HEAD_PARAMS_LANGUAGE1)))
	if len(language) <= 0 {
		language = strings.ToLower(strings.TrimSpace(r.Header.Get(HEAD_PARAMS_LANGUAGE2)))
	}
	body.Message = utils.GetLanguage(language, body.Message)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	//err := json.NewEncoder(w).Encode(body) // 先把信息返回
	//if err != nil {
	//	utils.LogErr("response", err)
	//	utils.LogWarn("response-body", fmt.Sprintf("%+v", body))
	//}
	//_, err := w.Write(response)
	//panic(err) // 再终止协程，正常情况下这里返回的是nil，系统检测到的panic不走这里
	// ---------------新方法-----------------
	// 注意，这里的异常一般都是client断开的，所以就不返回err了，直接给nil
	_ = json.NewEncoder(w).Encode(body)
	panic(nil) // 再终止协程，正常情况下这里返回的是nil，系统检测到的panic不走这里
}
