package controllers

import (
	"net/http"

	"services"
	"strconv"
)

func HandlerPostRead(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_read")
	if r.Method == http.MethodPost {
		PostPostRead(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostRead
func PostPostRead(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	pid, _ := strconv.ParseInt(values.Get("pid"), 10, 64)
	// 开始插入
	_, err := services.AddPostRead(user.Id, pid)
	response417ErrToast(w, r, err)
	// 返回
	response200Toast(w, r, "")
}
