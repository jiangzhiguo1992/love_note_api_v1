package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddTravel
func AddTravel(uid, cid int64, t *entity.Travel) (*entity.Travel, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if t == nil {
		return nil, errors.New("nil_travel")
	} else if t.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(t.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(t.Title)) > GetLimit().TravelTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// mysql
	t.UserId = uid
	t.CoupleId = cid
	t, err := mysql.AddTravel(t)
	if t == nil || err != nil {
		return nil, err
	}
	// 额外属性
	updateTravelForeign(t)
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_TRAVEL, t.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, t.Id, "push_title_note_update", t.Title, entity.PUSH_TYPE_NOTE_TRAVEL)
	}()
	return t, err
}

// DelTravel
func DelTravel(uid, cid, tid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if tid <= 0 {
		return errors.New("nil_travel")
	}
	// 旧数据检查
	t, err := mysql.GetTravelById(tid)
	if err != nil {
		return err
	} else if t == nil {
		return errors.New("nil_travel")
	} else if t.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelTravel(t)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_TRAVEL, tid)
		AddTrends(trends)
	}()
	return err
}

// UpdateTravel
func UpdateTravel(uid, cid int64, t *entity.Travel) (*entity.Travel, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if t == nil || t.Id <= 0 {
		return nil, errors.New("nil_travel")
	} else if t.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(t.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(t.Title)) > GetLimit().TravelTitleLength {
		return nil, errors.New("limit_title_over")
	}
	// 旧数据检查
	old, err := mysql.GetTravelById(t.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_travel")
	} else if old.CoupleId != cid {
		return old, errors.New("db_update_refuse")
	}
	// mysql
	old.HappenAt = t.HappenAt
	old.Title = t.Title
	old.TravelPlaceList = t.TravelPlaceList
	old.TravelAlbumList = t.TravelAlbumList
	old.TravelVideoList = t.TravelVideoList
	old.TravelFoodList = t.TravelFoodList
	old.TravelMovieList = t.TravelMovieList
	old.TravelDiaryList = t.TravelDiaryList
	t, err = mysql.UpdateTravel(old)
	if t == nil || err != nil {
		return old, err
	}
	// 额外属性
	updateTravelForeign(t)
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_TRAVEL, t.Id)
		AddTrends(trends)
	}()
	return t, err
}

// GetTravelByIdWithForeign
func GetTravelByIdWithForeign(uid, cid, tid int64) (*entity.Travel, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if tid <= 0 {
		return nil, errors.New("nil_travel")
	}
	// mysql
	t, err := mysql.GetTravelById(tid)
	if err != nil {
		return nil, err
	} else if t == nil {
		return nil, errors.New("nil_travel")
	} else if t.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// 额外属性
	LoadTravelWithForeign(t)
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_TRAVEL, tid)
		AddTrends(trends)
	}()
	return t, err
}

// GetTravelListByCouple
func GetTravelListByCouple(uid, cid int64, page int) ([]*entity.Travel, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Travel
	offset := page * limit
	list, err := mysql.GetTravelListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_travel")
		} else {
			return nil, nil
		}
	}
	// 额外属性
	//for _, t := range list {
	//	LoadTravelWithPlace(t)
	//}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_TRAVEL)
		AddTrends(trends)
	}()
	return list, err
}

// LoadTravelWithPlace
func LoadTravelWithPlace(t *entity.Travel) *entity.Travel {
	if t == nil || t.Id <= 0 || t.CoupleId <= 0 {
		return t
	}
	t.TravelPlaceList, _ = mysql.GetTravelPlaceListByCoupleTravel(t.CoupleId, t.Id)
	return t
}

