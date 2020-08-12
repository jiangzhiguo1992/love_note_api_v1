package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// UpdateWallPaper 增加/修改
func UpdateWallPaper(wp *entity.WallPaper) (*entity.WallPaper, error) {
	if wp == nil {
		return wp, errors.New("nil_wall_paper")
	} else if wp.CoupleId <= 0 {
		return wp, errors.New("nil_couple")
	}
	// 旧数据获取，用model的 纯正的err
	old, err := mysql.GetWallPaperByCouple(wp.CoupleId)
	if err != nil {
		return old, err
	}
	// 权限检查
	limit := GetVipLimitByCouple(wp.CoupleId).WallPaperCount
	if old == nil || old.Id <= 0 {
		// 新的
		if len(wp.ContentImageList) <= 0 {
			return nil, errors.New("limit_content_image_nil")
		} else if len(wp.ContentImageList) > limit {
			return old, errors.New("limit_content_image_over")
		}
		wp, err = mysql.AddWallPaper(wp)
	} else {
		// 有旧的，删图片不检查
		if (len(wp.ContentImageList) > limit) && (len(wp.ContentImageList) > len(old.ContentImageList)) {
			return old, errors.New("limit_content_image_over")
		}
		old.ContentImageList = wp.ContentImageList
		wp, err = mysql.UpdateWallPaper(old)
	}
	return wp, err
}

// GetWallPaperByCouple
func GetWallPaperByCouple(cid int64) (*entity.WallPaper, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	wp, err := mysql.GetWallPaperByCouple(cid)
	if wp == nil || wp.ContentImageList == nil || len(wp.ContentImageList) <= 0 {
		return nil, errors.New("no_data_wall_paper")
	}
	return wp, err
}
