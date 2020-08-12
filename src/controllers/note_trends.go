package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerTrends(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_trends")
	if r.Method == http.MethodGet {
		GetTrends(w, r)
	} else {
		response405(w, r)
	}
}

// GetTrends
func GetTrends(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 获取参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	group, _ := strconv.ParseBool(values.Get("group"))
	total, _ := strconv.ParseBool(values.Get("total"))
	if list {
		admin, _ := strconv.ParseBool(values.Get("admin"))
		page, _ := strconv.Atoi(values.Get("page"))
		var trendsList []*entity.Trends
		var err error
		if admin && services.IsAdminister(user) {
			uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
			cid, _ := strconv.ParseInt(values.Get("cid"), 10, 64)
			actType, _ := strconv.Atoi(values.Get("act_type"))
			conType, _ := strconv.Atoi(values.Get("con_type"))
			trendsList, err = services.GetTrendsList(uid, cid, actType, conType, page)
		} else {
			create, _ := strconv.ParseInt(values.Get("create"), 10, 64)
			trendsList, err = services.GetTrendsListByCreateUserCouple(create, user.Id, couple.Id, page)
		}
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			TrendsList []*entity.Trends `json:"trendsList"`
		}{trendsList})
	} else if group {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		infoList, err := services.GetTrendsContentTypeList(start, end)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			InfoList []*entity.FiledInfo `json:"infoList"`
		}{infoList})
	} else if total {
		if !services.GetVipLimitByCouple(couple.Id).NoteTotalEnable {
			response417Dialog(w, r, "vip_need")
		}
		noteTotal := services.GetTrendsNoteTotal(couple.Id)
		// 返回
		response200Data(w, r, struct {
			NoteTotal *services.NoteTotal `json:"noteTotal"`
		}{noteTotal})
	} else {
		response405(w, r)
	}
}
