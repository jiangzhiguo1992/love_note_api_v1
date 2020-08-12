package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerMatchWork(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_match_work")
	if r.Method == http.MethodPost {
		PostMatchWork(w, r)
	} else if r.Method == http.MethodDelete {
		DelMatchWork(w, r)
	} else if r.Method == http.MethodGet {
		GetMatchWork(w, r)
	} else {
		response405(w, r)
	}
}

// PostMatchWork
func PostMatchWork(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	couple := user.Couple
	// 接受参数
	work := &entity.MatchWork{}
	checkRequestBody(w, r, work)
	// 开始插入
	work, err := services.AddMatchWork(user.Id, couple.Id, work)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		MatchWork *entity.MatchWork `json:"matchWork"`
	}{work})
}

// DelMatchWork
func DelMatchWork(w http.ResponseWriter, r *http.Request) {
	user := checkTokenCouple(w, r)
	// 接受参数
	values := r.URL.Query()
	mwid, _ := strconv.ParseInt(values.Get("mwid"), 10, 64)
	// 开始删除
	err := services.DelMatchWork(user.Id, mwid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetMatchWork
func GetMatchWork(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	cid := services.GetCoupleIdByUser(user)
	// 接受参数
	values := r.URL.Query()
	mpid, _ := strconv.ParseInt(values.Get("mpid"), 10, 64)
	admin, _ := strconv.ParseBool(values.Get("admin"))
	report, _ := strconv.ParseBool(values.Get("report"))
	our, _ := strconv.ParseBool(values.Get("our"))
	total, _ := strconv.ParseBool(values.Get("total"))
	if mpid > 0 {
		order, _ := strconv.Atoi(values.Get("order"))
		page, _ := strconv.Atoi(values.Get("page"))
		list, err := services.GetMatchWorkListByPeriodOrder(user.Id, cid, mpid, order, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			MatchWorkList []*entity.MatchWork `json:"matchWorkList"`
		}{list})
	} else if admin {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
		mpid, _ := strconv.ParseInt(values.Get("mpid"), 10, 64)
		page, _ := strconv.Atoi(values.Get("page"))
		list, err := services.GetMatchWorkList(uid, mpid, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			MatchWorkList []*entity.MatchWork `json:"matchWorkList"`
		}{list})
	} else if report {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		list, err := services.GetMatchWorkReportList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			MatchWorkList []*entity.MatchWork `json:"matchWorkList"`
		}{list})
	} else if our {
		if cid <= 0 {
			response417NoCP(w, r)
		}
		kind, _ := strconv.Atoi(values.Get("kind"))
		page, _ := strconv.Atoi(values.Get("page"))
		list, err := services.GetMatchWorkListByCoupleKind(user.Id, cid, kind, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			MatchWorkList []*entity.MatchWork `json:"matchWorkList"`
		}{list})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
		kind, _ := strconv.Atoi(values.Get("kind"))
		total := services.GetMatchWorkTotalByKindWithDel(create, kind)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
