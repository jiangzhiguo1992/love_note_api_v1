package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"time"
)

// AddSleep
func AddSleep(uid, cid int64, s *entity.Sleep) (*entity.Sleep, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if s == nil {
		return nil, errors.New("nil_sleep")
	}
	// 数据检查
	nowUnix := time.Now().Unix()
	latest, _ := mysql.GetSleepLatestByUserCouple(uid, cid)
	if latest == nil {
		s.IsSleep = true
	} else {
		s.IsSleep = !latest.IsSleep
		if !s.IsSleep {
			// 醒了，时长检查
			if nowUnix < latest.CreateAt+GetLimit().NoteSleepSuccessMinSec {
				return nil, errors.New("sleep_continue_small")
			} else if nowUnix > latest.CreateAt+GetLimit().NoteSleepSuccessMaxSec {
				return nil, errors.New("sleep_continue_large")
			}
		}
	}
	// date
	now := utils.GetCSTDateByUnix(nowUnix)
	s.Year = now.Year()
	s.MonthOfYear = int(now.Month())
	s.DayOfMonth = now.Day()
	// limit
	dayCount := mysql.GetSleepTotalByUserCoupleDate(uid, cid, s.Year, s.MonthOfYear, s.DayOfMonth)
	if int(dayCount) >= GetLimit().SleepMaxPerDay {
		return nil, errors.New("limit_total_over")
	}
	// mysql
	s.UserId = uid
	s.CoupleId = cid
	s, err := mysql.AddSleep(s)
	if s == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_SLEEP, s.Id)
		AddTrends(trends)
		// push
		if s.IsSleep {
			AddPushInCouple(uid, s.Id, "push_title_note_update", "push_content_sleep_add", entity.PUSH_TYPE_NOTE_SLEEP)
		}
	}()
	return s, err
}

// DelSleep
func DelSleep(uid, cid, sid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if sid <= 0 {
		return errors.New("nil_sleep")
	}
	// 旧数据检查
	s, err := mysql.GetSleepById(sid)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("nil_sleep")
	} else if s.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelSleep(s)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_SLEEP, sid)
		AddTrends(trends)
	}()
	return err
}

// GetSleepLatestByUserCouple
func GetSleepLatestByUserCouple(uid, cid int64) (*entity.Sleep, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	s, err := mysql.GetSleepLatestByUserCouple(uid, cid)
	return s, err
}

// GetSleepListByCoupleYearMonth
func GetSleepListByCoupleYearMonth(uid, cid int64, year, month int) ([]*entity.Sleep, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if year <= 0 {
		return nil, errors.New("limit_year_nil")
	}
	// mysql
	limit := GetPageSizeLimit().Sleep
	list, err := mysql.GetSleepListByCoupleYearMonth(cid, year, month, 0, limit)
	if err != nil {
		return nil, err
	}
	// 不要了 防止toast
	//if list == nil || len(list) <= 0 {
	//	return nil, errors.New("no_data_sleep")
	//}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_SLEEP)
		AddTrends(trends)
	}()
	return list, err
}
