package controllers

import (
	"net/http"
	"strconv"

	"models/entity"
	"services"
)

func HandlerTopicMessage(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_message")
	if r.Method == http.MethodGet {
		GetTopicMessage(w, r)
	} else {
		response405(w, r)
	}
}

// GetTopicMessage
func GetTopicMessage(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	mine, _ := strconv.ParseBool(values.Get("mine"))
	if mine {
		kind, _ := strconv.Atoi(values.Get("kind"))
		page, _ := strconv.Atoi(values.Get("page"))
		messageList, err := services.GetTopicMessageListByToUserCoupleKind(user.Id, couple.Id, kind, page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			TopicMessageList []*entity.TopicMessage `json:"topicMessageList"`
		}{messageList})
	} else {
		response405(w, r)
	}
}
