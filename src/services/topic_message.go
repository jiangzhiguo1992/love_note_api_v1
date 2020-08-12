package services

import (
	"errors"
	"fmt"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// CreateTopicMessage
func CreateTopicMessage(uid, cid, tuid, tcid int64, kind int, conText string, conId int64) *entity.TopicMessage {
	message := &entity.TopicMessage{}
	message.UserId = uid
	message.CoupleId = cid
	message.ToUserId = tuid
	message.ToCoupleId = tcid
	message.Kind = kind
	message.ContentText = conText
	message.ContentId = conId
	return message
}

// AddTopicMessage
func AddTopicMessage(language string, tm *entity.TopicMessage) (*entity.TopicMessage, error) {
	if tm == nil {
		utils.LogErr("AddTopicMessage", "缺失 TopicMessage")
		return nil, errors.New("nil_message")
	}
	if tm.Kind == entity.TOPIC_MESSAGE_KIND_POST_BE_REPORT ||
		tm.Kind == entity.TOPIC_MESSAGE_KIND_COMMENT_BE_REPORT ||
		tm.Kind == entity.TOPIC_MESSAGE_KIND_POST_BE_POINT ||
		tm.Kind == entity.TOPIC_MESSAGE_KIND_COMMENT_BE_POINT ||
		tm.Kind == entity.TOPIC_MESSAGE_KIND_POST_BE_COLLECT {
		// 举报，点赞，收藏 暂时不要
		return nil, nil
	}
	if tm.UserId <= 0 {
		utils.LogErr("AddTopicMessage", "缺失 UserId "+fmt.Sprintf("%+v", tm))
		return nil, errors.New("nil_user")
	} else if tm.CoupleId <= 0 {
		utils.LogErr("AddTopicMessage", "缺失 CoupleId "+fmt.Sprintf("%+v", tm))
		return nil, errors.New("nil_couple")
	} else if tm.ToUserId <= 0 {
		utils.LogErr("AddTopicMessage", "缺失 ToUserId "+fmt.Sprintf("%+v", tm))
		return nil, errors.New("nil_user")
	} else if tm.ToCoupleId <= 0 {
		utils.LogErr("AddTopicMessage", "缺失 ToCoupleId "+fmt.Sprintf("%+v", tm))
		return nil, errors.New("nil_couple")
	} else if tm.Kind <= entity.TOPIC_MESSAGE_KIND_ALL {
		utils.LogErr("AddTopicMessage", "缺失 Kind "+fmt.Sprintf("%+v", tm))
		return nil, errors.New("data_err")
	} else if tm.ContentId <= 0 {
		utils.LogErr("AddTopicMessage", "缺失 ContentId "+fmt.Sprintf("%+v", tm))
		return nil, errors.New("data_err")
	}
	// mysql
	tm, err := mysql.AddTopicMessage(tm)
	// 同步
	go func() {
		// push
		if tm.Kind == entity.TOPIC_MESSAGE_KIND_OFFICIAL_TEXT {
			title := utils.GetLanguage(language, "push_title_new_notice")
			push := CreatePush(tm.UserId, tm.ToUserId, tm.ContentId, title, tm.ContentText, entity.PUSH_TYPE_TOPIC_MESSAGE)
			AddPush(push)
		} else if tm.Kind == entity.TOPIC_MESSAGE_KIND_JAB_IN_POST {
			title := utils.GetLanguage(language, "push_title_topic_jab")
			push := CreatePush(tm.UserId, tm.ToUserId, tm.ContentId, title, tm.ContentText, entity.PUSH_TYPE_TOPIC_POST)
			AddPush(push)
		} else if tm.Kind == entity.TOPIC_MESSAGE_KIND_JAB_IN_COMMENT {
			title := utils.GetLanguage(language, "push_title_topic_jab")
			push := CreatePush(tm.UserId, tm.ToUserId, tm.ContentId, title, tm.ContentText, entity.PUSH_TYPE_TOPIC_COMMENT)
			AddPush(push)
		} else if tm.Kind == entity.TOPIC_MESSAGE_KIND_POST_BE_COMMENT {
			title := utils.GetLanguage(language, "push_title_new_comment")
			push := CreatePush(tm.UserId, tm.ToUserId, tm.ContentId, title, tm.ContentText, entity.PUSH_TYPE_TOPIC_POST)
			AddPush(push)
		} else if tm.Kind == entity.TOPIC_MESSAGE_KIND_COMMENT_BE_REPLY {
			title := utils.GetLanguage(language, "push_title_new_comment")
			push := CreatePush(tm.UserId, tm.ToUserId, tm.ContentId, title, tm.ContentText, entity.PUSH_TYPE_TOPIC_COMMENT)
			AddPush(push)
		}
	}()
	return tm, err
}

// GetTopicMessageListByToUserCoupleKind
func GetTopicMessageListByToUserCoupleKind(uid, cid int64, kind, page int) ([]*entity.TopicMessage, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	limit := GetPageSizeLimit().TopicMessage
	offset := page * limit
	list, err := mysql.GetTopicMessageListByToUserCoupleKind(uid, cid, kind, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_message")
		} else {
			return nil, nil
		}
	}
	// 额外数据，不能缓存用户数据
	for _, v := range list {
		v.Couple, _ = GetCoupleVisibleByUser(v.UserId)
	}
	// 同步
	go func() {
		AddTopicMessageBrowse(uid, cid)
	}()
	return list, err
}
