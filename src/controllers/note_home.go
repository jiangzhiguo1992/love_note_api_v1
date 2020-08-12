package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerNoteHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetNoteHome(w, r)
	} else {
		response405(w, r)
	}
}

// GetNoteHome 获取note首页
func GetNoteHome(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 获取参数
	values := r.URL.Query()
	near, _ := strconv.ParseInt(values.Get("near"), 10, 64)
	// lock
	var lock *entity.Lock
	if getModelStatus("note_lock") {
		lock, _ = services.GetLockByUserCouple(user.Id, couple.Id)
	}
	// 返回最近的纪念日
	var latest *entity.Souvenir
	if getModelStatus("note_souvenir") {
		souvenirList, _ := services.GetSouvenirListByCouple(user.Id, couple.Id, true, -1)
		latest = services.GetSouvenirLatestByList(souvenirList, near)
	}
	// count
	taId := services.GetTaId(user)
	commonCount := &services.CommonCount{
		NoteTrendsNewCount: services.GetTrendsCountByUserCouple(user.Id, taId, couple.Id),
	}
	// 返回
	response200Data(w, r, struct {
		Lock           *entity.Lock          `json:"lock"`
		SouvenirLatest *entity.Souvenir      `json:"souvenirLatest"`
		CommonCount    *services.CommonCount `json:"commonCount"`
	}{lock, latest, commonCount})
}