// LoadTravelWithForeign
func LoadTravelWithForeign(t *entity.Travel) *entity.Travel {
	if t == nil || t.Id <= 0 || t.CoupleId <= 0 {
		return t
	}
	t.TravelPlaceList, _ = mysql.GetTravelPlaceListByCoupleTravel(t.CoupleId, t.Id)
	t.TravelAlbumList, _ = mysql.GetTravelAlbumListByCoupleTravel(t.CoupleId, t.Id)
	if t.TravelAlbumList != nil && len(t.TravelAlbumList) > 0 {
		for _, v := range t.TravelAlbumList {
			if v == nil || v.AlbumId <= 0 {
				continue
			}
			v.Album, _ = mysql.GetAlbumById(v.AlbumId)
		}
	}
	t.TravelVideoList, _ = mysql.GetTravelVideoListByCoupleTravel(t.CoupleId, t.Id)
	if t.TravelVideoList != nil && len(t.TravelVideoList) > 0 {
		for _, v := range t.TravelVideoList {
			if v == nil || v.VideoId <= 0 {
				continue
			}
			v.Video, _ = mysql.GetVideoById(v.VideoId)
		}
	}
	t.TravelFoodList, _ = mysql.GetTravelFoodListByCoupleTravel(t.CoupleId, t.Id)
	if t.TravelFoodList != nil && len(t.TravelFoodList) > 0 {
		for _, v := range t.TravelFoodList {
			if v == nil || v.FoodId <= 0 {
				continue
			}
			v.Food, _ = mysql.GetFoodById(v.FoodId)
		}
	}
	t.TravelMovieList, _ = mysql.GetTravelMovieListByCoupleTravel(t.CoupleId, t.Id)
	if t.TravelMovieList != nil && len(t.TravelMovieList) > 0 {
		for _, v := range t.TravelMovieList {
			if v == nil || v.MovieId <= 0 {
				continue
			}
			v.Movie, _ = mysql.GetMovieById(v.MovieId)
		}
	}
	t.TravelDiaryList, _ = mysql.GetTravelDiaryListByCoupleTravel(t.CoupleId, t.Id)
	if t.TravelDiaryList != nil && len(t.TravelDiaryList) > 0 {
		for _, v := range t.TravelDiaryList {
			if v == nil || v.DiaryId <= 0 {
				continue
			}
			v.Diary, _ = mysql.GetDiaryById(v.DiaryId)
		}
	}
	return t
}

// updateTravelForeign
func updateTravelForeign(t *entity.Travel) *entity.Travel {
	if t == nil || t.Id <= 0 || t.UserId <= 0 || t.CoupleId <= 0 {
		return t
	}
	// 地址
	addPlaceList, delPlaceList := splitTravelPlaceListByStatus(t.TravelPlaceList)
	if delPlaceList != nil && len(delPlaceList) > 0 {
		DelTravelPlaceList(t.UserId, delPlaceList)
	}
	if addPlaceList != nil && len(addPlaceList) > 0 {
		t.TravelPlaceList = AddTravelPlaceList(t.UserId, t.CoupleId, t.Id, addPlaceList)
	} else {
		t.TravelPlaceList = make([]*entity.TravelPlace, 0)
	}
	// 相册
	addAlbumList, delAlbumList := splitTravelAlbumListByStatus(t.TravelAlbumList)
	if len(delAlbumList) > 0 {
		DelTravelAlbumList(t.UserId, delAlbumList)
	}
	if len(addAlbumList) > 0 {
		t.TravelAlbumList = AddTravelAlbumList(t.UserId, t.CoupleId, t.Id, addAlbumList)
	} else {
		t.TravelAlbumList = make([]*entity.TravelAlbum, 0)
	}
	// 视频
	addVideoList, delVideoList := splitTravelVideoListByStatus(t.TravelVideoList)
	if len(delVideoList) > 0 {
		DelTravelVideoList(t.UserId, delVideoList)
	}
	if len(addVideoList) > 0 {
		t.TravelVideoList = AddTravelVideoList(t.UserId, t.CoupleId, t.Id, addVideoList)
	} else {
		t.TravelVideoList = make([]*entity.TravelVideo, 0)
	}
	// 美食
	addFoodList, delFoodList := splitTravelFoodListByStatus(t.TravelFoodList)
	if len(delFoodList) > 0 {
		DelTravelFoodList(t.UserId, delFoodList)
	}
	if len(addFoodList) > 0 {
		t.TravelFoodList = AddTravelFoodList(t.UserId, t.CoupleId, t.Id, addFoodList)
	} else {
		t.TravelFoodList = make([]*entity.TravelFood, 0)
	}
	// 电影
	addMovieList, delMovieList := splitTravelMovieListByStatus(t.TravelMovieList)
	if len(delMovieList) > 0 {
		DelTravelMovieList(t.UserId, delMovieList)
	}
	if len(addMovieList) > 0 {
		t.TravelMovieList = AddTravelMovieList(t.UserId, t.CoupleId, t.Id, addMovieList)
	} else {
		t.TravelMovieList = make([]*entity.TravelMovie, 0)
	}
	// 日记
	addDiaryList, delDiaryList := splitTravelDiaryListByStatus(t.TravelDiaryList)
	if len(delDiaryList) > 0 {
		DelTravelDiaryList(t.UserId, delDiaryList)
	}
	if len(addDiaryList) > 0 {
		t.TravelDiaryList = AddTravelDiaryList(t.UserId, t.CoupleId, t.Id, addDiaryList)
	} else {
		t.TravelDiaryList = make([]*entity.TravelDiary, 0)
	}
	return t
}

