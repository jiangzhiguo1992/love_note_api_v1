package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddFood
func AddFood(uid, cid int64, f *entity.Food) (*entity.Food, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if f == nil {
		return nil, errors.New("nil_food")
	} else if f.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(f.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(f.Title)) > GetLimit().FoodTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len([]rune(f.ContentText)) > GetLimit().FoodContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// limit
	if len(f.ContentImageList) > 0 {
		imgLimit := GetVipLimitByCouple(cid).FoodImageCount
		if imgLimit <= 0 {
			return nil, errors.New("limit_content_image_refuse")
		} else if len(f.ContentImageList) > imgLimit {
			return nil, errors.New("limit_content_image_over")
		}
	}
	// mysql
	f.UserId = uid
	f.CoupleId = cid
	f, err := mysql.AddFood(f)
	if f == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_FOOD, f.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, f.Id, "push_title_note_update", f.Title, entity.PUSH_TYPE_NOTE_FOOD)
	}()
	return f, err
}

// DelFood
func DelFood(uid, cid, fid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if fid <= 0 {
		return errors.New("nil_food")
	}
	// 旧数据检查
	f, err := mysql.GetFoodById(fid)
	if err != nil {
		return err
	} else if f == nil {
		return errors.New("nil_food")
	} else if f.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelFood(f)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_FOOD, fid)
		AddTrends(trends)
	}()
	return err
}

// UpdateFood
func UpdateFood(uid, cid int64, f *entity.Food) (*entity.Food, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if f == nil || f.Id <= 0 {
		return nil, errors.New("nil_food")
	} else if f.HappenAt == 0 {
		return nil, errors.New("limit_happen_nil")
	} else if len(f.Title) <= 0 {
		return nil, errors.New("limit_title_nil")
	} else if len([]rune(f.Title)) > GetLimit().FoodTitleLength {
		return nil, errors.New("limit_title_over")
	} else if len([]rune(f.ContentText)) > GetLimit().FoodContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 旧数据检查
	old, err := mysql.GetFoodById(f.Id)
	if err != nil {
		return old, err
	} else if old == nil {
		return old, errors.New("nil_food")
	} else if old.UserId != uid {
		return old, errors.New("db_update_refuse")
	}
	// 图片检查
	limit := GetVipLimitByCouple(cid).FoodImageCount
	if (len(f.ContentImageList) > limit) && (len(f.ContentImageList) > len(old.ContentImageList)) {
		// 修改的图数大于限制图数，如果是以前vip传上去的，则通过
		return old, errors.New("limit_content_image_over")
	}
	// mysql
	old.HappenAt = f.HappenAt
	old.Title = f.Title
	old.ContentImageList = f.ContentImageList
	old.ContentText = f.ContentText
	old.Longitude = f.Longitude
	old.Latitude = f.Latitude
	old.Address = f.Address
	old.CityId = f.CityId
	f, err = mysql.UpdateFood(old)
	if f == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_UPDATE, entity.TRENDS_CON_TYPE_FOOD, f.Id)
		AddTrends(trends)
	}()
	return f, err
}

// GetFoodListByCouple
func GetFoodListByCouple(uid, cid int64, page int) ([]*entity.Food, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Food
	offset := page * limit
	list, err := mysql.GetFoodListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_food")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_FOOD)
		AddTrends(trends)
	}()
	return list, err
}
