package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"time"
)

// AddMenses
func AddMenses(uid, cid int64, m *entity.Menses) (*entity.Menses, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if m == nil {
		return nil, errors.New("nil_menses")
	}
	// isStart
	latest, _ := mysql.GetMensesLatestByUserCouple(uid, cid)
	if latest == nil {
		m.IsStart = true
	} else {
		m.IsStart = !latest.IsStart
	}
	// date
	now := utils.GetCSTDateByUnix(time.Now().Unix())
	if m.Year <= 0 {
		m.Year = now.Year()
	}
	if m.MonthOfYear <= 0 {
		m.MonthOfYear = int(now.Month())
	}
	if m.DayOfMonth <= 0 {
		m.DayOfMonth = now.Day()
	}
	// mysql
	m.UserId = uid
	m.CoupleId = cid
	m, err := mysql.AddMenses(m)
	if m == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_MENSES, m.Id)
		_, _ = AddTrends(trends)
		// push
		content := "push_content_menses_come"
		if !m.IsStart {
			content = "push_content_menses_gone"
		}
		AddPushInCouple(uid, m.Id, "push_title_note_update", content, entity.PUSH_TYPE_NOTE_MENSES)
	}()
	return m, err
}

// DelMenses
func DelMenses(uid, cid, sid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if sid <= 0 {
		return errors.New("nil_menses")
	}
	// 旧数据检查
	s, err := mysql.GetMensesById(sid)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("nil_menses")
	} else if s.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelMenses(s)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_MENSES, sid)
		AddTrends(trends)
	}()
	return err
}

// DelMensesByDate
func DelMensesByDate(uid, cid int64, year, month, day int) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if year <= 0 || month <= 0 || day <= 0 {
		return errors.New("limit_happen_nil")
	}
	// 数据封装
	s := &entity.Menses{
		BaseCp: entity.BaseCp{
			UserId:   uid,
			CoupleId: cid,
		},
		Year:        year,
		MonthOfYear: month,
		DayOfMonth:  day,
	}
	// mysql
	err := mysql.DelMensesByUserCoupleDate(s)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_MENSES)
		AddTrends(trends)
	}()
	return err
}

// GetMensesLatestByUserCouple
func GetMensesLatestByUserCouple(uid, cid int64) (*entity.Menses, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	m, err := mysql.GetMensesLatestByUserCouple(uid, cid)
	return m, err
}

// GetMensesListByUserCoupleYearMonth
func GetMensesListByUserCoupleYearMonth(uid, cid int64, year, month int) ([]*entity.Menses, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if year <= 0 {
		return nil, errors.New("limit_year_nil")
	}
	// mysql
	limit := GetPageSizeLimit().Menses
	list, err := mysql.GetMensesListByUserCoupleYearMonth(uid, cid, year, month, 0, limit)
	if err != nil {
		return nil, err
	}
	// 不要了 防止toast
	//if list == nil || len(list) <= 0 {
	//	return nil, errors.New("no_data_menses")
	//}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_MENSES)
		AddTrends(trends)
	}()
	return list, err
}
