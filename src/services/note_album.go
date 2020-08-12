package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddAlbum
func AddAlbum(uid, cid int64, a *entity.Album) (*entity.Album, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil {
		return nil, errors.New("nil_album")
	} else if len(a.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(a.Title)) > GetLimit().AlbumTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// mysql
	a.UserId = uid
	a.CoupleId = cid
	a, err := mysql.AddAlbum(a)
	if a == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_ALBUM, a.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, a.Id, "push_title_note_update", a.Title, entity.PUSH_TYPE_NOTE_ALBUM)
	}()
	return a, err
}

// DelAlbum
func DelAlbum(uid, cid, aid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if aid <= 0 {
		return errors.New("nil_album")
	}
	// 旧数据检查
	a, err := mysql.GetAlbumById(aid)
	if err != nil {
		return err
	} else if a == nil {
		return errors.New("nil_album")
	} else if a.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// 图片检查，不要检查了，可以删除不是空的相册
	if mysql.GetPictureTotalByCoupleAlbum(cid, aid) > 0 {
		return errors.New("album_del_refuse_with_pic")
	}
	// mysql
	err = mysql.DelAlbum(a)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_ALBUM, aid)
		AddTrends(trends)
	}()
	return err
}

// UpdateAlbum
// 1.更新失败返回原数据
func UpdateAlbum(uid, cid int64, a *entity.Album, self bool) (*entity.Album, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil || a.Id <= 0 {
		return nil, errors.New("nil_album")
	} else if len(a.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(a.Title)) > GetLimit().AlbumTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// 旧数据检查
	old, err := mysql.GetAlbumById(a.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_album")
	} else if self && old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// mysql
	old.Title = a.Title
	old.Cover = a.Cover
	if ps, _ := mysql.GetPictureStartByAlbum(a.Id); ps != nil {
		old.StartAt = ps.HappenAt
	} else {
		old.StartAt = 0
	}
	if pe, _ := mysql.GetPictureEndByAlbum(a.Id); pe != nil {
		old.EndAt = pe.HappenAt
	} else {
		old.EndAt = 0
	}
	old.PictureCount = int(mysql.GetPictureTotalByCoupleAlbum(cid, a.Id))
	a, err = mysql.UpdateAlbum(old)
	if a == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_ALBUM, a.Id)
		AddTrends(trends)
	}()
	return a, err
}

// GetAlbumById
func GetAlbumById(uid, cid, aid int64) (*entity.Album, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if aid <= 0 {
		return nil, errors.New("nil_album")
	}
	// mysql
	a, err := mysql.GetAlbumById(aid)
	if err != nil {
		return nil, err
	} else if a == nil {
		return nil, errors.New("nil_album")
	} else if a.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// 同步，picture里已经做了这个了
	//go func() {
	//	trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_ALBUM, aid)
	//	AddTrends(trends)
	//}()
	return a, err
}

// GetAlbumListByCouple
func GetAlbumListByCouple(uid, cid int64, page int) ([]*entity.Album, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Album
	offset := page * limit
	list, err := mysql.GetAlbumListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_album")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_ALBUM)
		AddTrends(trends)
	}()
	return list, err
}
