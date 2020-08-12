package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"time"
)

// AddSign
func AddSign(uid, cid int64) (*entity.Sign, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// 重复性检查
	nowUnix := time.Now().Unix()
	now := utils.GetCSTDateByUnix(nowUnix)
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	repeat, err := mysql.GetSignByCoupleYearMonthDay(cid, year, month, day)
	if err != nil {
		return nil, err
	} else if repeat != nil {
		return nil, errors.New("sign_repeat")
	}
	// 生成参数
	s := &entity.Sign{}
	s.UserId = uid
	s.CoupleId = cid
	// 连续天数
	latestTime := utils.GetCSTDateByUnix(nowUnix - 60*60*24)
	latestYear := latestTime.Year()
	latestMonth := int(latestTime.Month())
	latestDay := latestTime.Day()
	yesterdaySign, _ := mysql.GetSignByCoupleYearMonthDay(cid, latestYear, latestMonth, latestDay)
	if yesterdaySign != nil {
		s.ContinueDay = yesterdaySign.ContinueDay + 1
	} else {
		s.ContinueDay = 1
	}
	// date
	s.Year = now.Year()
	s.MonthOfYear = int(now.Month())
	s.DayOfMonth = now.Day()
	// mysql
	s, err = mysql.AddSign(s)
	if s == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// 金币
		coin := &entity.Coin{
			Kind:   entity.COIN_KIND_ADD_BY_SIGN_DAY,
			Change: GetCoinChangeSign(s.ContinueDay),
		}
		AddCoinByFree(uid, cid, coin)
	}()
	return s, err
}

// GetSignByCoupleYearMonthDay
func GetSignByCoupleYearMonthDay(cid int64, year, month, day int) (*entity.Sign, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if year <= 0 {
		return nil, errors.New("limit_year_nil")
	}
	// mysql
	s, err := mysql.GetSignByCoupleYearMonthDay(cid, year, month, day)
	return s, err
}

// GetSignListByCoupleYearMonth
func GetSignListByCoupleYearMonth(cid int64, year, month int) ([]*entity.Sign, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if year <= 0 {
		return nil, errors.New("limit_year_nil")
	}
	limit := GetPageSizeLimit().Sign
	// mysql
	list, err := mysql.GetSignListByCoupleYearMonth(cid, year, month, 0, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		return nil, errors.New("no_data_sign")
	}
	// 没有额外属性和同步
	return list, err
}

// GetSignList
func GetSignList(uid, cid int64, page int) ([]*entity.Sign, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	limit := GetPageSizeLimit().Sign
	offset := page * limit
	// mysql
	list, err := mysql.GetSignList(uid, cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		return nil, errors.New("no_data_sign")
	}
	// 没有额外属性和同步
	return list, err
}

// GetSignTotalWithDel
func GetSignTotalWithDel(year, month, day int) int64 {
	if year == 0 {
		return 0
	}
	// mysql
	total := mysql.GetSignTotalWithDel(year, month, day)
	return total
}
