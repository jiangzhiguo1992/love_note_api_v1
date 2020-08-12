package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerSign(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_sign")
	if r.Method == http.MethodPost {
		PostSign(w, r)
	} else if r.Method == http.MethodGet {
		GetSign(w, r)
	} else {
		response405(w, r)
	}
}

// PostSign
func PostSign(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 开始插入
	sign, err := services.AddSign(user.Id, couple.Id)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "sign_success", struct {
		Sign *entity.Sign `json:"sign"`
	}{sign})
}

// GetSign
func GetSign(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	date, _ := strconv.ParseBool(values.Get("date"))
	list, _ := strconv.ParseBool(values.Get("list"))
	total, _ := strconv.ParseBool(values.Get("total"))
	if date {
		year, _ := strconv.Atoi(values.Get("year"))
		month, _ := strconv.Atoi(values.Get("month"))
		signList, err := services.GetSignListByCoupleYearMonth(couple.Id, year, month)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SignList []*entity.Sign `json:"signList"`
		}{signList})
	} else if list {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		cid, _ := strconv.ParseInt(values.Get("cid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		signList, err := services.GetSignList(uid, cid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SignList []*entity.Sign `json:"signList"`
		}{signList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		year, _ := strconv.Atoi(values.Get("year"))
		month, _ := strconv.Atoi(values.Get("month"))
		day, _ := strconv.Atoi(values.Get("day"))
		total := services.GetSignTotalWithDel(year, month, day)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
