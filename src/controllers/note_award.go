package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerAward(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_award")
	if r.Method == http.MethodPost {
		PostAward(w, r)
	} else if r.Method == http.MethodDelete {
		DelAward(w, r)
	} else if r.Method == http.MethodGet {
		GetAward(w, r)
	} else {
		response405(w, r)
	}
}

// PostAward
func PostAward(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	award := &entity.Award{}
	checkRequestBody(w, r, award)
	// 数据检查
	if award.HappenId != user.Id && award.HappenId != services.GetTaId(user) {
		response417Toast(w, r, "award_nil_happen_id")
	}
	// 开始插入
	award, err := services.AddAward(user.Id, couple.Id, award)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Award *entity.Award `json:"award"`
	}{award})
}

// DelAward
func DelAward(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	aid, _ := strconv.ParseInt(values.Get("aid"), 10, 64)
	// 开始删除
	err := services.DelAward(user.Id, couple.Id, aid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetAward
func GetAward(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	score, _ := strconv.ParseBool(values.Get("score"))
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
		awardList, err := services.GetAwardListByUserCouple(user.Id, suid, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			AwardList []*entity.Award `json:"awardList"`
		}{awardList})
	} else if score {
		taId := services.GetTaId(user)
		awardScoreMe, _ := services.GetAwardScoreByUserCouple(user.Id, couple.Id)
		awardScoreTa, _ := services.GetAwardScoreByUserCouple(taId, couple.Id)
		// 返回
		response200Data(w, r, struct {
			AwardScoreMe *entity.AwardScore `json:"awardScoreMe"`
			AwardScoreTa *entity.AwardScore `json:"awardScoreTa"`
		}{awardScoreMe, awardScoreTa})
	} else {
		response405(w, r)
	}
}
