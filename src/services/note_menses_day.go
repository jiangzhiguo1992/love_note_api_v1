package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"time"
)

// AddMensesDay
func AddMensesDay(uid, taId, cid int64, md *entity.MensesDay) (*entity.MensesDay, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if md == nil {
		return nil, errors.New("nil_menses_day")
	} else if md.Year <= 0 || md.MonthOfYear <= 0 || md.DayOfMonth <= 0 {
		return nil, errors.New("limit_happen_err")
	}
	// menses2
	unixDay := utils.GetUnixByCSTDate(time.Date(md.Year, time.Month(md.MonthOfYear), md.DayOfMonth, 0, 0, 0, 0, time.Local))
	m, err := mysql.GetMenses2ByUserCoupleAtBetween(uid, cid, unixDay)
	if err != nil {
		return nil, err
	} else if m == nil {
		m, err = mysql.GetMenses2ByUserCoupleAtBetween(taId, cid, unixDay)
		if err != nil {
			return nil, err
		} else if m == nil {
			return nil, errors.New("nil_menses")
		} else if m.CoupleId != cid {
			return nil, errors.New("db_add_refuse")
		}
	} else if m.CoupleId != cid {
		return nil, errors.New("db_add_refuse")
	}
	// old
	old, err := mysql.GetMensesDayByMensesDate(m.Id, md.Year, md.MonthOfYear, md.DayOfMonth)
	if err != nil {
		return nil, err
	} else if old == nil || old.Id <= 0 {
		// add
		md.UserId = uid
		md.CoupleId = cid
		md.Menses2Id = m.Id
		md, err = mysql.AddMensesDay(md)
	} else {
		// update
		old.Blood = md.Blood
		old.Pain = md.Pain
		old.Mood = md.Mood
		md, err = mysql.UpdateMensesDay(old)
	}
	if md == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_MENSES, m.Id)
		_, _ = AddTrends(trends)
		// push
		AddPushInCouple(uid, m.Id, "push_title_note_update", "push_content_menses_come", entity.PUSH_TYPE_NOTE_MENSES)
	}()
	return md, err
}
