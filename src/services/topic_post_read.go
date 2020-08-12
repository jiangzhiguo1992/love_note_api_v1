package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddPostRead
func AddPostRead(uid, pid int64) (*entity.PostRead, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if pid <= 0 {
		return nil, errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pid)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// old
	old, err := mysql.GetPostReadByUser(uid, pid)
	if err != nil {
		return nil, err
	} else if old == nil || old.Id <= 0 {
		pr := &entity.PostRead{
			UserId: uid,
			PostId: pid,
		}
		old, err = mysql.AddPostRead(pr)
	} else {
		old, err = mysql.UpdatePostRead(old)
	}
	if old == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		TopicInfoUpBrowse(p.Kind)
	}()
	return old, err
}

// IsPostReadByUserCouple
func IsPostReadByUserCouple(uid, pid int64) bool {
	if uid <= 0 || pid <= 0 {
		return false
	}
	read, _ := mysql.GetPostReadByUser(uid, pid)
	if read == nil || read.Id <= 0 {
		return false
	}
	return true
}
