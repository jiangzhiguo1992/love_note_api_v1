package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddMovie
func AddMovie(uid, cid int64, m *entity.Movie) (*entity.Movie, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if m == nil {
		return nil, errors.New("nil_movie")
	} else if m.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(m.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(m.Title)) > GetLimit().MovieTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len([]rune(m.ContentText)) > GetLimit().MovieContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// limit
	if len(m.ContentImageList) > 0 {
		imgLimit := GetVipLimitByCouple(cid).MovieImageCount
		if imgLimit <= 0 {
			return nil, errors.New("limit_content_image_refuse")
		} else if len(m.ContentImageList) > imgLimit {
			return nil, errors.New("limit_content_image_over")
		}
	}
	// mysql
	m.UserId = uid
	m.CoupleId = cid
	m, err := mysql.AddMovie(m)
	if m == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_MOVIE, m.Id)
		_, _ = AddTrends(trends)
		// push
		AddPushInCouple(uid, m.Id, "push_title_note_update", m.Title, entity.PUSH_TYPE_NOTE_MOVIE)
	}()
	return m, err
}

// DelMovie
func DelMovie(uid, cid, mid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if mid <= 0 {
		return errors.New("nil_movie")
	}
	// 旧数据检查
	m, err := mysql.GetMovieById(mid)
	if err != nil {
		return err
	} else if m == nil {
		return errors.New("nil_movie")
	} else if m.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelMovie(m)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_MOVIE, mid)
		AddTrends(trends)
	}()
	return err
}

// UpdateMovie
func UpdateMovie(uid, cid int64, m *entity.Movie) (*entity.Movie, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if m == nil || m.Id <= 0 {
		return nil, errors.New("nil_movie")
	} else if m.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(m.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(m.Title)) > GetLimit().MovieTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len([]rune(m.ContentText)) > GetLimit().MovieContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 旧数据检查
	old, err := mysql.GetMovieById(m.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_movie")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// 图片检查
	limit := GetVipLimitByCouple(cid).MovieImageCount
	if (len(m.ContentImageList) > limit) && (len(m.ContentImageList) > len(old.ContentImageList)) {
		// 修改的图数大于限制图数，如果是以前vip传上去的，则通过
		return old, errors.New("limit_content_image_over")
	}
	// mysql
	old.HappenAt = m.HappenAt
	old.Title = m.Title
	old.ContentImageList = m.ContentImageList
	old.ContentText = m.ContentText
	old.Longitude = m.Longitude
	old.Latitude = m.Latitude
	old.Address = m.Address
	old.CityId = m.CityId
	m, err = mysql.UpdateMovie(old)
	if m == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_MOVIE, m.Id)
		AddTrends(trends)
	}()
	return m, err
}

// GetMovieListByCouple
func GetMovieListByCouple(uid, cid int64, page int) ([]*entity.Movie, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Movie
	offset := page * limit
	list, err := mysql.GetMovieListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_movie")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_MOVIE)
		AddTrends(trends)
	}()
	return list, err
}
