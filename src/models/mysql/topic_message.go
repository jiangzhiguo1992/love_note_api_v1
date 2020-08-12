package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddTopicMessage
func AddTopicMessage(tm *entity.TopicMessage) (*entity.TopicMessage, error) {
	tm.Status = entity.STATUS_VISIBLE
	tm.CreateAt = time.Now().Unix()
	tm.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_TOPIC_MESSAGE).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,to_user_id=?,to_couple_id=?,kind=?,content_text=?,content_id=?").
		Exec(tm.Status, tm.CreateAt, tm.UpdateAt, tm.UserId, tm.CoupleId, tm.ToUserId, tm.ToCoupleId, tm.Kind, tm.ContentText, tm.ContentId)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	tm.Id, _ = db.Result().LastInsertId()
	return tm, nil
}

// GetTopicMessageListByToUserCoupleKind
func GetTopicMessageListByToUserCoupleKind(uid, cid int64, kind, offset, limit int) ([]*entity.TopicMessage, error) {
	hasKind := kind > entity.TOPIC_MESSAGE_KIND_ALL
	var where string
	if !hasKind {
		where = "status>=? AND to_user_id=? AND to_couple_id=?"
	} else {
		where = "status>=? AND to_user_id=? AND to_couple_id=? AND kind=?"
	}
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,couple_id,kind,content_text,content_id").
		Form(TABLE_TOPIC_MESSAGE).
		Where(where).
		OrderDown("create_at").
		Limit(offset, limit)
	if !hasKind {
		db.Query(entity.STATUS_VISIBLE, uid, cid)
	} else {
		db.Query(entity.STATUS_VISIBLE, uid, cid, kind)
	}
	defer db.Close()
	list := make([]*entity.TopicMessage, 0)
	for db.Next() {
		var tm entity.TopicMessage
		tm.ToUserId = uid
		tm.ToCoupleId = cid
		db.Scan(&tm.Id, &tm.CreateAt, &tm.UpdateAt, &tm.UserId, &tm.CoupleId, &tm.Kind, &tm.ContentText, &tm.ContentId)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &tm)
	}
	return list, nil
}

// GetTopicMessageTotalByUpdateToUserCouple
func GetTopicMessageTotalByUpdateToUserCouple(update, tuid, tcid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_TOPIC_MESSAGE).
		Where("status>=? AND update_at>=? AND to_user_id=? AND to_couple_id=?").
		Query(entity.STATUS_VISIBLE, update, tuid, tcid).
		NextScan(&total)
	defer db.Close()
	return total
}
