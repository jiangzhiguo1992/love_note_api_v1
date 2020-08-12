package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddMensesLength
func AddMensesLength(uid, cid int64, ml *entity.MensesLength) (*entity.MensesLength, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if ml == nil {
		return nil, errors.New("nil_menses_length")
	} else if ml.CycleDay <= 0 || ml.CycleDay > GetLimit().MensesMaxCycleDay {
		return nil, errors.New("limit_day_err")
	} else if ml.DurationDay <= 0 || ml.DurationDay > GetLimit().MensesMaxDurationDay {
		return nil, errors.New("limit_day_err")
	}
	// old
	old, err := mysql.GetMensesLengthByUserCouple(uid, cid)
	if err != nil {
		return old, err
	} else if old == nil || old.Id <= 0 {
		// 无记录
		ml.UserId = uid
		ml.CoupleId = cid
		ml, err = mysql.AddMensesLength(ml)
	} else {
		// 有记录
		old.CycleDay = ml.CycleDay
		old.DurationDay = ml.DurationDay
		ml, err = mysql.UpdateMensesLength(ml)
	}
	if ml == nil || err != nil {
		return old, err
	}
	return ml, err
}

// GetMensesLengthByUserCouple
func GetMensesLengthByUserCouple(uid, cid int64) (*entity.MensesLength, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	length, _ := mysql.GetMensesLengthByUserCouple(uid, cid)
	// default 不查错，不判空
	if length == nil {
		length = &entity.MensesLength{
			BaseObj: entity.BaseObj{
				Status: entity.STATUS_VISIBLE,
			},
			BaseCp: entity.BaseCp{
				UserId:   uid,
				CoupleId: cid,
			},
			CycleDay:    GetLimit().MensesDefaultCycleDay,
			DurationDay: GetLimit().MensesDefaultDurationDay,
		}
	}
	return length, nil
}
