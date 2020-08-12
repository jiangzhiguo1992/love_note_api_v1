package controllers

import (
	"net/http"
	"services"
)

func HandlerTopicHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetTopicHome(w, r)
	} else {
		response405(w, r)
	}
}

// GetTopicHome
func GetTopicHome(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// postKindInfo
	postKindInfoList := make([]*services.PostKindInfo, 0)
	if getModelStatus("topic_post") {
		postKindInfoList = services.GetPostKindInfoListWithTopicInfo(user)
	}
	// count
	commonCount := &services.CommonCount{}
	if user != nil && couple != nil {
		commonCount.TopicMsgNewCount = services.GetTopicMessageCountByUserCouple(user.Id, couple.Id)
	}
	// 返回
	response200Data(w, r, struct {
		PostKindInfoList []*services.PostKindInfo `json:"postKindInfoList"`
		CommonCount      *services.CommonCount    `json:"commonCount"`
	}{
		postKindInfoList,
		commonCount,
	})
}
