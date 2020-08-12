package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerApi(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "api")
	if r.Method == http.MethodGet {
		GetApi(w, r)
	} else {
		response405(w, r)
	}
}

// GetApi
func GetApi(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	uri, _ := strconv.ParseBool(values.Get("uri"))
	total, _ := strconv.ParseBool(values.Get("total"))
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	if list {
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		apiList, err := services.GetApiList(start, end, uid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			ApiList []*entity.Api `json:"apiList"`
		}{apiList})
	} else if uri {
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		infoList, err := services.GetApiUriListByCreate(start, end)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			InfoList []*entity.FiledInfo `json:"infoList"`
		}{infoList})
	} else if total {
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		total := services.GetApiTotalByCreateWithDel(start, end)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
