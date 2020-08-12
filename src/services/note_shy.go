package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// AddShy
func AddShy(uid, cid int64, s *entity.Shy) (*entity.Shy, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if s == nil {
		return nil, errors.New("nil_shy")
	} else if s.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len([]rune(s.Safe)) > GetLimit().ShySafeLength {
		return nil, errors.New("limit_content_text_over")
	} else if len([]rune(s.Desc)) > GetLimit().ShyDescLength {
		return nil, errors.New("limit_content_text_over")
	}
	if s.EndAt == 0 {
		// 兼容旧版本，自动+10分钟
		s.EndAt = s.HappenAt + 60*10
	} else if s.HappenAt > s.EndAt {
		return nil, errors.New("limit_happen_err")
	}
	// date
	happen := utils.GetCSTDateByUnix(s.HappenAt)
	s.Year = happen.Year()
	s.MonthOfYear = int(happen.Month())
	s.DayOfMonth = happen.Day()
	// limit
	dayCount := mysql.GetShyTotalByUserCoupleDate(uid, cid, s.Year, s.MonthOfYear, s.DayOfMonth)
	if int(dayCount) >= GetLimit().ShyMaxPerDay {
		return nil, errors.New("limit_total_over")
	}
	// mysql
	s.UserId = uid
	s.CoupleId = cid
	s, err := mysql.AddShy(s)
	if s == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_SHY, s.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, s.Id, "push_title_note_update", "push_content_shy_add", entity.PUSH_TYPE_NOTE_SHY)
	}()
	return s, err
}

// DelShy
func DelShy(uid, cid, sid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if sid <= 0 {
		return errors.New("nil_shy")
	}
	// 旧数据检查
	s, err := mysql.GetShyById(sid)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("nil_shy")
	} else if s.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelShy(s)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_SHY, sid)
		AddTrends(trends)
	}()
	return err
}

// GetShyListByCoupleYearMonth
func GetShyListByCoupleYearMonth(uid, cid int64, year, month int) ([]*entity.Shy, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if year <= 0 {
		return nil, errors.New("limit_year_nil")
	}
	// mysql
	limit := GetPageSizeLimit().Shy
	list, err := mysql.GetShyListByCoupleYearMonth(cid, year, month, 0, limit)
	if err != nil {
		return nil, err
	}
	// 不要了 防止toast
	//if list == nil || len(list) <= 0 {
	//	return nil, errors.New("no_data_shy")
	//}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_SHY)
		AddTrends(trends)
	}()
	return list, err
}
