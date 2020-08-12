package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerVersion(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "version")
	if r.Method == http.MethodPost {
		PostVersion(w, r)
	} else if r.Method == http.MethodDelete {
		DelVersion(w, r)
	} else if r.Method == http.MethodGet {
		GetVersion(w, r)
	} else {
		response405(w, r)
	}
}

// PostVersion
func PostVersion(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接受参数
	version := &entity.Version{}
	checkRequestBody(w, r, version)
	// 开始插入
	version, err := services.AddVersion(version)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_add_success")
}

// DelVersion
func DelVersion(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// admin检查
	if !services.IsAdminister(user) {
		response200Toast(w, r, "")
	}
	// 接收数据
	values := r.URL.Query()
	vid, _ := strconv.ParseInt(values.Get("vid"), 10, 64)
	// 开始删除
	err := services.DelVersion(vid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "db_delete_success")
}

// GetVersion 获取最新的version
func GetVersion(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	list, _ := strconv.ParseBool(values.Get("list"))
	code, _ := strconv.Atoi(values.Get("code"))
	if list {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		versionList, err := services.GetVersionList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			VersionList []*entity.Version `json:"versionList"`
		}{versionList})
	} else if code > 0 {
		// 查询最新版本
		versionList, err := services.GetVersionListByCode(user.Id, code)
		response417ErrDialog(w, r, err)
		if versionList == nil || len(versionList) <= 0 {
			response200Toast(w, r, "version_is_latest")
		}
		// 返回
		response200Data(w, r, struct {
			VersionList []*entity.Version `json:"versionList"`
		}{versionList})
	} else {
		response405(w, r)
	}
}
