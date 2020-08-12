package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// AddSouvenir
func AddSouvenir(uid, cid int64, s *entity.Souvenir) (*entity.Souvenir, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if s == nil {
		return nil, errors.New("nil_souvenir")
	} else if len(s.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(s.Title)) > GetLimit().SouvenirTitleLength {
		return nil, errors.New("limit_title_over")
	} else if s.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else {
		//if s.Done {
		//	// 已完成的要大于生日，小于现在 1000秒的误差
		//	if s.HappenAt > (time.Now().Unix() + 1000) {
		//		return nil, errors.New("limit_happen_err")
		//	}
		//} else {
		//	// 未完成的要大于现在 1000秒的误差
		//	if s.HappenAt < (time.Now().Unix() - 1000) {
		//		return nil, errors.New("limit_happen_err")
		//	}
		//}
	}
	// 数量检查
	if mysql.GetSouvenirTotalByCouple(cid) >= int64(GetVipLimitByCouple(cid).SouvenirCount) {
		return nil, errors.New("limit_total_over")
	}
	// mysql
	s.UserId = uid
	s.CoupleId = cid
	s, err := mysql.AddSouvenir(s)
	if s == nil || err != nil {
		return nil, err
	}
	// 没有额外属性
	// 同步
	go func() {
		// trends
		pushType := entity.PUSH_TYPE_NOTE_SOUVENIR
		contentType := entity.TRENDS_CON_TYPE_SOUVENIR
		if s == nil || !s.Done {
			contentType = entity.TRENDS_CON_TYPE_WISH
			pushType = entity.PUSH_TYPE_NOTE_WISH
		}
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, contentType, s.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, s.Id, "push_title_note_update", s.Title, pushType)
	}()
	return s, err
}

// DelSouvenir
func DelSouvenir(uid, cid, sid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if sid <= 0 {
		return errors.New("nil_souvenir")
	}
	// 旧数据检查
	s, err := mysql.GetSouvenirById(sid)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("nil_souvenir")
	} else if s.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelSouvenir(s)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		contentType := entity.TRENDS_CON_TYPE_SOUVENIR
		if s == nil || !s.Done {
			contentType = entity.TRENDS_CON_TYPE_WISH
		}
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, contentType, sid)
		AddTrends(trends)
	}()
	return err
}

// UpdateSouvenir
func UpdateSouvenir(uid, cid int64, s *entity.Souvenir) (*entity.Souvenir, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if s == nil || s.Id <= 0 {
		return nil, errors.New("nil_souvenir")
	} else if len(s.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(s.Title)) > GetLimit().SouvenirTitleLength {
		return nil, errors.New("limit_title_over")
	} else if s.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else {
		//if s.Done {
		//	// 已完成的要大于生日，小于现在
		//	if s.HappenAt > time.Now().Unix() {
		//		return nil, errors.New("limit_happen_err")
		//	}
		//} else {
		//	// 未完成的要大于现在
		//	if s.HappenAt < time.Now().Unix() {
		//		return nil, errors.New("limit_happen_err")
		//	}
		//}
	}
	// 旧数据检查
	old, err := mysql.GetSouvenirById(s.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_souvenir")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// mysql
	old.HappenAt = s.HappenAt
	old.Title = s.Title
	old.Done = s.Done
	old.Longitude = s.Longitude
	old.Latitude = s.Latitude
	old.Address = s.Address
	old.CityId = s.CityId
	s, err = mysql.UpdateSouvenir(old)
	if s == nil || err != nil {
		return old, err
	}
	// 不需要额外属性
	// 同步
	go func() {
		contentType := entity.TRENDS_CON_TYPE_SOUVENIR
		if s == nil || !s.Done {
			contentType = entity.TRENDS_CON_TYPE_WISH
		}
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, contentType, s.Id)
		AddTrends(trends)
	}()
	return s, err
}

