package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddSuggestComment
func AddSuggestComment(uid int64, sc *entity.SuggestComment) (*entity.SuggestComment, error) {
	if sc == nil {
		return nil, errors.New("nil_comment")
	} else if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if sc.SuggestId <= 0 {
		return nil, errors.New("nil_suggest")
	} else if len(strings.TrimSpace(sc.ContentText)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(sc.ContentText)) > GetLimit().SuggestCommentContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// 检查数据
	s, err := mysql.GetSuggestById(sc.SuggestId)
	if err != nil {
		return nil, err
	} else if s == nil {
		return nil, errors.New("nil_suggest")
	}
	// admin
	u, _ := GetUserById(uid)
	sc.Official = IsAdminister(u)
	// mysql
	sc.UserId = uid
	sc, err = mysql.AddSuggestComment(sc)
	if sc == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// 评论数
		s.CommentCount = s.CommentCount + 1
		mysql.UpdateSuggestCount(s, true)
		// language
		language := "zh-cn"
		entry, err := mysql.GetEntryLatestByUser(s.UserId)
		if err == nil && entry != nil {
			language = entry.Language
		}
		// push
		title := utils.GetLanguage(language, "push_title_new_comment")
		push := CreatePush(sc.UserId, s.UserId, s.Id, title, sc.ContentText, entity.PUSH_TYPE_SUGGEST)
		AddPush(push)
	}()
	return sc, err
}

// DelSuggestComment
func DelSuggestComment(uid, scid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if scid <= 0 {
		return errors.New("nil_comment")
	}
	// 旧数据检查
	sc, err := mysql.GetSuggestCommentById(scid)
	if err != nil {
		return err
	} else if sc == nil {
		return errors.New("nil_comment")
	}
	s, err := mysql.GetSuggestById(sc.SuggestId)
	if err != nil {
		return err
	} else if s == nil {
		return errors.New("nil_suggest")
	}
	// admin
	u, _ := GetUserById(uid)
	if !IsAdminister(u) && sc.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelSuggestComment(sc)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		s.CommentCount = s.CommentCount - 1
		mysql.UpdateSuggestCount(s, false)
	}()
	return err
}

// GetSuggestCommentList
func GetSuggestCommentList(uid, sid int64, page int) ([]*entity.SuggestComment, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().SuggestComment
	offset := page * limit
	list, err := mysql.GetSuggestCommentList(uid, sid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	return list, nil
}

// GetSuggestCommentListWithAll
func GetSuggestCommentListWithAll(uid, sid int64, page int) ([]*entity.SuggestComment, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if sid <= 0 {
		return nil, errors.New("nil_suggest")
	}
	list, err := GetSuggestCommentList(0, sid, page)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_comment")
		} else {
			return nil, nil
		}
	}
	// 额外属性
	for _, sc := range list {
		sc.Mine = (uid == sc.UserId) && (uid > 0)
	}
	return list, nil
}

// GetSuggestCommentTotalByCreate
func GetSuggestCommentTotalByCreate(create int64) int64 {
	if create == 0 {
		return 0
	}
	// mysql
	total := mysql.GetSuggestCommentTotalByCreateWithDel(create)
	return total
}

// IsSuggestCommentByUser
func IsSuggestCommentByUser(uid, sid int64) bool {
	if uid <= 0 {
		return false
	} else if sid <= 0 {
		return false
	}
	comment := mysql.GetSuggestCommentTotalByUser(uid, sid) > 0
	return comment
}
