package controllers

import (
	"encoding/json"
	"net/http"

	"libs/aliyun"
	"libs/utils"
	"models/entity"
	"services"
	"strconv"
)

func HandlerSms(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "sms")
	if r.Method == http.MethodPost {
		PostSms(w, r)
	} else if r.Method == http.MethodGet {
		GetSms(w, r)
	} else {
		response405(w, r)
	}
}

// PostSms
func PostSms(w http.ResponseWriter, r *http.Request) {
	sms := &entity.Sms{}
	checkRequestBody(w, r, sms)
	// 检验是否可以发送短信
	err := services.SmsEnableSend(sms.Phone, sms.SendType)
	response417ErrDialog(w, r, err)
	// 发送短信参数
	templateCode := ""
	templateParam := ""
	if sms.SendType == entity.SMS_TYPE_REGISTER || sms.SendType == entity.SMS_TYPE_LOGIN ||
		sms.SendType == entity.SMS_TYPE_FORGET || sms.SendType == entity.SMS_TYPE_PHONE ||
		sms.SendType == entity.SMS_TYPE_LOCK {
		// 验证码生成
		code := services.CreateSmsValidateCode()
		sms.Content = code
		param := &struct {
			Code string `json:"code"`
		}{code}
		bytes, _ := json.Marshal(param)
		// 短信模板生成
		templateParam = string(bytes)
		templateCode = utils.GetConfigStr("conf", "third.conf", "sms", "template_code_validate")
	} else {
		response405(w, r)
	}
	// 先在数据库里增加发送短信记录
	sms, err = services.AddSms(sms)
	response417ErrDialog(w, r, err)
	// 然后发送短信
	err = aliyun.SendSms(sms.Phone, templateCode, templateParam)
	response417ErrDialog(w, r, err)
	// 返回
	response200Toast(w, r, "sms_send_out")
}

// GetSms
func GetSms(w http.ResponseWriter, r *http.Request) {
	user := checkTokenUser(w, r)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	total, _ := strconv.ParseBool(values.Get("total"))
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	if list {
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		phone := values.Get("phone")
		sendType, _ := strconv.Atoi(values.Get("type"))
		page, _ := strconv.Atoi(values.Get("page"))
		smsList, err := services.GetSmsListByCreate(start, end, phone, sendType, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SmsList []*entity.Sms `json:"smsList"`
		}{smsList})
	} else if total {
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		phone := values.Get("phone")
		sendType, _ := strconv.Atoi(values.Get("type"))
		total := services.GetSmsTotalByCreateWithDel(start, end, phone, sendType)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