// splitTravelPlaceListByStatus
func splitTravelPlaceListByStatus(objList []*entity.TravelPlace) ([]*entity.TravelPlace, []*entity.TravelPlace) {
	addList := make([]*entity.TravelPlace, 0)
	delList := make([]*entity.TravelPlace, 0)
	if objList == nil || len(objList) <= 0 {
		return addList, delList
	}
	for _, v := range objList {
		if v == nil {
			continue
		}
		if v.Status >= entity.STATUS_VISIBLE {
			addList = append(addList, v)
		} else {
			delList = append(delList, v)
		}
	}
	return addList, delList
}

// splitTravelAlbumListByStatus
func splitTravelAlbumListByStatus(objList []*entity.TravelAlbum) ([]*entity.TravelAlbum, []*entity.TravelAlbum) {
	addList := make([]*entity.TravelAlbum, 0)
	delList := make([]*entity.TravelAlbum, 0)
	if objList == nil || len(objList) <= 0 {
		return addList, delList
	}
	for _, v := range objList {
		if v == nil {
			continue
		}
		if v.Status >= entity.STATUS_VISIBLE {
			addList = append(addList, v)
		} else {
			delList = append(delList, v)
		}
	}
	return addList, delList
}

// splitTravelVideoListByStatus
func splitTravelVideoListByStatus(objList []*entity.TravelVideo) ([]*entity.TravelVideo, []*entity.TravelVideo) {
	addList := make([]*entity.TravelVideo, 0)
	delList := make([]*entity.TravelVideo, 0)
	if objList == nil || len(objList) <= 0 {
		return addList, delList
	}
	for _, v := range objList {
		if v == nil {
			continue
		}
		if v.Status >= entity.STATUS_VISIBLE {
			addList = append(addList, v)
		} else {
			delList = append(delList, v)
		}
	}
	return addList, delList
}

// splitTravelFoodListByStatus
func splitTravelFoodListByStatus(objList []*entity.TravelFood) ([]*entity.TravelFood, []*entity.TravelFood) {
	addList := make([]*entity.TravelFood, 0)
	delList := make([]*entity.TravelFood, 0)
	if objList == nil || len(objList) <= 0 {
		return addList, delList
	}
	for _, v := range objList {
		if v == nil {
			continue
		}
		if v.Status >= entity.STATUS_VISIBLE {
			addList = append(addList, v)
		} else {
			delList = append(delList, v)
		}
	}
	return addList, delList
}

