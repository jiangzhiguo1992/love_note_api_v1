package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerAwardRule(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_award_rule")
	if r.Method == http.MethodPost {
		PostAwardRule(w, r)
	} else if r.Method == http.MethodDelete {
		DelAwardRule(w, r)
	} else if r.Method == http.MethodGet {
		GetAwardRule(w, r)
	} else {
		response405(w, r)
	}
}

// PostAwardRule
func PostAwardRule(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	awardRule := &entity.AwardRule{}
	checkRequestBody(w, r, awardRule)
	// 开始插入
	awardRule, err := services.AddAwardRule(user.Id, couple.Id, awardRule)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		AwardRule *entity.AwardRule `json:"awardRule"`
	}{awardRule})
}

// DelAwardRule
func DelAwardRule(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	arid, _ := strconv.ParseInt(values.Get("arid"), 10, 64)
	// 开始删除
	err := services.DelAwardRule(user.Id, couple.Id, arid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetAwardRule
func GetAwardRule(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		awardRuleList, err := services.GetAwardRuleListByCouple(user.Id, couple.Id, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			AwardRuleList []*entity.AwardRule `json:"awardRuleList"`
		}{awardRuleList})
	} else {
		response405(w, r)
	}
}