// UpdateSouvenirForeign
func UpdateSouvenirForeign(uid, cid int64, year int, s *entity.Souvenir) *entity.Souvenir {
	if uid <= 0 {
		return s
	} else if cid <= 0 {
		return s
	} else if s == nil || s.Id <= 0 {
		return s
	}
	// 旧数据检查
	old, err := mysql.GetSouvenirById(s.Id)
	if err != nil {
		return s
	} else if old == nil {
		return s
	} else if old.CoupleId != cid {
		return s
	}
	// mysql
	old.SouvenirGiftList = s.SouvenirGiftList
	old.SouvenirTravelList = s.SouvenirTravelList
	old.SouvenirAlbumList = s.SouvenirAlbumList
	old.SouvenirVideoList = s.SouvenirVideoList
	old.SouvenirFoodList = s.SouvenirFoodList
	old.SouvenirMovieList = s.SouvenirMovieList
	old.SouvenirDiaryList = s.SouvenirDiaryList
	s = updateSouvenirForeignByYear(old, year)
	if s == nil || err != nil {
		return old
	}
	// 同步
	go func() {
		contentType := entity.TRENDS_CON_TYPE_SOUVENIR
		if s == nil || !s.Done {
			contentType = entity.TRENDS_CON_TYPE_WISH
		}
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, contentType, s.Id)
		AddTrends(trends)
	}()
	return s
}

// GetSouvenirByIdWithForeign
func GetSouvenirByIdWithForeign(uid, cid, sid int64) (*entity.Souvenir, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if sid <= 0 {
		return nil, errors.New("nil_souvenir")
	}
	// mysql
	s, err := mysql.GetSouvenirById(sid)
	if err != nil {
		return nil, err
	} else if s == nil {
		return nil, errors.New("nil_souvenir")
	} else if s.CoupleId != cid {
		return nil, errors.New("db_query_refuse")
	}
	// 额外属性
	LoadSouvenirWithForeign(s)
	// 同步
	go func() {
		contentType := entity.TRENDS_CON_TYPE_SOUVENIR
		if s == nil || !s.Done {
			contentType = entity.TRENDS_CON_TYPE_WISH
		}
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, contentType, sid)
		AddTrends(trends)
	}()
	return s, err
}

// GetSouvenirListByCouple
func GetSouvenirListByCouple(uid, cid int64, done bool, page int) ([]*entity.Souvenir, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// limit
	limit := GetPageSizeLimit().Souvenir
	offset := page * limit
	if page < 0 {
		limit = GetVipLimitByCouple(cid).SouvenirCount
		offset = 0
	}
	// mysql
	var list []*entity.Souvenir
	var err error
	if done {
		list, err = mysql.GetSouvenirDoneListByCouple(cid, offset, limit)
	} else {
		list, err = mysql.GetSouvenirWishListByUserCouple(cid, offset, limit)
	}
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_souvenir")
		} else {
			return nil, nil
		}
	}
	// list没有额外属性
	// 同步
	go func() {
		contentType := entity.TRENDS_CON_TYPE_SOUVENIR
		if !done {
			contentType = entity.TRENDS_CON_TYPE_WISH
		}
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, contentType)
		AddTrends(trends)
	}()
	return list, err
}

// LoadSouvenirWithForeign
func LoadSouvenirWithForeign(s *entity.Souvenir) *entity.Souvenir {
	if s == nil || s.Id <= 0 || s.CoupleId <= 0 {
		return s
	}
	s.SouvenirGiftList, _ = mysql.GetSouvenirGiftListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirGiftList != nil && len(s.SouvenirGiftList) > 0 {
		for _, v := range s.SouvenirGiftList {
			if v == nil || v.GiftId <= 0 {
				continue
			}
			v.Gift, _ = mysql.GetGiftById(v.GiftId)
		}
	}
	s.SouvenirTravelList, _ = mysql.GetSouvenirTravelListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirTravelList != nil && len(s.SouvenirTravelList) > 0 {
		for _, v := range s.SouvenirTravelList {
			if v == nil || v.TravelId <= 0 {
				continue
			}
			v.Travel, _ = mysql.GetTravelById(v.TravelId)
			//LoadTravelWithPlace(v.Travel)
		}
	}
	s.SouvenirAlbumList, _ = mysql.GetSouvenirAlbumListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirAlbumList != nil && len(s.SouvenirAlbumList) > 0 {
		for _, v := range s.SouvenirAlbumList {
			if v == nil || v.AlbumId <= 0 {
				continue
			}
			v.Album, _ = mysql.GetAlbumById(v.AlbumId)
		}
	}
	s.SouvenirVideoList, _ = mysql.GetSouvenirVideoListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirVideoList != nil && len(s.SouvenirVideoList) > 0 {
		for _, v := range s.SouvenirVideoList {
			if v == nil || v.VideoId <= 0 {
				continue
			}
			v.Video, _ = mysql.GetVideoById(v.VideoId)
		}
	}
	s.SouvenirFoodList, _ = mysql.GetSouvenirFoodListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirFoodList != nil && len(s.SouvenirFoodList) > 0 {
		for _, v := range s.SouvenirFoodList {
			if v == nil || v.FoodId <= 0 {
				continue
			}
			v.Food, _ = mysql.GetFoodById(v.FoodId)
		}
	}
	s.SouvenirMovieList, _ = mysql.GetSouvenirMovieListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirMovieList != nil && len(s.SouvenirMovieList) > 0 {
		for _, v := range s.SouvenirMovieList {
			if v == nil || v.MovieId <= 0 {
				continue
			}
			v.Movie, _ = mysql.GetMovieById(v.MovieId)
		}
	}
	s.SouvenirDiaryList, _ = mysql.GetSouvenirDiaryListByCoupleSouvenir(s.CoupleId, s.Id)
	if s.SouvenirDiaryList != nil && len(s.SouvenirDiaryList) > 0 {
		for _, v := range s.SouvenirDiaryList {
			if v == nil || v.DiaryId <= 0 {
				continue
			}
			v.Diary, _ = mysql.GetDiaryById(v.DiaryId)
		}
	}
	return s
}

