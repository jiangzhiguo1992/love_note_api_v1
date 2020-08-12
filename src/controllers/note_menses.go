package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerMenses(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_menses")
	if r.Method == http.MethodPost {
		PostMenses(w, r)
	} else if r.Method == http.MethodDelete {
		DelMenses(w, r)
	} else if r.Method == http.MethodGet {
		GetMenses(w, r)
	} else {
		response405(w, r)
	}
}

// PostMenses 发表
func PostMenses(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	menses := &entity.Menses{}
	checkRequestBody(w, r, menses)
	// 数据检查
	if user.Sex != entity.USER_SEX_GIRL {
		response417Toast(w, r, "menses_sex_limit")
	}
	// 开始插入
	menses, err := services.AddMenses(user.Id, couple.Id, menses)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Menses *entity.Menses `json:"menses"`
	}{menses})
}

// DelMenses
func DelMenses(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	year, _ := strconv.Atoi(values.Get("year"))
	month, _ := strconv.Atoi(values.Get("month"))
	day, _ := strconv.Atoi(values.Get("day"))
	mid, _ := strconv.ParseInt(values.Get("mid"), 10, 64)
	// 开始删除
	if year > 0 && month > 0 && day > 0 {
		err := services.DelMensesByDate(user.Id, couple.Id, year, month, day)
		response417ErrToast(w, r, err)
	} else if mid > 0 {
		err := services.DelMenses(user.Id, couple.Id, mid)
		response417ErrToast(w, r, err)
	}
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetMenses 获取
func GetMenses(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	latest, _ := strconv.ParseBool(values.Get("latest"))
	date, _ := strconv.ParseBool(values.Get("date"))
	if latest {
		mensesInfo := &services.MensesInfo{CanMe: false, CanTa: false}
		// 我的
		var mensesMe *entity.Menses
		if user.Sex == entity.USER_SEX_GIRL {
			mensesInfo.CanMe = true
			mensesMe, _ = services.GetMensesLatestByUserCouple(user.Id, couple.Id)
		}
		// ta的
		taId := services.GetTaId(user)
		ta, _ := services.GetUserById(taId)
		var mensesTa *entity.Menses
		if ta == nil {
			response200Show(w, r, "nil_user")
		} else if ta.Sex == entity.USER_SEX_GIRL {
			mensesInfo.CanTa = true
			mensesTa, _ = services.GetMensesLatestByUserCouple(taId, couple.Id)
		}
		// 返回
		response200Data(w, r, struct {
			MensesInfo *services.MensesInfo `json:"mensesInfo"`
			MensesMe   *entity.Menses       `json:"mensesMe"`
			MensesTa   *entity.Menses       `json:"mensesTa"`
		}{mensesInfo, mensesMe, mensesTa})
	} else if date {
		mine, _ := strconv.ParseBool(values.Get("mine"))
		var searchUid int64
		// 性别检查
		if mine {
			if user.Sex != entity.USER_SEX_GIRL {
				response200Show(w, r, "menses_sex_limit")
			} else {
				searchUid = user.Id
			}
		} else {
			taId := services.GetTaId(user)
			ta, _ := services.GetUserById(taId)
			if ta == nil {
				response200Show(w, r, "nil_user")
			} else if ta.Sex != entity.USER_SEX_GIRL {
				response200Show(w, r, "menses_sex_limit")
			} else {
				searchUid = ta.Id
			}
		}
		year, _ := strconv.Atoi(values.Get("year"))
		month, _ := strconv.Atoi(values.Get("month"))
		// 数据
		mensesList, err := services.GetMensesListByUserCoupleYearMonth(searchUid, couple.Id, year, month)
		response417ErrDialog(w, r, err) // dialog
		// 返回
		response200Data(w, r, struct {
			MensesList []*entity.Menses `json:"mensesList"`
		}{mensesList})
	} else {
		response405(w, r)
	}
}
