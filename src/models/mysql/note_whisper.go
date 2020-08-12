package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddWhisper
func AddWhisper(w *entity.Whisper) (*entity.Whisper, error) {
	w.Status = entity.STATUS_VISIBLE
	w.CreateAt = time.Now().Unix()
	w.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_WHISPER).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,channel=?,is_image=?,content=?").
		Exec(w.Status, w.CreateAt, w.UpdateAt, w.UserId, w.CoupleId, w.Channel, w.IsImage, w.Content)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	w.Id, _ = db.Result().LastInsertId()
	return w, nil
}

// GetWhisperListByCoupleChannel
func GetWhisperListByCoupleChannel(cid int64, channel string, offset, limit int) ([]*entity.Whisper, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,is_image,content").
		Form(TABLE_WHISPER).
		Where("status>=? AND couple_id=? AND channel=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid, channel)
	defer db.Close()
	list := make([]*entity.Whisper, 0)
	for db.Next() {
		var w entity.Whisper
		w.CoupleId = cid
		w.Channel = channel
		db.Scan(&w.Id, &w.CreateAt, &w.UpdateAt, &w.UserId, &w.IsImage, &w.Content)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &w)
	}
	return list, nil
}
