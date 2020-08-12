package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// TogglePostPoint
func TogglePostPoint(uid, cid int64, pp *entity.PostPoint) (*entity.PostPoint, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pp == nil || pp.PostId <= 0 {
		return nil, errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pp.PostId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// mysql
	old, err := mysql.GetPostPointByUserCouple(uid, cid, p.Id)
	if err != nil {
		return old, err
	} else if old == nil || old.Id <= 0 {
		// 没点赞
		pp.UserId = uid
		pp.CoupleId = cid
		pp, err = mysql.AddPostPoint(pp)
	} else {
		// 已点赞
		if old.Status >= entity.STATUS_VISIBLE {
			old.Status = entity.STATUS_DELETE
		} else {
			old.Status = entity.STATUS_VISIBLE
		}
		pp, err = mysql.UpdatePostPoint(old)
	}
	if pp == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		// post
		if pp.Status >= entity.STATUS_VISIBLE {
			p.PointCount = p.PointCount + 1
		} else {
			p.PointCount = p.PointCount - 1
		}
		UpdatePostCount(p, false)
		// message
		if uid != p.UserId &&
			pp.Status >= entity.STATUS_VISIBLE &&
			p.Kind != entity.POST_KIND_LIMIT_UNKNOWN {
			language := "zh-cn"
			entry, err := mysql.GetEntryLatestByUser(p.UserId)
			if err == nil && entry != nil {
				language = entry.Language
			}
			content := utils.GetLanguage(language, "topic_message_post_point") + p.Title
			message := CreateTopicMessage(uid, cid, p.UserId, p.CoupleId, entity.TOPIC_MESSAGE_KIND_POST_BE_POINT, content, p.Id)
			AddTopicMessage(language, message)
		}
		// topicInfo
		TopicInfoUpdatePoint(p.Kind, pp.Status >= entity.STATUS_VISIBLE)
	}()
	return pp, err
}

// IsPostPointByUserCouple
func IsPostPointByUserCouple(uid, cid, pid int64) bool {
	if uid <= 0 || cid <= 0 || pid <= 0 {
		return false
	}
	point, _ := mysql.GetPostPointByUserCouple(uid, cid, pid)
	if point == nil || point.Id <= 0 {
		return false
	} else if point.Status < entity.STATUS_VISIBLE {
		return false
	}
	return true
}
