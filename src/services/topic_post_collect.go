package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// TogglePostCollect
func TogglePostCollect(uid, cid int64, pc *entity.PostCollect) (*entity.PostCollect, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pc == nil || pc.PostId <= 0 {
		return nil, errors.New("nil_post")
	}
	// post检查
	p, _ := GetPostById(pc.PostId)
	// mysql
	old, err := mysql.GetPostCollectByUserCouple(uid, cid, pc.PostId)
	if err != nil {
		return old, err
	} else if old == nil || old.Id <= 0 {
		// 没收藏
		if p == nil {
			return nil, errors.New("nil_post")
		}
		pc.UserId = uid
		pc.CoupleId = cid
		pc, err = mysql.AddPostCollect(pc)
	} else {
		// 已收藏
		if old.Status >= entity.STATUS_VISIBLE {
			old.Status = entity.STATUS_DELETE
		} else {
			if p == nil {
				return nil, errors.New("nil_post")
			}
			old.Status = entity.STATUS_VISIBLE
		}
		pc, err = mysql.UpdatePostCollect(old)
	}
	if pc == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		// 可能会被删除
		if p != nil {
			// post
			if pc.Status >= entity.STATUS_VISIBLE {
				p.CollectCount = p.CollectCount + 1
			} else {
				p.CollectCount = p.CollectCount - 1
			}
			UpdatePostCount(p, false)
			// message
			if uid != p.UserId &&
				pc.Status >= entity.STATUS_VISIBLE &&
				p.Kind != entity.POST_KIND_LIMIT_UNKNOWN {
				language := "zh-cn"
				entry, err := mysql.GetEntryLatestByUser(p.UserId)
				if err == nil && entry != nil {
					language = entry.Language
				}
				content := utils.GetLanguage(language, "topic_message_post_collect") + p.Title
				message := CreateTopicMessage(uid, cid, p.UserId, p.CoupleId, entity.TOPIC_MESSAGE_KIND_POST_BE_COLLECT, content, p.Id)
				AddTopicMessage(language, message)
			}
			// topicInfo
			TopicInfoUpdateCollect(p.Kind, pc.Status >= entity.STATUS_VISIBLE)
		}
	}()
	return pc, err
}

// IsPostCollectByUserCouple
func IsPostCollectByUserCouple(uid, cid, pid int64) bool {
	if uid <= 0 || cid <= 0 || pid <= 0 {
		return false
	}
	collect, _ := mysql.GetPostCollectByUserCouple(uid, cid, pid)
	if collect == nil || collect.Id <= 0 {
		return false
	} else if collect.Status < entity.STATUS_VISIBLE {
		return false
	}
	return true
}