// updateSouvenirForeignByYear
func updateSouvenirForeignByYear(s *entity.Souvenir, year int) *entity.Souvenir {
	if s == nil || s.Id <= 0 || s.UserId <= 0 || s.CoupleId <= 0 || year <= 0 {
		return s
	}
	limitCount := GetLimit().SouvenirForeignYearCount
	// 礼物
	addGiftList, delGiftList := splitSouvenirGiftListByStatus(s.SouvenirGiftList)
	if delGiftList != nil && len(delGiftList) > 0 {
		DelSouvenirGiftList(s.UserId, year, delGiftList)
	}
	if addGiftList != nil && len(addGiftList) > 0 && len(addGiftList) <= limitCount {
		s.SouvenirGiftList = AddSouvenirGiftList(s.UserId, s.CoupleId, s.Id, year, addGiftList)
	} else {
		s.SouvenirGiftList = make([]*entity.SouvenirGift, 0)
	}
	// 游记
	addTravelList, delTravelList := splitSouvenirTravelListByStatus(s.SouvenirTravelList)
	if delTravelList != nil && len(delTravelList) > 0 {
		DelSouvenirTravelList(s.UserId, year, delTravelList)
	}
	if addTravelList != nil && len(addTravelList) > 0 && len(addTravelList) <= limitCount {
		s.SouvenirTravelList = AddSouvenirTravelList(s.UserId, s.CoupleId, s.Id, year, addTravelList)
	} else {
		s.SouvenirTravelList = make([]*entity.SouvenirTravel, 0)
	}
	// 相册
	addAlbumList, delAlbumList := splitSouvenirAlbumListByStatus(s.SouvenirAlbumList)
	if delAlbumList != nil && len(delAlbumList) > 0 {
		DelSouvenirAlbumList(s.UserId, year, delAlbumList)
	}
	if addAlbumList != nil && len(addAlbumList) > 0 && len(addAlbumList) <= limitCount {
		s.SouvenirAlbumList = AddSouvenirAlbumList(s.UserId, s.CoupleId, s.Id, year, addAlbumList)
	} else {
		s.SouvenirAlbumList = make([]*entity.SouvenirAlbum, 0)
	}
	// 视频
	addVideoList, delVideoList := splitSouvenirVideoListByStatus(s.SouvenirVideoList)
	if delVideoList != nil && len(delVideoList) > 0 {
		DelSouvenirVideoList(s.UserId, year, delVideoList)
	}
	if addVideoList != nil && len(addVideoList) > 0 && len(addVideoList) <= limitCount {
		s.SouvenirVideoList = AddSouvenirVideoList(s.UserId, s.CoupleId, s.Id, year, addVideoList)
	} else {
		s.SouvenirVideoList = make([]*entity.SouvenirVideo, 0)
	}
	// 美食
	addFoodList, delFoodList := splitSouvenirFoodListByStatus(s.SouvenirFoodList)
	if delFoodList != nil && len(delFoodList) > 0 {
		DelSouvenirFoodList(s.UserId, year, delFoodList)
	}
	if addFoodList != nil && len(addFoodList) > 0 && len(addFoodList) <= limitCount {
		s.SouvenirFoodList = AddSouvenirFoodList(s.UserId, s.CoupleId, s.Id, year, addFoodList)
	} else {
		s.SouvenirFoodList = make([]*entity.SouvenirFood, 0)
	}
	// 电影
	addMovieList, delMovieList := splitSouvenirMovieListByStatus(s.SouvenirMovieList)
	if delMovieList != nil && len(delMovieList) > 0 {
		DelSouvenirMovieList(s.UserId, year, delMovieList)
	}
	if addMovieList != nil && len(addMovieList) > 0 && len(addMovieList) <= limitCount {
		s.SouvenirMovieList = AddSouvenirMovieList(s.UserId, s.CoupleId, s.Id, year, addMovieList)
	} else {
		s.SouvenirMovieList = make([]*entity.SouvenirMovie, 0)
	}
	// 日记
	addDiaryList, delDiaryList := splitSouvenirDiaryListByStatus(s.SouvenirDiaryList)
	if delDiaryList != nil && len(delDiaryList) > 0 {
		DelSouvenirDiaryList(s.UserId, year, delDiaryList)
	}
	if addDiaryList != nil && len(addDiaryList) > 0 && len(addDiaryList) <= limitCount {
		s.SouvenirDiaryList = AddSouvenirDiaryList(s.UserId, s.CoupleId, s.Id, year, addDiaryList)
	} else {
		s.SouvenirDiaryList = make([]*entity.SouvenirDiary, 0)
	}
	return s
}

