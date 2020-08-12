package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// UpdateAwardScore
func UpdateAwardScore(uid, cid int64, a *entity.Award, insert bool) (*entity.AwardScore, error) {
	if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if a == nil || a.ScoreChange == 0 {
		return nil, errors.New("nil_award")
	} else if a.HappenId <= 0 {
		return nil, errors.New("nil_user")
	}
	// 数据检查
	as, err := mysql.GetAwardScoreByUserCouple(a.HappenId, cid)
	if err != nil {
		return as, err
	} else if as == nil || as.Id <= 0 {
		as = &entity.AwardScore{
			BaseCp: entity.BaseCp{
				UserId:   a.HappenId,
				CoupleId: cid,
			},
		}
		as, err = mysql.AddAwardScore(as)
	}
	if err != nil {
		return as, err
	} else if as == nil {
		return as, errors.New("nil_award_score")
	}
	// mysql
	as.ChangeCount = as.ChangeCount + 1
	if insert {
		as.TotalScore = as.TotalScore + int64(a.ScoreChange)
	} else {
		as.TotalScore = as.TotalScore - int64(a.ScoreChange)
	}
	as, err = mysql.UpdateAwardScore(as)
	// 没有trends
	return as, err
}

// GetAwardScoreByUserCouple
func GetAwardScoreByUserCouple(uid, cid int64) (*entity.AwardScore, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// mysql
	as, err := mysql.GetAwardScoreByUserCouple(uid, cid)
	return as, err
}
