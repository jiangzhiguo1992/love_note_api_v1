package services

import (
	"errors"
	"models/entity"
	"models/mysql"
)

// AddWord
func AddWord(uid, cid int64, w *entity.Word) (*entity.Word, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if w == nil {
		return nil, errors.New("nil_word")
	} else if len(w.ContentText) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len([]rune(w.ContentText)) > GetLimit().WordContentLength {
		return nil, errors.New("limit_content_text_over")
	}
	// mysql
	w.UserId = uid
	w.CoupleId = cid
	w, err := mysql.AddWord(w)
	if w == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_WORD, w.Id)
		AddTrends(trends)
		// push
		AddPushInCouple(uid, w.Id, "push_title_note_update", w.ContentText, entity.PUSH_TYPE_NOTE_WORD)
	}()
	return w, err
}

// DelWord
func DelWord(uid, cid, wid int64) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	} else if wid <= 0 {
		return errors.New("nil_word")
	}
	// 旧数据检查
	w, err := mysql.GetWordById(wid)
	if err != nil {
		return err
	} else if w == nil {
		return errors.New("nil_word")
	} else if w.UserId != uid {
		return errors.New("db_delete_refuse")
	}
	// mysql
	err = mysql.DelWord(w)
	if err != nil {
		return err
	}
	// 同步
	go func() {
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_DELETE, entity.TRENDS_CON_TYPE_WORD, wid)
		AddTrends(trends)
	}()
	return err
}

// GetWordListByCouple
func GetWordListByCouple(uid, cid int64, page int) ([]*entity.Word, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Word
	offset := page * limit
	list, err := mysql.GetWordListByCouple(cid, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_word")
		} else {
			return nil, nil
		}
	}
	if page > 0 {
		return list, err
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_WORD)
		AddTrends(trends)
	}()
	return list, err
}
