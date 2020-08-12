package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// TogglePostCommentPoint
func TogglePostCommentPoint(uid, cid int64, pcp *entity.PostCommentPoint) (*entity.PostCommentPoint, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pcp == nil || pcp.PostCommentId <= 0 {
		return nil, errors.New("nil_comment")
	}
	// comment检查
	pc, err := GetPostCommentById(pcp.PostCommentId)
	if err != nil {
		return nil, err
	} else if pc == nil {
		return nil, errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		return nil, errors.New("nil_post")
	}
	// post检查
	p, err := GetPostById(pc.PostId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// mysql
	old, err := mysql.GetPostCommentPointByUserCouple(uid, cid, pcp.PostCommentId)
	if err != nil {
		return old, err
	} else if old == nil || old.Id <= 0 {
		// 没点赞
		pcp.UserId = uid
		pcp.CoupleId = cid
		pcp, err = mysql.AddPostCommentPoint(pcp)
	} else {
		// 已点赞
		if old.Status >= entity.STATUS_VISIBLE {
			old.Status = entity.STATUS_DELETE
		} else {
			old.Status = entity.STATUS_VISIBLE
		}
		pcp, err = mysql.UpdatePostCommentPoint(old)
	}
	if pcp == nil || err != nil {
		return old, err
	}
	// 同步
	go func() {
		// post
		if pcp.Status >= entity.STATUS_VISIBLE {
			pc.PointCount = pc.PointCount + 1
		} else {
			pc.PointCount = pc.PointCount - 1
		}
		UpdatePostCommentCount(pc, false)
		// message
		if uid != pc.UserId &&
			pcp.Status >= entity.STATUS_VISIBLE &&
			p.Kind != entity.POST_KIND_LIMIT_UNKNOWN {
			language := "zh-cn"
			entry, err := mysql.GetEntryLatestByUser(pc.UserId)
			if err == nil && entry != nil {
				language = entry.Language
			}
			content := utils.GetLanguage(language, "topic_message_comment_point") + pc.ContentText
			var conId int64
			if pc.ToCommentId > 0 {
				conId = pc.ToCommentId
			} else {
				conId = pc.Id
			}
			message := CreateTopicMessage(uid, cid, pc.UserId, pc.CoupleId, entity.TOPIC_MESSAGE_KIND_COMMENT_BE_POINT, content, conId)
			AddTopicMessage(language, message)
		}
		// topicInfo
		TopicInfoUpdatePoint(p.Kind, pcp.Status >= entity.STATUS_VISIBLE)
	}()
	return pcp, err
}

// IsPostCommentPointByUserCouple
func IsPostCommentPointByUserCouple(uid, cid, pcid int64) bool {
	if uid <= 0 || cid <= 0 || pcid <= 0 {
		return false
	}
	point, _ := mysql.GetPostCommentPointByUserCouple(uid, cid, pcid)
	if point == nil || point.Id <= 0 {
		return false
	} else if point.Status < entity.STATUS_VISIBLE {
		return false
	}
	return true
}
