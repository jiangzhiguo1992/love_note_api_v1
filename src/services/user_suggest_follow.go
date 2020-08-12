package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// ToggleSuggestFollow
func ToggleSuggestFollow(sf *entity.SuggestFollow) (*entity.SuggestFollow, error) {
	if sf == nil || sf.SuggestId <= 0 {
		return nil, errors.New("nil_suggest")
	} else if sf.UserId <= 0 {
		return nil, errors.New("nil_user")
	}
	// 数据检查
	s, _ := mysql.GetSuggestById(sf.SuggestId)
	// mysql
	old, err := mysql.GetSuggestFollowByUser(sf.UserId, sf.SuggestId)
	if err != nil {
		return old, err
	} else if old == nil || old.Id <= 0 {
		if s == nil {
			return nil, errors.New("nil_suggest")
		}
		// 没关注
		sf, err = mysql.AddSuggestFollow(sf)
	} else {
		// 已关注
		if old.Status >= entity.STATUS_VISIBLE {
			old.Status = entity.STATUS_DELETE
		} else {
			if s == nil {
				return nil, errors.New("nil_suggest")
			}
			old.Status = entity.STATUS_VISIBLE
		}
		sf, err = mysql.UpdateSuggestFollowStatus(old)
	}
	if sf == nil || err != nil {
		return sf, err
	}
	// 同步
	go func() {
		// 可能会被删除
		if s != nil {
			// post
			if sf.Status >= entity.STATUS_VISIBLE {
				s.FollowCount = s.FollowCount + 1
			} else {
				s.FollowCount = s.FollowCount - 1
			}
			mysql.UpdateSuggestCount(s, false)
		}
	}()
	return sf, err
}

// GetSuggestFollowListByUser
func GetSuggestFollowListByUser(uid int64, page int) ([]*entity.SuggestFollow, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Suggest
	offset := page * limit
	list, err := mysql.GetSuggestFollowListByUser(uid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_suggest")
		} else {
			return nil, nil
		}
	}
	return list, nil
}

// IsSuggestFollowByUser
func IsSuggestFollowByUser(uid, sid int64) bool {
	if uid <= 0 {
		return false
	} else if sid <= 0 {
		return false
	}
	sf, _ := mysql.GetSuggestFollowByUser(uid, sid)
	if sf == nil || sf.Id <= 0 {
		return false
	} else if sf.Status < entity.STATUS_VISIBLE {
		return false
	}
	return true
}
