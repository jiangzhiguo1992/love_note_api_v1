package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddPictureList
func AddPictureList(uid, cid int64, list []*entity.Picture) ([]*entity.Picture, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if list == nil || len(list) <= 0 {
		return nil, errors.New("nil_picture")
	}
	// 相册检查
	albumId := list[0].AlbumId
	for _, p := range list {
		if p.AlbumId <= 0 {
			return nil, errors.New("picture_album_nil")
		} else if albumId != p.AlbumId {
			return nil, errors.New("album_no_same")
		}
	}
	a, err := mysql.GetAlbumById(albumId)
	if err != nil {
		return nil, err
	} else if a == nil {
		return nil, errors.New("nil_album")
	} else if a.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// limit
	totalLimit := GetVipLimitByCouple(cid).PictureTotalCount
	if totalLimit <= 0 {
		return nil, errors.New("limit_content_image_refuse")
	} else if mysql.GetPictureTotalByCouple(cid)+int64(len(list)) > int64(totalLimit) {
		return nil, errors.New("limit_content_image_over")
	}
	// 遍历多个图片
	pushList := make([]*entity.Picture, 0)
	for _, p := range list {
		if p == nil {
			continue
		}
		if p.AlbumId <= 0 {
			return nil, errors.New("nil_album")
		} else if p.HappenAt == 0 {
			return nil, errors.New("limit_happen_nil")
		} else if len(strings.TrimSpace(p.ContentImage)) <= 0 {
			return nil, errors.New("limit_content_image_nil")
		}
		// mysql
		picture := &entity.Picture{}
		picture.UserId = uid
		picture.CoupleId = cid
		picture.AlbumId = p.AlbumId
		picture.HappenAt = p.HappenAt
		picture.ContentImage = p.ContentImage
		picture.Longitude = p.Longitude
		picture.Latitude = p.Latitude
		picture.Address = p.Address
		picture.CityId = p.CityId
		picture, _ = mysql.AddPicture(picture)
		if picture != nil {
			pushList = append(pushList, picture)
		}
	}
	// 检查成功数量
	if pushList == nil || len(pushList) <= 0 {
		return nil, errors.New("nil_picture")
	}
	// 同步
	go func() {
		// album
		UpdateAlbum(uid, cid, a, false)
		// 动态
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_ALBUM, a.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, a.Id, "push_title_note_update", "push_content_picture_add", entity.PUSH_TYPE_NOTE_PICTURE)
	}()
	return pushList, err
}

// DelPicture
func DelPicture(uid, cid, pid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if pid <= 0 {
		return errors.New("nil_picture")
	}
	// 数据检查
	p, err := mysql.GetPictureById(pid)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("nil_picture")
	} else if p.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	a, err := mysql.GetAlbumById(p.AlbumId)
	if err != nil {
		return err
	} else if a == nil {
		return errors.New("nil_album")
	} else if a.CoupleId != cid {
		return errors.New("db_query_refuse")
	}
	// mysql
	err = mysql.DelPicture(p)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		// album
		UpdateAlbum(uid, cid, a, false)
		// 动态
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_ALBUM, a.Id)
		AddTrends(trends)
	}()
	return err
}

// UpdatePicture
func UpdatePicture(uid, cid int64, p *entity.Picture) (*entity.Picture, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return p, errors.New("nil_couple")
	} else if p == nil || p.Id <= 0 {
		return nil, errors.New("nil_picture")
	} else if p.AlbumId <= 0 {
		return nil, errors.New("nil_album")
	} else if p.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	}
	// 相册检查
	a, err := mysql.GetAlbumById(p.AlbumId)
	if err != nil {
		return nil, err
	} else if a == nil {
		return nil, errors.New("nil_album")
	} else if a.CoupleId != cid {
		return nil, errors.New("db_update_refuse")
	}
	// 旧数据检查
	old, err := mysql.GetPictureById(p.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_picture")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// mysql
	old.AlbumId = p.AlbumId
	old.HappenAt = p.HappenAt
	old.Longitude = p.Longitude
	old.Latitude = p.Latitude
	old.Address = p.Address
	old.CityId = p.CityId
	p, err = mysql.UpdatePicture(old)
	if p == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		// album
		UpdateAlbum(uid, cid, a, false)
		// 动态
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_ALBUM, a.Id)
		AddTrends(trends)
	}()
	return p, err
}

// GetPictureListByCoupleAlbum
func GetPictureListByCoupleAlbum(uid, cid, aid int64, page int) ([]*entity.Picture, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if aid <= 0 {
		return nil, errors.New("nil_album")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Picture
	offset := page * limit
	list, err := mysql.GetPictureListByCoupleAlbum(cid, aid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_picture")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_ALBUM, aid)
		AddTrends(trends)
	}()
	return list, err
}
