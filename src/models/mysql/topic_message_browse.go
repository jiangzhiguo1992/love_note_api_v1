package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTopicMessageBrowse
func AddTopicMessageBrowse(tmb *entity.TopicMessageBrowse) (*entity.TopicMessageBrowse, error) {
	tmb.Status = entity.STATUS_VISIBLE
	tmb.CreateAt = time.Now().Unix()
	tmb.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TOPIC_MESSAGE_BROWSE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?").
		Exec(tmb.Status, tmb.CreateAt, tmb.UpdateAt, tmb.UserId, tmb.CoupleId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	tmb.Id, _ = db.Result().LastInsertId()
	return tmb, nil
}

// UpdateTopicMessageBrowse
func UpdateTopicMessageBrowse(tmb *entity.TopicMessageBrowse) (*entity.TopicMessageBrowse, error) {
	tmb.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_TOPIC_MESSAGE_BROWSE).
		Set("update_at=?").
		Where("id=?").
		Exec(tmb.UpdateAt, tmb.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return tmb, nil
}

// GetTopicMessageBrowseByUserCouple
func GetTopicMessageBrowseByUserCouple(uid, cid int64) (*entity.TopicMessageBrowse, error) {
	var tmb entity.TopicMessageBrowse
	tmb.UserId = uid
	tmb.CoupleId = cid
	db := mysqlDB().
		Select("id,create_at,update_at").
		Form(TABLE_TOPIC_MESSAGE_BROWSE).
		Where("status>=? AND user_id=? AND couple_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&tmb.Id, &tmb.CreateAt, &tmb.UpdateAt)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if tmb.Id <= 0 {
		return nil, nil
	}
	return &tmb, nil
}
