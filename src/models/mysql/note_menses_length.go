package mysql

import (
	"errors"
	"models/entity"
	"time"
)

// AddMensesLength
func AddMensesLength(ml *entity.MensesLength) (*entity.MensesLength, error) {
	ml.Status = entity.STATUS_VISIBLE
	ml.CreateAt = time.Now().Unix()
	ml.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Insert(TABLE_MENSES_LENGTH).
		Set("status=?,create_at=?,update_at=?,user_id=?,couple_id=?,cycle_day=?,duration_day=?").
		Exec(ml.Status, ml.CreateAt, ml.UpdateAt, ml.UserId, ml.CoupleId, ml.CycleDay, ml.DurationDay)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_add_fail")
	}
	ml.Id, _ = db.Result().LastInsertId()
	return ml, nil
}

// UpdateMensesLength
func UpdateMensesLength(ml *entity.MensesLength) (*entity.MensesLength, error) {
	ml.UpdateAt = time.Now().Unix()
	db := mysqlDB().
		Update(TABLE_MENSES_LENGTH).
		Set("update_at=?,cycle_day=?,duration_day=?").
		Where("id=?").
		Exec(ml.UpdateAt, ml.CycleDay, ml.DurationDay, ml.Id)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_update_fail")
	}
	return ml, nil
}

// GetMensesLengthByUserCouple
func GetMensesLengthByUserCouple(uid, cid int64) (*entity.MensesLength, error) {
	var ml entity.MensesLength
	ml.UserId = uid
	ml.CoupleId = cid
	db := mysqlDB().
		Select("id,status,create_at,update_at,cycle_day,duration_day").
		Form(TABLE_MENSES_LENGTH).
		Where("status>=? AND user_id=? AND couple_id=?").
		Limit(0, 1).
		Query(entity.STATUS_VISIBLE, uid, cid).
		NextScan(&ml.Id, &ml.Status, &ml.CreateAt, &ml.UpdateAt, &ml.CycleDay, &ml.DurationDay)
	defer db.Close()
	if db.Err() != nil {
		return nil, errors.New("db_query_fail")
	} else if ml.Id <= 0 {
		return nil, nil
	} else if ml.Status < entity.STATUS_VISIBLE {
		return nil, nil
	}
	return &ml, nil
}
