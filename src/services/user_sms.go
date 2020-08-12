package services

import (
	"errors"
	"libs/utils"
	"math"
	"models/entity"
	"models/mysql"
	"strconv"
	"strings"
	"time"
)

// CreateSmsValidateCode 生成6位验证码
func CreateSmsValidateCode() string {
	codeLength := GetLimit().SmsCodeLength
	max := math.Pow10(codeLength) - 1
	min := math.Pow10(codeLength - 1)
	code := utils.GetRandRange(int(max), int(min))
	return strconv.Itoa(code) // 100000-999999
}

// AddSms
// 1.外部做防黑措施
func AddSms(s *entity.Sms) (*entity.Sms, error) {
	if s == nil {
		return s, errors.New("nil_sms")
	} else if len(strings.TrimSpace(s.Phone)) != PHONE_LENGTH {
		return nil, errors.New("limit_phone_err")
	} else if s.SendType <= 0 {
		return nil, errors.New("sms_type_nil")
	}
	// mysql
	s, err := mysql.AddSms(s)
	return s, err
}

// SmsEnableSend 是否可以发送短信
func SmsEnableSend(phone string, sendType int) error {
	if len(strings.TrimSpace(phone)) != PHONE_LENGTH {
		return errors.New("limit_phone_err")
	} else if sendType <= 0 {
		return errors.New("sms_type_nil")
	}
	// 检查最多发送次数
	limit := GetLimit()
	maxCount := limit.SmsMaxCount
	// mysql
	list, err := mysql.GetSmsList(phone, sendType, 0, maxCount)
	if err != nil {
		return err
	} else if list != nil && len(list) > 0 {
		now := time.Now().Unix()
		// 相连sms间隔
		currentBetween := now - list[0].CreateAt
		if currentBetween < int64(limit.SmsBetweenSec) {
			//sec := utils.ConvertInt642String(shouldBetween*60-currentBetween)
			return errors.New("sms_twice_send_too_frequent")
		}
		// 最大max间隔
		shouldStart := list[len(list)-1].CreateAt + int64(limit.SmsMaxSec)
		if len(list) >= maxCount && shouldStart > now {
			//sec := utils.ConvertInt642String(shouldBetween*60-currentBetween)
			return errors.New("sms_send_frequent")
		}
	}
	return nil
}

// SmsCheckCode 检查code是否正确
func SmsCheckCode(phone string, sendType int, code string) error {
	limit := GetLimit()
	if len(strings.TrimSpace(phone)) != PHONE_LENGTH {
		return errors.New("limit_phone_err")
	} else if sendType <= 0 {
		return errors.New("sms_type_nil")
	} else if len(code) != limit.SmsCodeLength {
		return errors.New("sms_code_length_nil")
	}
	s, err := mysql.GetSmsByPhoneType(phone, sendType)
	if err != nil {
		return err
	} else if s == nil || s.Id <= 0 {
		return errors.New("sms_code_nil")
	} else if s.CreateAt < time.Now().Unix()-int64(limit.SmsEffectSec) {
		return errors.New("sms_code_nil")
	} else if strings.TrimSpace(s.Content) != strings.TrimSpace(code) {
		return errors.New("sms_code_err")
	}
	return nil
}

// GetSmsListByCreate
func GetSmsListByCreate(start, end int64, phone string, sendType, page int) ([]*entity.Sms, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Sms
	offset := page * limit
	list, err := mysql.GetSmsList(phone, sendType, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_common")
		} else {
			return nil, nil
		}
	}
	return list, err
}

// GetSmsTotalByCreateWithDel
func GetSmsTotalByCreateWithDel(start, end int64, phone string, sendType int) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetSmsTotalByCreateWithDel(start, end, phone, sendType)
	return total
}
