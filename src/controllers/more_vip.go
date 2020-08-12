package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerVip(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_vip")
	if r.Method == http.MethodPost {
		PostVip(w, r)
	} else if r.Method == http.MethodGet {
		GetVip(w, r)
	} else {
		response405(w, r)
	}
}

// PostVip
func PostVip(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接受参数
	vip := &entity.Vip{}
	checkRequestBody(w, r, vip)
	// 开始插入
	vip, err := services.AddVipByAdmin(vip)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Vip *entity.Vip `json:"vip"`
	}{vip})
}

// GetVip
func GetVip(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	home, _ := strconv.ParseBool(values.Get("home"))
	list, _ := strconv.ParseBool(values.Get("list"))
	expireDays, _ := strconv.ParseBool(values.Get("expire_days"))
	total, _ := strconv.ParseBool(values.Get("total"))
	if home {
		vip, err := services.GetVipLatest(couple.Id)
		vipLimit := services.GetVipLimitByCouple(couple.Id)
		vipYesLimit := services.GetVipLimit(true)
		vipNoLimit := services.GetVipLimit(false)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Vip         *entity.Vip        `json:"vip"`
			VipLimit    *services.VipLimit `json:"vipLimit"`
			VipYesLimit *services.VipLimit `json:"vipYesLimit"`
			VipNoLimit  *services.VipLimit `json:"vipNoLimit"`
		}{vip, vipLimit, vipYesLimit, vipNoLimit})
	} else if list {
		page, _ := strconv.Atoi(values.Get("page"))
		admin, _ := strconv.ParseBool(values.Get("admin"))
		var vipList []*entity.Vip
		var err error
		if admin && services.IsAdminister(user) {
			uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
			cid, _ := strconv.ParseInt(values.Get("cid"), 10, 64)
			bid, _ := strconv.ParseInt(values.Get("bid"), 10, 64)
			fromType, _ := strconv.Atoi(values.Get("from_type"))
			vipList, err = services.GetVipList(uid, cid, bid, fromType, page)
		} else {
			vipList, err = services.GetVipListByCouple(couple.Id, page)
		}
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			VipList []*entity.Vip `json:"vipList"`
		}{vipList})
	} else if expireDays {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		infoList, err := services.GetVipExpireDaysListByCreate(start, end)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			InfoList []*entity.FiledInfo `json:"infoList"`
		}{infoList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		total := services.GetVipTotalByCreateWithDel(start, end)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
