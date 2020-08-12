package services

import (
	"encoding/json"
	"errors"
	"libs/aliyun"
	"libs/utils"
	"models/mysql"
	"strings"
)

// GetPushInfo
func GetPushInfo(uid int64) *PushInfo {
	// platform
	entry, err := mysql.GetEntryLatestByUser(uid)
	if err != nil {
		return nil
	} else if entry == nil {
		return nil
	} else if len(entry.Platform) <= 0 {
		return nil
	}
	platform := strings.ToLower(strings.TrimSpace(entry.Platform))
	platformAndroid := utils.GetConfigStr("conf", "app.conf", "phone", "platform_android")
	platformIos := utils.GetConfigStr("conf", "app.conf", "phone", "platform_ios")
	info := &PushInfo{}
	if platform == platformAndroid {
		// android
		info.AliAppKey = utils.GetConfigStr("conf", "third.conf", "push", "ali_android_app_key")
		info.AliAppSecret = utils.GetConfigStr("conf", "third.conf", "push", "ali_android_app_secret")
		info.MiAppId = utils.GetConfigStr("conf", "third.conf", "push", "mi_app_id")
		info.MiAppKey = utils.GetConfigStr("conf", "third.conf", "push", "mi_app_key")
		info.OppoAppKey = utils.GetConfigStr("conf", "third.conf", "push", "oppo_app_key")
		info.OppoAppSecret = utils.GetConfigStr("conf", "third.conf", "push", "oppo_app_secret")
		info.ChannelId = utils.GetConfigStr("conf", "third.conf", "push", "channel_id")
		info.NoticeLight = utils.GetConfigBool("conf", "third.conf", "push", "notice_light")
		info.NoticeSound = utils.GetConfigBool("conf", "third.conf", "push", "notice_sound")
		info.NoticeVibrate = utils.GetConfigBool("conf", "third.conf", "push", "notice_vibrate")
		info.NoStartHour = utils.GetConfigInt("conf", "third.conf", "push", "no_start_hour")
		info.NoEndHour = utils.GetConfigInt("conf", "third.conf", "push", "no_end_hour")
	} else if platform == platformIos {
		// ios
		info.AliAppKey = utils.GetConfigStr("conf", "third.conf", "push", "ali_ios_app_key")
		info.AliAppSecret = utils.GetConfigStr("conf", "third.conf", "push", "ali_ios_app_secret")
		info.ChannelId = utils.GetConfigStr("conf", "third.conf", "push", "channel_id")
		info.NoticeLight = utils.GetConfigBool("conf", "third.conf", "push", "notice_light")
		info.NoticeSound = utils.GetConfigBool("conf", "third.conf", "push", "notice_sound")
		info.NoticeVibrate = utils.GetConfigBool("conf", "third.conf", "push", "notice_vibrate")
		info.NoStartHour = utils.GetConfigInt("conf", "third.conf", "push", "no_start_hour")
		info.NoEndHour = utils.GetConfigInt("conf", "third.conf", "push", "no_end_hour")
	}
	return info
}

func CreatePush(uid, toUid, conId int64, title, body string, tp int) *Push {
	push := &Push{
		UserId:      uid,
		ToUserId:    toUid,
		Title:       title,
		ContentText: body,
		ContentType: tp,
		ContentId:   conId,
	}
	return push
}

func AddPushInCouple(uid, conId int64, title, content string, tp int) {
	// push
	couple, err := GetCoupleVisibleByUser(uid)
	if err != nil {
		return
	} else if couple == nil || couple.Id <= 0 {
		return
	}
	toUid := couple.InviteeId
	if toUid == uid {
		toUid = couple.CreatorId
	}
	entry, err := mysql.GetEntryLatestByUser(toUid)
	if err != nil {
		return
	} else if entry == nil {
		return
	}
	title = utils.GetLanguage(entry.Language, title)
	content = utils.GetLanguage(entry.Language, content)
	push := CreatePush(uid, toUid, conId, title, content, tp)
	_, _ = AddPush(push)
}

func AddPush(push *Push) (*Push, error) {
	if push == nil {
		return nil, nil
	} else if push.ToUserId <= 0 {
		return nil, errors.New("nil_user")
	} else if !utils.GetConfigBool("conf", "model.conf", "model", "push") {
		return nil, errors.New("request_model_close")
	}
	//else if len(push.ContentText) > 20 {
	//	push.ContentText = push.ContentText[1 : 10-1]
	//}
	// platform
	entry, err := mysql.GetEntryLatestByUser(push.ToUserId)
	if err != nil {
		return push, err
	} else if entry == nil {
		return push, errors.New("nil_entry")
	} else if len(strings.TrimSpace(entry.Platform)) <= 0 {
		return push, errors.New("nil_platform")
	}
	push.Platform = entry.Platform
	// length
	if len(strings.TrimSpace(push.Title)) <= 0 {
		push.Title = "YuSheng"
	} else if len(push.Title) > 30 {
		push.Title = push.Title[0:29]
	}
	if len(strings.TrimSpace(push.ContentText)) <= 0 {
		push.ContentText = "~"
	} else if len(push.ContentText) > 50 {
		push.ContentText = push.ContentText[0:49]
	}
	// extraJson
	bytes, err := json.Marshal(push)
	if err != nil {
		return push, err
	}
	// 开始推送
	err = aliyun.SendPushNotice2User(push.ToUserId, push.Platform, push.Title, push.ContentText, string(bytes))
	return push, err
}
