package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerMenses2(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_menses")
	if r.Method == http.MethodPost {
		PostMenses2(w, r)
	} else if r.Method == http.MethodGet {
		GetMenses2(w, r)
	} else {
		response405(w, r)
	}
}

// PostMenses2 发表
func PostMenses2(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	menses := &entity.Menses{}
	checkRequestBody(w, r, menses)
	var uid int64
	// 数据检查
	if menses.IsMe {
		// 我的
		if user.Sex != entity.USER_SEX_GIRL {
			response417Toast(w, r, "menses_sex_limit")
		}
		uid = user.Id
	} else {
		// ta的
		taId := services.GetTaId(user)
		ta, _ := services.GetUserById(taId)
		if ta == nil {
			response200Show(w, r, "nil_user")
		} else if ta.Sex != entity.USER_SEX_GIRL {
			response417Toast(w, r, "menses_sex_limit")
		}
		uid = ta.Id
	}
	// 开始插入
	menses2, err := services.AddMenses2(uid, couple.Id, menses)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		Menses2 *entity.Menses2 `json:"menses2"`
	}{menses2})
}

// GetMenses2 获取
func GetMenses2(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	date, _ := strconv.ParseBool(values.Get("date"))
	if date {
		mine, _ := strconv.ParseBool(values.Get("mine"))
		year, _ := strconv.Atoi(values.Get("year"))
		month, _ := strconv.Atoi(values.Get("month"))
		var searchUid int64
		// 性别检查
		if mine {
			if user.Sex != entity.USER_SEX_GIRL {
				response200Show(w, r, "menses_sex_limit") // show，App要展示
			} else {
				searchUid = user.Id
			}
		} else {
			taId := services.GetTaId(user)
			ta, _ := services.GetUserById(taId)
			if ta == nil {
				response200Show(w, r, "nil_user")
			} else if ta.Sex != entity.USER_SEX_GIRL {
				response200Show(w, r, "menses_sex_limit") // show，App要展示
			} else {
				searchUid = ta.Id
			}
		}
		// 数据
		menses2List, err := services.GetMenses2ListByUserCoupleYearMonth(searchUid, couple.Id, year, month)
		response417ErrDialog(w, r, err) // dialog
		// 返回
		response200Data(w, r, struct {
			Menses2List []*entity.Menses2 `json:"menses2List"`
		}{menses2List})
	} else {
		response405(w, r)
	}
}