// splitTravelMovieListByStatus
func splitTravelMovieListByStatus(objList []*entity.TravelMovie) ([]*entity.TravelMovie, []*entity.TravelMovie) {
	addList := make([]*entity.TravelMovie, 0)
	delList := make([]*entity.TravelMovie, 0)
	if objList == nil || len(objList) <= 0 {
		return addList, delList
	}
	for _, v := range objList {
		if v == nil {
			continue
		}
		if v.Status >= entity.STATUS_VISIBLE {
			addList = append(addList, v)
		} else {
			delList = append(delList, v)
		}
	}
	return addList, delList
}

// splitTravelDiaryListByStatus
func splitTravelDiaryListByStatus(objList []*entity.TravelDiary) ([]*entity.TravelDiary, []*entity.TravelDiary) {
	addList := make([]*entity.TravelDiary, 0)
	delList := make([]*entity.TravelDiary, 0)
	if objList == nil || len(objList) <= 0 {
		return addList, delList
	}
	for _, v := range objList {
		if v == nil {
			continue
		}
		if v.Status >= entity.STATUS_VISIBLE {
			addList = append(addList, v)
		} else {
			delList = append(delList, v)
		}
	}
	return addList, delList
}

// AddTravelPlaceList
func AddTravelPlaceList(uid, cid, tid int64, tpList []*entity.TravelPlace) []*entity.TravelPlace {
	list := make([]*entity.TravelPlace, 0)
	if tpList == nil || len(tpList) <= 0 {
		return list
	}
	// limit
	limitCount := GetLimit().TravelPlaceCount
	total := mysql.GetTravelPlaceTotalByCoupleTravel(cid, tid)
	if total+len(list) > limitCount {
		return list
	}
	for _, v := range tpList {
		if v == nil {
			continue
		}
		// 旧数据
		if v.Id > 0 {
			tp, err := mysql.GetTravelPlaceById(v.Id)
			if err != nil {
				continue
			} else if tp != nil && tp.CoupleId == cid {
				list = append(list, tp)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.TravelId = tid
		if v.HappenAt == 0 || len(v.ContentText) <= 0 ||
			len([]rune(v.ContentText)) > GetLimit().TravelPlaceContentLength {
			continue
		}
		tp, _ := mysql.AddTravelPlace(v)
		if tp != nil {
			list = append(list, tp)
		}
	}
	return list
}

// DelTravelPlaceList
func DelTravelPlaceList(uid int64, tpList []*entity.TravelPlace) []*entity.TravelPlace {
	list := make([]*entity.TravelPlace, 0)
	if tpList == nil || len(tpList) <= 0 {
		return list
	}
	for _, v := range tpList {
		if v == nil {
			continue
		}
		tp, _ := mysql.GetTravelPlaceById(v.Id)
		if tp == nil || tp.UserId != uid {
			continue
		}
		err := mysql.DelTravelPlace(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddTravelVideoList
func AddTravelVideoList(uid, cid, tid int64, tvList []*entity.TravelVideo) []*entity.TravelVideo {
	list := make([]*entity.TravelVideo, 0)
	if tvList == nil || len(tvList) <= 0 {
		return list
	}
	// limit
	limitCount := GetLimit().TravelVideoCount
	total := mysql.GetTravelVideoTotalByCoupleTravel(cid, tid)
	if total+len(list) > limitCount {
		return list
	}
	for _, v := range tvList {
		if v == nil || v.VideoId <= 0 {
			continue
		}
		// 关联数据
		video, err := mysql.GetVideoById(v.VideoId)
		if video == nil || video.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		tv, err := mysql.GetTravelVideoByCoupleTravelVideo(cid, tid, v.VideoId)
		if err != nil {
			continue
		} else if tv != nil && tv.CoupleId == cid {
			tv.Video = video
			list = append(list, tv)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			tv, err := mysql.GetTravelVideoById(v.Id)
			if err != nil {
				continue
			} else if tv != nil && tv.CoupleId == cid {
				tv.Video = video
				list = append(list, tv)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.TravelId = tid
		tv, _ = mysql.AddTravelVideo(v)
		if tv != nil {
			tv.Video = video
			list = append(list, tv)
		}
	}
	return list
}

// DelTravelVideoList
func DelTravelVideoList(uid int64, tvList []*entity.TravelVideo) []*entity.TravelVideo {
	list := make([]*entity.TravelVideo, 0)
	if tvList == nil || len(tvList) <= 0 {
		return list
	}
	for _, v := range tvList {
		if v == nil {
			continue
		}
		tv, _ := mysql.GetTravelVideoById(v.Id)
		if tv == nil || tv.UserId != uid {
			continue
		}
		err := mysql.DelTravelVideo(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddTravelAlbumList
func AddTravelAlbumList(uid, cid, tid int64, taList []*entity.TravelAlbum) []*entity.TravelAlbum {
	list := make([]*entity.TravelAlbum, 0)
	if taList == nil || len(taList) <= 0 {
		return list
	}
	// limit
	limitCount := GetLimit().TravelAlbumCount
	total := mysql.GetTravelAlbumTotalByCoupleTravel(cid, tid)
	if total+len(list) > limitCount {
		return list
	}
	for _, v := range taList {
		if v == nil || v.AlbumId <= 0 {
			continue
		}
		// 关联数据
		album, err := mysql.GetAlbumById(v.AlbumId)
		if album == nil || album.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		ta, err := mysql.GetTravelAlbumByCoupleTravelAlbum(cid, tid, v.AlbumId)
		if err != nil {
			continue
		} else if ta != nil && ta.CoupleId == cid {
			ta.Album = album
			list = append(list, ta)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			ta, err := mysql.GetTravelAlbumById(v.Id)
			if err != nil {
				continue
			} else if ta != nil && ta.CoupleId == cid {
				ta.Album = album
				list = append(list, ta)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.TravelId = tid
		ta, _ = mysql.AddTravelAlbum(v)
		if ta != nil {
			ta.Album = album
			list = append(list, ta)
		}
	}
	return list
}

// DelTravelAlbumList
func DelTravelAlbumList(uid int64, taList []*entity.TravelAlbum) []*entity.TravelAlbum {
	list := make([]*entity.TravelAlbum, 0)
	if taList == nil || len(taList) <= 0 {
		return list
	}
	for _, v := range taList {
		if v == nil {
			continue
		}
		ta, _ := mysql.GetTravelAlbumById(v.Id)
		if ta == nil || ta.UserId != uid {
			continue
		}
		err := mysql.DelTravelAlbum(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddTravelDiaryList
func AddTravelDiaryList(uid, cid, tid int64, tdList []*entity.TravelDiary) []*entity.TravelDiary {
	list := make([]*entity.TravelDiary, 0)
	if tdList == nil || len(tdList) <= 0 {
		return list
	}
	// limit
	limitCount := GetLimit().TravelDiaryCount
	total := mysql.GetTravelDiaryTotalByCoupleTravel(cid, tid)
	if total+len(list) > limitCount {
		return list
	}
	for _, v := range tdList {
		if v == nil || v.DiaryId <= 0 {
			continue
		}
		// 关联数据
		diary, err := mysql.GetDiaryById(v.DiaryId)
		if diary == nil || diary.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		td, _ := mysql.GetTravelDiaryByCoupleTravelDiary(cid, tid, v.DiaryId)
		if err != nil {
			continue
		} else if td != nil && td.CoupleId == cid {
			td.Diary = diary
			list = append(list, td)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			td, _ := mysql.GetTravelDiaryById(v.Id)
			if err != nil {
				continue
			} else if td != nil && td.CoupleId == cid {
				td.Diary = diary
				list = append(list, td)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.TravelId = tid
		td, _ = mysql.AddTravelDiary(v)
		if td != nil {
			td.Diary = diary
			list = append(list, td)
		}
	}
	return list
}

// DelTravelDiaryList
func DelTravelDiaryList(uid int64, tdList []*entity.TravelDiary) []*entity.TravelDiary {
	list := make([]*entity.TravelDiary, 0)
	if tdList == nil || len(tdList) <= 0 {
		return list
	}
	for _, v := range tdList {
		if v == nil {
			continue
		}
		td, _ := mysql.GetTravelDiaryById(v.Id)
		if td == nil || td.UserId != uid {
			continue
		}
		err := mysql.DelTravelDiary(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddTravelFoodList
func AddTravelFoodList(uid, cid, tid int64, tfList []*entity.TravelFood) []*entity.TravelFood {
	list := make([]*entity.TravelFood, 0)
	if tfList == nil || len(tfList) <= 0 {
		return list
	}
	// limit
	limitCount := GetLimit().TravelFoodCount
	total := mysql.GetTravelFoodTotalByCoupleTravel(cid, tid)
	if total+len(list) > limitCount {
		return list
	}
	for _, v := range tfList {
		if v == nil || v.FoodId <= 0 {
			continue
		}
		// 关联数据
		food, err := mysql.GetFoodById(v.FoodId)
		if food == nil || food.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		tf, _ := mysql.GetTravelFoodByCoupleTravelFood(cid, tid, v.FoodId)
		if err != nil {
			continue
		} else if tf != nil && tf.CoupleId == cid {
			tf.Food = food
			list = append(list, tf)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			tf, _ := mysql.GetTravelFoodById(v.Id)
			if err != nil {
				continue
			} else if tf != nil && tf.CoupleId == cid {
				tf.Food = food
				list = append(list, tf)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.TravelId = tid
		tf, _ = mysql.AddTravelFood(v)
		if tf != nil {
			tf.Food = food
			list = append(list, tf)
		}
	}
	return list
}

// DelTravelFoodList
func DelTravelFoodList(uid int64, tfList []*entity.TravelFood) []*entity.TravelFood {
	list := make([]*entity.TravelFood, 0)
	if tfList == nil || len(tfList) <= 0 {
		return list
	}
	for _, v := range tfList {
		if v == nil {
			continue
		}
		tf, _ := mysql.GetTravelFoodById(v.Id)
		if tf == nil || tf.UserId != uid {
			continue
		}
		err := mysql.DelTravelFood(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddTravelMovieList
func AddTravelMovieList(uid, cid, tid int64, tmList []*entity.TravelMovie) []*entity.TravelMovie {
	list := make([]*entity.TravelMovie, 0)
	if tmList == nil || len(tmList) <= 0 {
		return list
	}
	// limit
	limitCount := GetLimit().TravelMovieCount
	total := mysql.GetTravelMovieTotalByCoupleTravel(cid, tid)
	if total+len(list) > limitCount {
		return list
	}
	for _, v := range tmList {
		if v == nil || v.MovieId <= 0 {
			continue
		}
		// 关联数据
		movie, err := mysql.GetMovieById(v.MovieId)
		if movie == nil || movie.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		tm, _ := mysql.GetTravelMovieByCoupleTravelMovie(cid, tid, v.MovieId)
		if err != nil {
			continue
		} else if tm != nil && tm.CoupleId == cid {
			tm.Movie = movie
			list = append(list, tm)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			tm, _ := mysql.GetTravelMovieById(v.Id)
			if err != nil {
				continue
			} else if tm != nil && tm.CoupleId == cid {
				tm.Movie = movie
				list = append(list, tm)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.TravelId = tid
		tm, _ = mysql.AddTravelMovie(v)
		if tm != nil {
			tm.Movie = movie
			list = append(list, tm)
		}
	}
	return list
}

// DelTravelMovieList
func DelTravelMovieList(uid int64, tmList []*entity.TravelMovie) []*entity.TravelMovie {
	list := make([]*entity.TravelMovie, 0)
	if tmList == nil || len(tmList) <= 0 {
		return list
	}
	for _, v := range tmList {
		if v == nil {
			continue
		}
		tm, _ := mysql.GetTravelMovieById(v.Id)
		if tm == nil || tm.UserId != uid {
			continue
		}
		err := mysql.DelTravelMovie(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}
