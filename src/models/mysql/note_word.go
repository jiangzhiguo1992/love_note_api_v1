package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddWord
func AddWord(w *entity.Word) (*entity.Word, error) {
	w.Status = entity.STATUS_VISIBLE
	w.CreateAt = time.Now().Unix()
	w.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_WORD).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,content_text=?").
		Exec(w.Status, w.CreateAt, w.UpdateAt, w.UserId, w.CoupleId, w.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	w.Id, _ = db.Result().LastInsertId()
	return w, nil
}

// DelWord
func DelWord(w *entity.Word) error {
	w.Status = entity.STATUS_DELETE
	w.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_WORD).
		Set("status=?,update_at=?").
		Where("id=?").
		Exec(w.Status, w.UpdateAt, w.Id)
	defer db.Close()
	if db.Err() != nil {
		return errors.New("db_delete_fail")
	}
	return nil
}

//GetWordById
func GetWordById(wid int64) (*entity.Word, error) {
	var w entity.Word
	db := mysqlDB().
		Select("id,status,create_at,update_at,user_id,couple_id,content_text").
		Form(TABLE_WORD).
		Where("id=?").
		Query(wid).
		NextScan(&w.Id, &w.Status, &w.CreateAt, &w.UpdateAt, &w.UserId, &w.CoupleId, &w.ContentText)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if w.Id <= 0 {
		return nil, nil
	} else if w.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &w, nil
}

// GetWordListByCouple
func GetWordListByCouple(cid int64, offset, limit int) ([]*entity.Word, error) {
	db := mysqlDB().
		Select("id,create_at,update_at,user_id,content_text").
		Form(TABLE_WORD).
		Where("status>=? AND couple_id=?").
		OrderDown("create_at").
		Limit(offset, limit).
		Query(entity.STATUS_VISIBLE, cid)
	defer db.Close()
	list := make([]*entity.Word, 0)
	for db.Next() {
		var w entity.Word
		w.CoupleId = cid
		db.Scan(&w.Id, &w.CreateAt, &w.UpdateAt, &w.UserId, &w.ContentText)
		if db.Err() != nil {
			return nil, errors.New("db_query_fail")
		}
		list = append(list, &w)
	}
	return list, nil
}

// GetWordTotalByCouple
func GetWordTotalByCouple(cid int64) int64 {
	var total int64 = 0
	db := mysqlDB().
		Select(SQL_COUNT).
		Form(TABLE_WORD).
		Where("status>=? AND couple_id=?").
		Query(entity.STATUS_VISIBLE, cid).
		NextScan(&total)
	defer db.Close()
	return total
}
