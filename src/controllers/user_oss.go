package controllers

import (
	"net/http"
	"services"
	"strconv"
)

func HandlerOss(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "oss")
	if r.Method == http.MethodGet {
		GetOss(w, r)
	} else {
		response405(w, r)
	}
}

// GetOss
func GetOss(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	cid := services.GetCoupleIdByUser(user)
	// 接收数据
	values := r.URL.Query()
	admin, _ := strconv.ParseBool(values.Get("admin"))
	// ossInfo
	var info *services.OssInfo
	var err error
	if admin && services.IsAdminister(user) {
		info, err = services.GetOssInfoByAdmin(user.Id, cid, user.UserToken)
	} else {
		info, err = services.GetOssInfoByUserCouple(user.Id, cid, user.UserToken)
	}
	response417ErrDialog(w, r, err)
	// 返回
	response200Data(w, r, struct {
		OssInfo *services.OssInfo `json:"ossInfo"`
	}{info})
}
