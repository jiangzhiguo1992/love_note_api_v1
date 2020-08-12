package controllers

import (
	"libs/utils"
	"models/entity"
	"net/http"
	"services"
	"strconv"
	"time"
)

func HandlerSouvenir(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "note_souvenir")
	if r.Method == http.MethodPost {
		PostSouvenir(w, r)
	} else if r.Method == http.MethodDelete {
		DelSouvenir(w, r)
	} else if r.Method == http.MethodPut {
		PutSouvenir(w, r)
	} else if r.Method == http.MethodGet {
		GetSouvenir(w, r)
	} else {
		response405(w, r)
	}
}

// PostSouvenir
func PostSouvenir(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	souvenir := &entity.Souvenir{}
	checkRequestBody(w, r, souvenir)
	// 防止通过年限来增加关联数据 1000秒的误差
	//if souvenir.HappenAt < (user.Birthday - 1000) {
	//	response417Toast(w, r, "limit_happen_err")
	//}
	// 开始插入
	souvenir, err := services.AddSouvenir(user.Id, couple.Id, souvenir)
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "db_add_success", struct {
		Souvenir *entity.Souvenir `json:"souvenir"`
	}{souvenir})
}

// DelSouvenir
func DelSouvenir(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	// 开始删除
	err := services.DelSouvenir(user.Id, couple.Id, sid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// PutSouvenir
func PutSouvenir(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	year, _ := strconv.Atoi(values.Get("year"))
	souvenir := &entity.Souvenir{}
	checkRequestBody(w, r, souvenir)
	var err error
	if year <= 0 {
		// 修改实体
		//if souvenir.HappenAt < (user.Birthday - 1000) {
		//	// 防止通过年限来增加关联数据 1000秒的误差
		//	response417Toast(w, r, "limit_happen_err")
		//}
		souvenir, err = services.UpdateSouvenir(user.Id, couple.Id, souvenir)
		response417ErrDialog(w, r, err)
	} else {
		// 修改关联
		if year < utils.GetCSTDateByUnix(souvenir.HappenAt).Year() || year > utils.GetCSTDateByUnix(time.Now().Unix()).Year() {
			// 防止添加不能显示的年限关联
			response417Toast(w, r, "limit_happen_err")
		}
		souvenir = services.UpdateSouvenirForeign(user.Id, couple.Id, year, souvenir)
	}
	// 返回
	response200DataToast(w, r, "db_update_success", struct {
		Souvenir *entity.Souvenir `json:"souvenir"`
	}{souvenir})
}

// GetSouvenir
func GetSouvenir(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	sid, _ := strconv.ParseInt(values.Get("sid"), 10, 64)
	if list {
		page, _ := strconv.Atoi(values.Get("page"))
		done, _ := strconv.ParseBool(values.Get("done"))
		souvenirList, err := services.GetSouvenirListByCouple(user.Id, couple.Id, done, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			SouvenirList []*entity.Souvenir `json:"souvenirList"`
		}{souvenirList})
	} else if sid > 0 {
		souvenir, err := services.GetSouvenirByIdWithForeign(user.Id, couple.Id, sid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Souvenir *entity.Souvenir `json:"souvenir"`
		}{souvenir})
	} else {
		response405(w, r)
	}
}
