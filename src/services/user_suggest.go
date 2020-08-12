package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddSuggest
func AddSuggest(uid int64, s *entity.Suggest) (*entity.Suggest, error) {
	if s == nil {
		return s, errors.New("nil_suggest")
	} else if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if s.Kind <= 0 {
		return nil, errors.New("limit_kind_nil")
	} else if len(strings.TrimSpace(s.Title)) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(s.Title)) > GetLimit().SuggestTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len(strings.TrimSpace(s.ContentText)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(s.ContentText)) > GetLimit().SuggestContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// admin
	u, _ := GetUserById(uid)
	s.Official = IsAdminister(u)
	// mysql
	s.UserId = uid
	s, err := mysql.AddSuggest(s)
	return s, err
}

// DelSuggest
func DelSuggest(uid, sid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if sid <= 0 {
		return errors.New("nil_suggest")
	}
	// 旧数据检查
	s, err := mysql.GetSuggestById(sid)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("nil_suggest")
	}
	// admin
	u, _ := GetUserById(uid)
	if !IsAdminister(u) && s.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelSuggest(s)
	return err
}

// UpdateSuggest
func UpdateSuggest(s *entity.Suggest) (*entity.Suggest, error) {
	if s == nil || s.Id <= 0 {
		return nil, errors.New("nil_suggest")
	}
	// mysql
	s, err := mysql.UpdateSuggest(s)
	return s, err
}

// GetSuggestByIdWithAll
func GetSuggestByIdWithAll(uid, sid int64) (*entity.Suggest, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if sid <= 0 {
		return nil, errors.New("nil_suggest")
	}
	// mysql
	s, err := mysql.GetSuggestById(sid)
	if err != nil {
		return nil, err
	} else if s == nil {
		return nil, errors.New("nil_suggest")
	}
	LoadSuggestWithAll(uid, s)
	return s, err
}

// GetSuggestListByStatusKind
func GetSuggestListByStatusKind(uid int64, status, kind, page int) ([]*entity.Suggest, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Suggest
	offset := page * limit
	list, err := mysql.GetSuggestListByStatusKind(status, kind, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_suggest")
		} else {
			return nil, nil
		}
	}
	// 额外属性
	//for _, s := range list {
	//	LoadSuggestWithAll(uid, s)
	//}
	return list, nil
}

// GetSuggestListByUser
func GetSuggestListByUser(uid int64, page int) ([]*entity.Suggest, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Suggest
	offset := page * limit
	list, err := mysql.GetSuggestListByUser(uid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_suggest")
		} else {
			return nil, nil
		}
	}
	// 额外属性
	//for _, s := range list {
	//	LoadSuggestWithAll(uid, s)
	//}
	return list, nil
}

// GetSuggestListByUserFollow
func GetSuggestListByUserFollow(uid int64, page int) ([]*entity.Suggest, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	followList, err := GetSuggestFollowListByUser(uid, page)
	if err != nil {
		return nil, err
	} else if followList == nil || len(followList) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_suggest")
		} else {
			return nil, nil
		}
	}
	list := make([]*entity.Suggest, 0)
	for _, v := range followList {
		if v == nil {
			continue
		}
		s, err := mysql.GetSuggestById(v.SuggestId)
		if s != nil && s.Id > 0 && s.Status >= entity.STATUS_VISIBLE {
			// 可能会被删除
			//LoadSuggestWithAll(uid, s)
			list = append(list, s)
		} else if err == nil {
			// 不是查询出错，就取消关注
			ToggleSuggestFollow(v)
		}
	}
	if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_suggest")
		} else {
			return nil, nil
		}
	}
	return list, nil
}

// GetSuggestTotalByCreateWithDel
func GetSuggestTotalByCreateWithDel(create int64) int64 {
	if create == 0 {
		return 0
	}
	// mysql
	total := mysql.GetSuggestTotalByCreateWithDel(create)
	return total
}

// LoadSuggestWithAll
func LoadSuggestWithAll(uid int64, s *entity.Suggest) *entity.Suggest {
	if uid <= 0 || s == nil || s.Id <= 0 {
		return s
	}
	// 额外属性
	s.Mine = (s.UserId == uid) && (uid > 0)
	s.Follow = IsSuggestFollowByUser(uid, s.Id)
	s.Comment = IsSuggestCommentByUser(uid, s.Id)
	return s
}