// splitSouvenirGiftListByStatus
func splitSouvenirGiftListByStatus(objList []*entity.SouvenirGift) ([]*entity.SouvenirGift, []*entity.SouvenirGift) {
	addList := make([]*entity.SouvenirGift, 0)
	delList := make([]*entity.SouvenirGift, 0)
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

// splitSouvenirTravelListByStatus
func splitSouvenirTravelListByStatus(objList []*entity.SouvenirTravel) ([]*entity.SouvenirTravel, []*entity.SouvenirTravel) {
	addList := make([]*entity.SouvenirTravel, 0)
	delList := make([]*entity.SouvenirTravel, 0)
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

// splitSouvenirAlbumListByStatus
func splitSouvenirAlbumListByStatus(objList []*entity.SouvenirAlbum) ([]*entity.SouvenirAlbum, []*entity.SouvenirAlbum) {
	addList := make([]*entity.SouvenirAlbum, 0)
	delList := make([]*entity.SouvenirAlbum, 0)
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

// splitSouvenirVideoListByStatus
func splitSouvenirVideoListByStatus(objList []*entity.SouvenirVideo) ([]*entity.SouvenirVideo, []*entity.SouvenirVideo) {
	addList := make([]*entity.SouvenirVideo, 0)
	delList := make([]*entity.SouvenirVideo, 0)
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

// splitSouvenirFoodListByStatus
func splitSouvenirFoodListByStatus(objList []*entity.SouvenirFood) ([]*entity.SouvenirFood, []*entity.SouvenirFood) {
	addList := make([]*entity.SouvenirFood, 0)
	delList := make([]*entity.SouvenirFood, 0)
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

// splitSouvenirMovieListByStatus
func splitSouvenirMovieListByStatus(objList []*entity.SouvenirMovie) ([]*entity.SouvenirMovie, []*entity.SouvenirMovie) {
	addList := make([]*entity.SouvenirMovie, 0)
	delList := make([]*entity.SouvenirMovie, 0)
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

// splitSouvenirDiaryListByStatus
func splitSouvenirDiaryListByStatus(objList []*entity.SouvenirDiary) ([]*entity.SouvenirDiary, []*entity.SouvenirDiary) {
	addList := make([]*entity.SouvenirDiary, 0)
	delList := make([]*entity.SouvenirDiary, 0)
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

// AddSouvenirAlbumList
func AddSouvenirAlbumList(uid, cid, sid int64, year int, saList []*entity.SouvenirAlbum) []*entity.SouvenirAlbum {
	list := make([]*entity.SouvenirAlbum, 0)
	if saList == nil || len(saList) <= 0 {
		return list
	}
	for _, v := range saList {
		if v == nil || v.AlbumId <= 0 {
			continue
		}
		// 关联数据
		album, err := mysql.GetAlbumById(v.AlbumId)
		if album == nil || album.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		sa, err := mysql.GetSouvenirAlbumByCoupleSouvenirAlbum(cid, sid, v.AlbumId)
		if err != nil {
			continue
		} else if sa != nil && sa.CoupleId == cid && sa.Year == year {
			sa.Album = album
			list = append(list, sa)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			sa, err := mysql.GetSouvenirAlbumById(v.Id)
			if err != nil {
				continue
			} else if sa != nil && sa.CoupleId == cid && sa.Year == year {
				sa.Album = album
				list = append(list, sa)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		sa, _ = mysql.AddSouvenirAlbum(v)
		if sa != nil {
			sa.Album = album
			list = append(list, sa)
		}
	}
	return list
}

// DelSouvenirAlbumList
func DelSouvenirAlbumList(uid int64, year int, saList []*entity.SouvenirAlbum) []*entity.SouvenirAlbum {
	list := make([]*entity.SouvenirAlbum, 0)
	if saList == nil || len(saList) <= 0 {
		return list
	}
	for _, v := range saList {
		if v == nil {
			continue
		}
		sa, _ := mysql.GetSouvenirAlbumById(v.Id)
		if sa == nil || sa.UserId != uid || sa.Year != year {
			continue
		}
		err := mysql.DelSouvenirAlbum(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddSouvenirDiaryList
func AddSouvenirDiaryList(uid, cid, sid int64, year int, sdList []*entity.SouvenirDiary) []*entity.SouvenirDiary {
	list := make([]*entity.SouvenirDiary, 0)
	if sdList == nil || len(sdList) <= 0 {
		return list
	}
	for _, v := range sdList {
		if v == nil || v.DiaryId <= 0 {
			continue
		}
		// 关联数据
		diary, err := mysql.GetDiaryById(v.DiaryId)
		if diary == nil || diary.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		sd, err := mysql.GetSouvenirDiaryByCoupleSouvenirDiary(cid, sid, v.DiaryId)
		if err != nil {
			continue
		} else if sd != nil && sd.CoupleId == cid && sd.Year == year {
			sd.Diary = diary
			list = append(list, sd)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			sd, err := mysql.GetSouvenirDiaryById(v.Id)
			if err != nil {
				continue
			} else if sd != nil && sd.CoupleId == cid && sd.Year == year {
				sd.Diary = diary
				list = append(list, sd)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		sd, _ = mysql.AddSouvenirDiary(v)
		if sd != nil {
			sd.Diary = diary
			list = append(list, sd)
		}
	}
	return list
}

// DelSouvenirDiaryList
func DelSouvenirDiaryList(uid int64, year int, sdList []*entity.SouvenirDiary) []*entity.SouvenirDiary {
	list := make([]*entity.SouvenirDiary, 0)
	if sdList == nil || len(sdList) <= 0 {
		return list
	}
	for _, v := range sdList {
		if v == nil {
			continue
		}
		sd, _ := mysql.GetSouvenirDiaryById(v.Id)
		if sd == nil || sd.UserId != uid || sd.Year != year {
			continue
		}
		err := mysql.DelSouvenirDiary(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddSouvenirFoodList
func AddSouvenirFoodList(uid, cid, sid int64, year int, sfList []*entity.SouvenirFood) []*entity.SouvenirFood {
	list := make([]*entity.SouvenirFood, 0)
	if sfList == nil || len(sfList) <= 0 {
		return list
	}
	for _, v := range sfList {
		if v == nil || v.FoodId <= 0 {
			continue
		}
		// 关联数据
		food, err := mysql.GetFoodById(v.FoodId)
		if food == nil || food.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		sf, err := mysql.GetSouvenirFoodByCoupleSouvenirFood(cid, sid, v.FoodId)
		if err != nil {
			continue
		} else if sf != nil && sf.CoupleId == cid && sf.Year == year {
			sf.Food = food
			list = append(list, sf)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			sf, err := mysql.GetSouvenirFoodById(v.Id)
			if err != nil {
				continue
			} else if sf != nil && sf.CoupleId == cid && sf.Year == year {
				sf.Food = food
				list = append(list, sf)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		sf, _ = mysql.AddSouvenirFood(v)
		if sf != nil {
			sf.Food = food
			list = append(list, sf)
		}
	}
	return list
}

// DelSouvenirFoodList
func DelSouvenirFoodList(uid int64, year int, sfList []*entity.SouvenirFood) []*entity.SouvenirFood {
	list := make([]*entity.SouvenirFood, 0)
	if sfList == nil || len(sfList) <= 0 {
		return list
	}
	for _, v := range sfList {
		if v == nil {
			continue
		}
		sf, _ := mysql.GetSouvenirFoodById(v.Id)
		if sf == nil || sf.UserId != uid || sf.Year != year {
			continue
		}
		err := mysql.DelSouvenirFood(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddSouvenirMovieList
func AddSouvenirMovieList(uid, cid, sid int64, year int, smList []*entity.SouvenirMovie) []*entity.SouvenirMovie {
	list := make([]*entity.SouvenirMovie, 0)
	if smList == nil || len(smList) <= 0 {
		return list
	}
	for _, v := range smList {
		if v == nil || v.MovieId <= 0 {
			continue
		}
		// 关联数据
		movie, err := mysql.GetMovieById(v.MovieId)
		if movie == nil || movie.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		sm, err := mysql.GetSouvenirMovieByCoupleSouvenirMovie(cid, sid, v.MovieId)
		if err != nil {
			continue
		} else if sm != nil && sm.CoupleId == cid && sm.Year == year {
			sm.Movie = movie
			list = append(list, sm)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			sm, err := mysql.GetSouvenirMovieById(v.Id)
			if err != nil {
				continue
			} else if sm != nil && sm.CoupleId == cid && sm.Year == year {
				sm.Movie = movie
				list = append(list, sm)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		sm, _ = mysql.AddSouvenirMovie(v)
		if sm != nil {
			sm.Movie = movie
			list = append(list, sm)
		}
	}
	return list
}

// DelSouvenirMovieList
func DelSouvenirMovieList(uid int64, year int, smList []*entity.SouvenirMovie) []*entity.SouvenirMovie {
	list := make([]*entity.SouvenirMovie, 0)
	if smList == nil || len(smList) <= 0 {
		return list
	}
	for _, v := range smList {
		if v == nil {
			continue
		}
		sm, _ := mysql.GetSouvenirMovieById(v.Id)
		if sm == nil || sm.UserId != uid || sm.Year != year {
			continue
		}
		err := mysql.DelSouvenirMovie(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddSouvenirGiftList
func AddSouvenirGiftList(uid, cid, sid int64, year int, sgList []*entity.SouvenirGift) []*entity.SouvenirGift {
	list := make([]*entity.SouvenirGift, 0)
	if sgList == nil || len(sgList) <= 0 {
		return list
	}
	for _, v := range sgList {
		if v == nil || v.GiftId <= 0 {
			continue
		}
		// 关联数据
		gift, err := mysql.GetGiftById(v.GiftId)
		if gift == nil || gift.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		sg, err := mysql.GetSouvenirGiftByCoupleSouvenirGift(cid, sid, v.GiftId)
		if err != nil {
			continue
		} else if sg != nil && sg.CoupleId == cid && sg.Year == year {
			sg.Gift = gift
			list = append(list, sg)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			sg, err := mysql.GetSouvenirGiftById(v.Id)
			if err != nil {
				continue
			} else if sg != nil && sg.CoupleId == cid && sg.Year == year {
				sg.Gift = gift
				list = append(list, sg)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		sg, _ = mysql.AddSouvenirGift(v)
		if sg != nil {
			sg.Gift = gift
			list = append(list, sg)
		}
	}
	return list
}

// DelSouvenirGiftList
func DelSouvenirGiftList(uid int64, year int, sgList []*entity.SouvenirGift) []*entity.SouvenirGift {
	list := make([]*entity.SouvenirGift, 0)
	if sgList == nil || len(sgList) <= 0 {
		return list
	}
	for _, v := range sgList {
		if v == nil {
			continue
		}
		sg, _ := mysql.GetSouvenirGiftById(v.Id)
		if sg == nil || sg.UserId != uid || sg.Year != year {
			continue
		}
		err := mysql.DelSouvenirGift(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddSouvenirTravelList
func AddSouvenirTravelList(uid, cid, sid int64, year int, stList []*entity.SouvenirTravel) []*entity.SouvenirTravel {
	list := make([]*entity.SouvenirTravel, 0)
	if stList == nil || len(stList) <= 0 {
		return list
	}
	for _, v := range stList {
		if v == nil || v.TravelId <= 0 {
			continue
		}
		// 关联数据
		travel, err := mysql.GetTravelById(v.TravelId)
		if travel == nil || travel.CoupleId != cid || err != nil {
			continue
		}
		//LoadTravelWithPlace(travel)
		// 旧数据
		st, err := mysql.GetSouvenirTravelByCoupleSouvenirTravel(cid, sid, v.TravelId)
		if err != nil {
			continue
		} else if st != nil && st.CoupleId == cid && st.Year == year {
			st.Travel = travel
			list = append(list, st)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			st, err := mysql.GetSouvenirTravelById(v.Id)
			if err != nil {
				continue
			} else if st != nil && st.CoupleId == cid && st.Year == year {
				st.Travel = travel
				list = append(list, st)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		st, _ = mysql.AddSouvenirTravel(v)
		if st != nil {
			st.Travel = travel
			list = append(list, st)
		}
	}
	return list
}

// DelSouvenirTravelList
func DelSouvenirTravelList(uid int64, year int, stList []*entity.SouvenirTravel) []*entity.SouvenirTravel {
	list := make([]*entity.SouvenirTravel, 0)
	if stList == nil || len(stList) <= 0 {
		return list
	}
	for _, v := range stList {
		if v == nil {
			continue
		}
		st, _ := mysql.GetSouvenirTravelById(v.Id)
		if st == nil || st.UserId != uid || st.Year != year {
			continue
		}
		err := mysql.DelSouvenirTravel(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// AddSouvenirVideoList
func AddSouvenirVideoList(uid, cid, sid int64, year int, svList []*entity.SouvenirVideo) []*entity.SouvenirVideo {
	list := make([]*entity.SouvenirVideo, 0)
	if svList == nil || len(svList) <= 0 {
		return list
	}
	for _, v := range svList {
		if v == nil || v.VideoId <= 0 {
			continue
		}
		// 关联数据
		video, err := mysql.GetVideoById(v.VideoId)
		if video == nil || video.CoupleId != cid || err != nil {
			continue
		}
		// 旧数据
		sv, err := mysql.GetSouvenirVideoByCoupleSouvenirVideo(cid, sid, v.VideoId)
		if err != nil {
			continue
		} else if sv != nil && sv.CoupleId == cid && sv.Year == year {
			sv.Video = video
			list = append(list, sv)
			continue
		}
		// 依然旧数据
		if v.Id > 0 {
			sv, err := mysql.GetSouvenirVideoById(v.Id)
			if err != nil {
				continue
			} else if sv != nil && sv.CoupleId == cid && sv.Year == year {
				sv.Video = video
				list = append(list, sv)
				continue
			}
		}
		// 加新的
		v.UserId = uid
		v.CoupleId = cid
		v.SouvenirId = sid
		v.Year = year
		sv, _ = mysql.AddSouvenirVideo(v)
		if sv != nil {
			sv.Video = video
			list = append(list, sv)
		}
	}
	return list
}

// DelSouvenirVideoList
func DelSouvenirVideoList(uid int64, year int, svList []*entity.SouvenirVideo) []*entity.SouvenirVideo {
	list := make([]*entity.SouvenirVideo, 0)
	if svList == nil || len(svList) <= 0 {
		return list
	}
	for _, v := range svList {
		if v == nil {
			continue
		}
		sv, _ := mysql.GetSouvenirVideoById(v.Id)
		if sv == nil || sv.UserId != uid || sv.Year != year {
			continue
		}
		err := mysql.DelSouvenirVideo(v)
		if err != nil {
			list = append(list, v)
		}
	}
	return list
}

// GetSouvenirLatestByList 获取最近的souvenir
func GetSouvenirLatestByList(souvenirList []*entity.Souvenir, nearAt int64) *entity.Souvenir {
	if souvenirList == nil || len(souvenirList) <= 0 {
		return nil
	}
	var latest *entity.Souvenir
	if len(souvenirList) > 0 {
		nSouvenir := &entity.Souvenir{HappenAt: nearAt}
		latest = souvenirList[0]
		for i := 1; i < len(souvenirList); i++ {
			iSouvenir := souvenirList[i]
			// 开始比较 lSouvenir 和 iSouvenir 中最先到达的
			ilList := sortSouvenirUpInMonth(iSouvenir, latest)
			if ilList[0] == iSouvenir {
				// iMonth < lMonth
				inList := sortSouvenirUpInMonth(iSouvenir, nSouvenir)
				if inList[1] == iSouvenir {
					// nMonth < iMonth < lMonth
					latest = iSouvenir
					continue
				} else {
					// iMonth < lMonth/nMonth
					lnList := sortSouvenirUpInMonth(latest, nSouvenir)
					if lnList[1] == latest {
						// iMonth < nMonth < lMonth
						continue
					} else {
						// iMonth < lMonth < nMonth
						latest = iSouvenir
						continue
					}
				}
			} else {
				// lMonth < iMonth
				lnList := sortSouvenirUpInMonth(latest, nSouvenir)
				if lnList[1] == latest {
					// nMonth < lMonth < iMonth
					continue
				} else {
					// lMonth < iMonth/nMonth
					inList := sortSouvenirUpInMonth(iSouvenir, nSouvenir)
					if inList[1] == iSouvenir {
						// lMonth < nMonth < iMonth
						latest = iSouvenir
						continue
					} else {
						// lMonth < iMonth < nMonth
						continue
					}
				}
			}
		}
	}
	return latest
}

// sortSouvenirUpInMonth 比较两个纪念日的年份以下的时间，升序返回
func sortSouvenirUpInMonth(s1, s2 *entity.Souvenir) []*entity.Souvenir {
	list := make([]*entity.Souvenir, 2)
	s1Happen := utils.GetCSTDateByUnix(s1.HappenAt)
	s2Happen := utils.GetCSTDateByUnix(s2.HappenAt)
	_, s1Month, _ := s1Happen.Date()
	_, s2Month, _ := s2Happen.Date()
	if s1Month < s2Month {
		list[0] = s1
		list[1] = s2
	} else if s2Month < s1Month {
		list[0] = s2
		list[1] = s1
	} else {
		return sortSouvenirUpInDay(s1, s2)
	}
	return list
}

// sortSouvenirUpInDay
func sortSouvenirUpInDay(s1, s2 *entity.Souvenir) []*entity.Souvenir {
	list := make([]*entity.Souvenir, 2)
	s1Happen := utils.GetCSTDateByUnix(s1.HappenAt)
	s2Happen := utils.GetCSTDateByUnix(s2.HappenAt)
	_, _, s1Day := s1Happen.Date()
	_, _, s2Day := s2Happen.Date()
	if s1Day < s2Day {
		list[0] = s1
		list[1] = s2
	} else if s2Day < s1Day {
		list[0] = s2
		list[1] = s1
	} else {
		return sortSouvenirUpInHour(s1, s2)
	}
	return list
}

// sortSouvenirUpInHour
func sortSouvenirUpInHour(s1, s2 *entity.Souvenir) []*entity.Souvenir {
	list := make([]*entity.Souvenir, 2)
	s1Happen := utils.GetCSTDateByUnix(s1.HappenAt)
	s2Happen := utils.GetCSTDateByUnix(s2.HappenAt)
	s1Hour, _, _ := s1Happen.Clock()
	s2Hour, _, _ := s2Happen.Clock()
	if s1Hour < s2Hour {
		list[0] = s1
		list[1] = s2
	} else if s2Hour < s1Hour {
		list[0] = s2
		list[1] = s1
	} else {
		return sortSouvenirUpInMin(s1, s2)
	}
	return list
}

// sortSouvenirUpInMin
func sortSouvenirUpInMin(s1, s2 *entity.Souvenir) []*entity.Souvenir {
	list := make([]*entity.Souvenir, 2)
	s1Happen := utils.GetCSTDateByUnix(s1.HappenAt)
	s2Happen := utils.GetCSTDateByUnix(s2.HappenAt)
	_, s1Min, _ := s1Happen.Clock()
	_, s2Min, _ := s2Happen.Clock()
	if s1Min < s2Min {
		list[0] = s1
		list[1] = s2
	} else if s2Min < s1Min {
		list[0] = s2
		list[1] = s1
	} else {
		return sortSouvenirUpInSec(s1, s2)
	}
	return list
}

// sortSouvenirUpInSec
func sortSouvenirUpInSec(s1, s2 *entity.Souvenir) []*entity.Souvenir {
	list := make([]*entity.Souvenir, 2)
	s1Happen := utils.GetCSTDateByUnix(s1.HappenAt)
	s2Happen := utils.GetCSTDateByUnix(s2.HappenAt)
	_, _, s1Sec := s1Happen.Clock()
	_, _, s2Sec := s2Happen.Clock()
	if s1Sec <= s2Sec {
		list[0] = s1
		list[1] = s2
	} else {
		list[0] = s2
		list[1] = s1
	}
	return list
}
