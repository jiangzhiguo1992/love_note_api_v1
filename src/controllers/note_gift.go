package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerGift(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_gift")
	if r.Method == http.MethodPost {
		PostGift(w, r)
	} else if r.Method == http.MethodDelete {
		DelGift(w, r)
	} else if r.Method == http.MethodPut {
		PutGift(w, r)
	} else if r.Method == http.MethodGet {
		GetGift(w, r)
	} else {
		response405(w, r)
	}
}

// PostGift
func PostGift(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	gift := &entity.Gift{}
	checkRequestBody(w, r, gift)
	// 数据检查
	if gift.ReceiveId != user.Id && gift.ReceiveId != services.GetTaId(user) {
		response417Toast(w, r, "gift_receive_nil")
	}
	// 开始插入
	gift, err := services.AddGift(user.Id, couple.Id, gift)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Gift *entity.Gift `json:"gift"`
	}{gift})
}

// DelGift
func DelGift(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	gid, _ := strconv.ParseInt(values.Get("gid"), 10, 64)
	// 开始删除
	err := services.DelGift(user.Id, couple.Id, gid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutGift
func PutGift(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	gift := &entity.Gift{}
	checkRequestBody(w, r, gift)
	// 数据检查
	if gift.ReceiveId != user.Id && gift.ReceiveId != services.GetTaId(user) {
		response417Toast(w, r, "gift_receive_nil")
	}
	// 开始插入
	gift, err := services.UpdateGift(user.Id, couple.Id, gift)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Gift *entity.Gift `json:"gift"`
	}{gift})
}

// GetGift
func GetGift(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		who, _ := strconv.Atoi(values.Get("who"))
		page, _ := strconv.Atoi(values.Get("page"))
		var suid int64
		if who == services.LIST_WHO_BY_ME {
			suid = user.Id
		} else if who == services.LIST_WHO_BY_TA {
			suid = services.GetTaId(user)
		} else {
			suid = 0
		}
		giftList, err := services.GetGiftListByUserCouple(user.Id, suid, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			GiftList []*entity.Gift `json:"giftList"`
		}{giftList})
	} else {
		response405(w, r)
	}
}
