package services

import (
	"errors"
	"fmt"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddEntry
func AddEntry(e *entity.Entry) (*entity.Entry, error) {
	if e == nil {
		return nil, errors.New("nil_entry")
	} else if e.UserId <= 0 {
		utils.LogErr("AddEntry", "缺失重要数据 "+fmt.Sprintf("%+v", e))
		return nil, errors.New("nil_user")
	}
	e.DeviceName = strings.TrimSpace(e.DeviceName)
	e.Market = strings.ToLower(strings.TrimSpace(e.Market))
	e.Language = strings.ToLower(strings.TrimSpace(e.Language))
	e.Platform = strings.ToLower(strings.TrimSpace(e.Platform))
	// mysql
	old, err := mysql.GetEntryLatestByUser(e.UserId)
	if err != nil {
		return nil, err
	}
	if old == nil || old.Id <= 0 {
		// 新加
		e, err = mysql.AddEntry(e)
	} else {
		// 更新
		old.DeviceId = e.DeviceId
		old.DeviceName = e.DeviceName
		old.Market = e.Market
		old.Language = e.Language
		old.Platform = e.Platform
		old.OsVersion = e.OsVersion
		old.AppVersion = e.AppVersion
		e, err = mysql.UpdateEntry(old)
	}
	return e, err
}

// GetEntryList
func GetEntryList(uid int64, page int) ([]*entity.Entry, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Entry
	offset := page * limit
	list, err := mysql.GetEntryList(uid, offset, limit)
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

// GetEntryTotalByCreate 时间活跃度
func GetEntryTotalByCreate(start, end int64) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetEntryTotalByCreateWithDel(start, end)
	return total
}

// GetEntryTotalByUpdate 时间活跃度
func GetEntryTotalByUpdate(start, end int64) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetEntryTotalByUpdateWithDel(start, end)
	return total
}

// GetEntryGroupDeviceNameList 设备名称占比
func GetEntryGroupDeviceNameList(at string, start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetEntryGroupListByFiled("device_name", at, start, end)
}

// GetEntryGroupMarketList 应用市场占比
func GetEntryGroupMarketList(at string, start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetEntryGroupListByFiled("market", at, start, end)
}

// GetEntryGroupLanguageList 语言占比
func GetEntryGroupLanguageList(at string, start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetEntryGroupListByFiled("language", at, start, end)
}

// GetEntryGroupPlatformList 移动平台占比
func GetEntryGroupPlatformList(at string, start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetEntryGroupListByFiled("platform", at, start, end)
}

// GetEntryGroupOsVersionList 操作系统占比
func GetEntryGroupOsVersionList(at string, start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetEntryGroupListByFiled("os_version", at, start, end)
}

// GetEntryGroupAppVersionList 软件版本占比
func GetEntryGroupAppVersionList(at string, start, end int64) ([]*entity.FiledInfo, error) {
	if start >= end {
		return nil, errors.New("limit_happen_err")
	}
	return mysql.GetEntryGroupListByFiled("app_version", at, start, end)
}
