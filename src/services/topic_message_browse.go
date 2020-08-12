package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddTopicMessageBrowse
func AddTopicMessageBrowse(uid, cid int64) (*entity.TopicMessageBrowse, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// old
	old, err := mysql.GetTopicMessageBrowseByUserCouple(uid, cid)
	if err != nil {
		return nil, err
	} else if old == nil || old.Id <= 0 {
		tmb := &entity.TopicMessageBrowse{
			BaseCp: entity.BaseCp{
				UserId:   uid,
				CoupleId: cid,
			},
		}
		old, err = mysql.AddTopicMessageBrowse(tmb)
	} else {
		old, err = mysql.UpdateTopicMessageBrowse(old)
	}
	if old == nil || err != nil {
		return old, err
	}
	return old, err
}

// GetTopicMessageCountByUserCouple
func GetTopicMessageCountByUserCouple(uid, cid int64) int {
	if uid <= 0 {
		return 0
	} else if cid <= 0 {
		return 0
	}
	// TopicMessageBrowseMe
	tmbMe, err := mysql.GetTopicMessageBrowseByUserCouple(uid, cid)
	if err != nil {
		return 0
	} else if tmbMe == nil {
		return 0
	}
	// mysql
	return int(mysql.GetTopicMessageTotalByUpdateToUserCouple(tmbMe.UpdateAt, uid, cid))
}
