package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
	"time"
)

// AddMatchPeriod
func AddMatchPeriod(mp *entity.MatchPeriod) (*entity.MatchPeriod, error) {
	if mp == nil {
		return nil, errors.New("nil_period")
	} else if mp.StartAt == 0 || mp.EndAt == 0 {
		return nil, errors.New("limit_happen_err")
	} else if len(strings.TrimSpace(mp.Title)) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if mp.Kind != entity.MATCH_KIND_WIFE_PICTURE &&
		mp.Kind != entity.MATCH_KIND_LETTER_SHOW &&
		mp.Kind != entity.MATCH_KIND_DISCUSS_MEET {
		return nil, errors.New("limit_kind_nil")
	} else if mp.CoinChange <= 0 {
		return nil, errors.New("nil_coin")
	}
	// period
	latest, err := mysql.GetMatchPeriodLatest(mp.Kind)
	if err != nil {
		return nil, err
	}
	if latest == nil {
		mp.Period = 1
	} else {
		mp.Period = latest.Period + 1
	}
	// mysql
	mp, err = mysql.AddMatchPeriod(mp)
	if mp == nil || err != nil {
		return mp, err
	}
	return mp, err
}

// DelMatchPeriod
func DelMatchPeriod(mpid int64) error {
	if mpid <= 0 {
		return errors.New("nil_period")
	}
	// period检查
	mp, err := mysql.GetMatchPeriodById(mpid)
	if err != nil {
		return err
	} else if mp == nil {
		return errors.New("nil_period")
	}
	// mysql
	err = mysql.DelMatchPeriod(mp)
	if err != nil {
		return err
	}
	return err
}

// UpdateMatchPeriodCount
func UpdateMatchPeriodCount(mp *entity.MatchPeriod) (*entity.MatchPeriod, error) {
	if mp == nil || mp.Id <= 0 {
		return nil, errors.New("nil_period")
	} else if !IsMatchPeriodPlay(mp) {
		return nil, errors.New("period_at_err")
	}
	// mysql
	mp, err := mysql.UpdateMatchPeriodCount(mp)
	return mp, err
}

// GetMatchPeriodNow
func GetMatchPeriodNow(kind int) (*entity.MatchPeriod, error) {
	if kind != entity.MATCH_KIND_WIFE_PICTURE &&
		kind != entity.MATCH_KIND_LETTER_SHOW &&
		kind != entity.MATCH_KIND_DISCUSS_MEET {
		return nil, errors.New("limit_kind_nil")
	}
	// mysql
	mp, err := mysql.GetMatchPeriodNow(kind)
	return mp, err
}

// GetMatchPeriodList
func GetMatchPeriodList(kind, page int) ([]*entity.MatchPeriod, error) {
	if kind != entity.MATCH_KIND_WIFE_PICTURE &&
		kind != entity.MATCH_KIND_LETTER_SHOW &&
		kind != entity.MATCH_KIND_DISCUSS_MEET {
		return nil, errors.New("limit_kind_nil")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().MatchPeriod
	offset := page * limit
	list, err := mysql.GetMatchPeriodList(kind, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_match_period")
		} else {
			return nil, nil
		}
	}
	// 没有额外数据和同步
	return list, err
}

// IsMatchPeriodPlay
func IsMatchPeriodPlay(mp *entity.MatchPeriod) bool {
	if mp == nil || mp.StartAt == 0 || mp.EndAt == 0 {
		return false
	}
	now := time.Now().Unix()
	if mp.StartAt > now || mp.EndAt < now {
		return false
	}
	return true
}
