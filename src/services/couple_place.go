package services

import (
	"errors"
	"math"
	"models/entity"
	"models/mysql"
)

// AddPlace 添加地区信息
func AddPlace(uid, cid int64, p *entity.Place) (*entity.Place, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if p == nil {
		return nil, errors.New("nil_place")
	} else if p.Longitude == 0 && p.Latitude == 0 {
		return nil, errors.New("limit_lat_lon_nil")
	}
	// old
	old, _ := GetPlaceLatestByUserCouple(uid, cid)
	if old != nil {
		// 距离小于1000米，就不更新
		if earthDistance(old.Longitude, old.Latitude, p.Longitude, p.Latitude) < 1000 {
			// 记住返回原来的数据
			return old, nil
		}
	}
	// mysql
	p.UserId = uid
	p.CoupleId = cid
	p, err := mysql.AddPlace(p)
	if p == nil || err != nil {
		return nil, err
	}
	return p, err
}

// GetPlaceLatestByUserCouple 获取最新的地区记录
func GetPlaceLatestByUserCouple(uid, cid int64) (*entity.Place, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	p, err := mysql.GetPlaceLatestByUserCouple(uid, cid)
	return p, err
}

// GetPlaceListByCouple
func GetPlaceListByCouple(cid int64, page int) ([]*entity.Place, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	limit := GetPageSizeLimit().Place
	offset := page * limit
	list, err := mysql.GetPlaceListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_place")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetPlaceList
func GetPlaceList(uid int64, page int) ([]*entity.Place, error) {
	// mysql
	limit := GetPageSizeLimit().Place
	offset := page * limit
	list, err := mysql.GetPlaceList(uid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_place")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetPlaceGroupCountryList 国家
func GetPlaceGroupCountryList(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetPlaceFilerListByCreate("country", start, end)
}

// GetPlaceGroupProvinceList 省份
func GetPlaceGroupProvinceList(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetPlaceFilerListByCreate("province", start, end)
}

// GetPlaceGroupCityList 城市
func GetPlaceGroupCityList(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetPlaceFilerListByCreate("city", start, end)
}

// GetPlaceGroupDistrictList 辖区
func GetPlaceGroupDistrictList(start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetPlaceFilerListByCreate("district", start, end)
}

// 计算经纬度距离
func earthDistance(lng1, lat1, lng2, lat2 float64) float64 {
	radius := 6371000 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * float64(radius)
}
