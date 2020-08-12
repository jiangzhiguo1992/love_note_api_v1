package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
)

// AddPostCommentReport
func AddPostCommentReport(uid, cid int64, pcr *entity.PostCommentReport) (*entity.PostCommentReport, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if pcr == nil || pcr.PostCommentId <= 0 {
		return nil, errors.New("nil_comment")
	}
	// comment检查
	pc, err := GetPostCommentById(pcr.PostCommentId)
	if err != nil {
		return nil, err
	} else if pc == nil {
		return nil, errors.New("nil_comment")
	} else if pc.PostId <= 0 {
		return nil, errors.New("nil_post")
	} else if pc.Official {
		return nil, errors.New("report_refuse")
	}
	// post检查
	p, err := GetPostById(pc.PostId)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, errors.New("nil_post")
	}
	// old
	old, err := mysql.GetPostCommentReportByUserCouple(uid, cid, pc.Id)
	if err != nil {
		return nil, err
	} else if old != nil {
		return nil, errors.New("report_repeat")
	}
	// mysql
	pcr.UserId = uid
	pcr.CoupleId = cid
	pcr, err = mysql.AddPostCommentReport(pcr)
	if pcr == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// post
		pc.ReportCount = pc.ReportCount + 1
		UpdatePostCommentCount(pc, false)
		// message
		if uid != pc.UserId && p.Kind != entity.POST_KIND_LIMIT_UNKNOWN {
			language := "zh-cn"
			entry, err := mysql.GetEntryLatestByUser(pc.UserId)
			if err == nil && entry != nil {
				language = entry.Language
			}
			content := utils.GetLanguage(language, "topic_message_comment_report") + pc.ContentText
			var conId int64
			if pc.ToCommentId > 0 {
				conId = pc.ToCommentId
			} else {
				conId = pc.Id
			}
			message := CreateTopicMessage(uid, cid, pc.UserId, pc.CoupleId, entity.TOPIC_MESSAGE_KIND_COMMENT_BE_REPORT, content, conId)
			AddTopicMessage(language, message)
		}
		// topicInfo
		TopicInfoUpReport(p.Kind)
	}()
	return pcr, err
}

// IsPostCommentReportByUserCouple
func IsPostCommentReportByUserCouple(uid, cid, pcid int64) bool {
	if uid <= 0 || cid <= 0 || pcid <= 0 {
		return false
	}
	report, _ := mysql.GetPostCommentReportByUserCouple(uid, cid, pcid)
	if report == nil || report.Id <= 0 {
		return false
	}
	return true
}
