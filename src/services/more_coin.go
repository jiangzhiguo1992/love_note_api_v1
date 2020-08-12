package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"time"
)

// AddCoinByPay
func AddCoinByPay(uid, cid, bid int64, c *entity.Coin) (*entity.Coin, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if c == nil || c.Change == 0 {
		return nil, errors.New("nil_coin")
	} else if bid <= 0 {
		return nil, errors.New("nil_bill")
	}
	// bill
	b, err := mysql.GetBillById(bid)
	if err != nil {
		return nil, err
	} else if b.UserId != uid {
		return nil, errors.New("db_add_refuse")
	} else if b.CoupleId != cid {
		return nil, errors.New("db_add_refuse")
	}
	// mysql
	c.UserId = uid
	c.CoupleId = cid
	c.Kind = entity.COIN_KIND_ADD_BY_PLAY_PAY
	c.Count = GetCoinCountNowByCouple(cid) + c.Change
	c, err = mysql.AddCoin(c)
	return c, err
}

// AddCoinByFree
func AddCoinByFree(uid, cid int64, c *entity.Coin) (*entity.Coin, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if c == nil || c.Change == 0 {
		return nil, errors.New("nil_coin")
	} else if c.Kind == 0 {
		return nil, errors.New("limit_kind_nil")
	}
	// 检查余额
	c.Count = GetCoinCountNowByCouple(cid) + c.Change
	if c.Count < 0 {
		return nil, errors.New("coin_need")
	}
	// mysql
	c.UserId = uid
	c.CoupleId = cid
	//c.BillId = 0
	c, err := mysql.AddCoin(c)
	return c, err
}

// AddCoinByAdmin
func AddCoinByAdmin(c *entity.Coin) (*entity.Coin, error) {
	if c == nil || c.Change == 0 {
		return nil, errors.New("nil_coin")
	} else if c.UserId <= 0 {
		return nil, errors.New("nil_user")
	} else if c.CoupleId <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	c.Kind = entity.COIN_KIND_ADD_BY_SYS
	//c.BillId = 0
	c.Count = GetCoinCountNowByCouple(c.CoupleId) + c.Change
	c, err := mysql.AddCoin(c)
	return c, err
}

// GetCoinLatest
func GetCoinLatest(cid int64) (*entity.Coin, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	c, err := mysql.GetCoinLatest(cid)
	return c, err
}

// GetCoinList
func GetCoinList(uid, cid, bid int64, kind int, page int) ([]*entity.Coin, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Coin
	offset := page * limit
	list, err := mysql.GetCoinList(uid, cid, bid, kind, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_coin")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetCoinListByCouple
func GetCoinListByCouple(cid int64, page int) ([]*entity.Coin, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Coin
	offset := page * limit
	list, err := mysql.GetCoinListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_coin")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性和同步
	return list, err
}

// AddCoinByAd
func AddCoinByAd(uid, cid int64, kind int) (*entity.Coin, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if kind != entity.COIN_KIND_ADD_BY_AD_WATCH && kind != entity.COIN_KIND_ADD_BY_AD_CLICK {
		return nil, errors.New("limit_kind_nil")
	}
	now := time.Now()
	limit := GetLimit()
	// old
	oldWatch, _ := mysql.GetCoinLatestByUserCoupleKind(uid, cid, entity.COIN_KIND_ADD_BY_AD_WATCH)
	if oldWatch != nil && (oldWatch.CreateAt+int64(limit.CoinAdBetweenSec) > now.Unix()) {
		return nil, errors.New("coin_ad_frequent")
	}
	oldClick, _ := mysql.GetCoinLatestByUserCoupleKind(uid, cid, entity.COIN_KIND_ADD_BY_AD_CLICK)
	if oldClick != nil && (oldClick.CreateAt+int64(limit.CoinAdBetweenSec) > now.Unix()) {
		return nil, errors.New("coin_ad_frequent")
	}
	// count
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	countWatch := mysql.GetCoinCountByCreateUserCoupleKind(today, uid, cid, entity.COIN_KIND_ADD_BY_AD_WATCH)
	countClick := mysql.GetCoinCountByCreateUserCoupleKind(today, uid, cid, entity.COIN_KIND_ADD_BY_AD_CLICK)
	currentCount := countWatch + countClick
	// limit
	coinMax := limit.CoinAdMaxPerDayCount
	if currentCount >= coinMax {
		return nil, errors.New("limit_total_over")
	}
	canGetCount := coinMax - currentCount
	var change = limit.CoinAdWatchCount
	if kind == entity.COIN_KIND_ADD_BY_AD_CLICK {
		change = limit.CoinAdClickCount
	}
	if change > canGetCount {
		change = canGetCount
	}
	// add
	coin := &entity.Coin{
		Kind:   kind,
		Change: change,
	}
	c, err := AddCoinByFree(uid, cid, coin)
	return c, err
}

// GetCoinCountNowByCouple
func GetCoinCountNowByCouple(cid int64) int {
	if cid <= 0 {
		return 0
	}
	coin, _ := mysql.GetCoinLatest(cid)
	if coin == nil {
		return 0
	}
	return coin.Count
}

// GetCoinChangeSign
func GetCoinChangeSign(days int) int {
	limit := GetLimit()
	change := limit.CoinSignMinCount
	max := limit.CoinSignMaxCount
	if days > 1 { // 有连续才叠加
		increase := (days - 1) * limit.CoinSignIncreaseCount
		change += increase
	}
	if change > max {
		change = max
	}
	return change
}

// GetCoinChangeListByCreateWithPay
func GetCoinChangeListByCreateWithPay(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	// mysql
	list, err := mysql.GetCoinChangeListByCreateWithPay(start, end)
	return list, err
}

// GetCoinTotalByCreateKindWithDel
func GetCoinTotalByCreateKindWithDel(start, end int64, kind int) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetCoinTotalByCreateKindWithDel(start, end, kind)
	return total
}
