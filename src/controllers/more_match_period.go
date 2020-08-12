package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerMatchPeriod(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_match_period")
	if r.Method == http.MethodPost {
		PostMatchPeriod(w, r)
	} else if r.Method == http.MethodDelete {
		DelMatchPeriod(w, r)
	} else if r.Method == http.MethodGet {
		GetMatchPeriod(w, r)
	} else {
		response405(w, r)
	}
}

// PostMatchPeriod
func PostMatchPeriod(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接受参数
	period := &entity.MatchPeriod{}
	checkRequestBody(w, r, period)
	// 开始插入
	period, err := services.AddMatchPeriod(period)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		MatchPeriod *entity.MatchPeriod `json:"matchPeriod"`
	}{period})
}

// DelMatchPeriod
func DelMatchPeriod(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接受参数
	values := r.URL.Query()
	mpid, _ := strconv.ParseInt(values.Get("mpid"), 10, 64)
	// 开始删除
	err := services.DelMatchPeriod(mpid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetMatchPeriod
func GetMatchPeriod(w http.ResponseWriter, r *http.Request) {
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		kind, _ := strconv.Atoi(values.Get("kind"))
		page, _ := strconv.Atoi(values.Get("page"))
		matchPeriodList, err := services.GetMatchPeriodList(kind, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			MatchPeriodList []*entity.MatchPeriod `json:"matchPeriodList"`
		}{matchPeriodList})
	} else {
		response405(w, r)
	}
}
