package services

import (
	"errors"
	"models/entity"
	"models/mysql"
	"strings"
)

// AddVersion
func AddVersion(v *entity.Version) (*entity.Version, error) {
	if v == nil || v.VersionCode <= 0 {
		return nil, errors.New("nil_version")
	} else if len(strings.TrimSpace(v.VersionName)) <= 0 ||
		len(strings.TrimSpace(v.UpdateUrl)) <= 0 ||
		len(strings.TrimSpace(v.UpdateLog)) <= 0 {
		return nil, errors.New("limit_content_text_nil")
	} else if len(strings.TrimSpace(v.Platform)) <= 0 {
		return nil, errors.New("nil_platform")
	}
	// mysql
	v, err := mysql.AddVersion(v)
	return v, err
}

// DelVersion
func DelVersion(vid int64) error {
	if vid <= 0 {
		return errors.New("nil_version")
	}
	// 旧数据检查
	v, err := mysql.GetVersionById(vid)
	if err != nil {
		return err
	} else if v == nil {
		return errors.New("nil_version")
	}
	// mysql
	err = mysql.DelVersion(v)
	return err
}

// GetVersionList
func GetVersionList(page int) ([]*entity.Version, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Version
	offset := page * limit
	list, err := mysql.GetVersionList(offset, limit)
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

// GetVersionListByCode
func GetVersionListByCode(uid int64, code int) ([]*entity.Version, error) {
	if code <= 0 {
		return nil, nil
	}
	// platform
	entry, err := mysql.GetEntryLatestByUser(uid)
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, errors.New("nil_entry")
	} else if len(entry.Platform) <= 0 {
		return nil, errors.New("nil_platform")
	}
	// mysql
	list, err := mysql.GetVersionListByCode(entry.Platform, code)
	return list, err
}

// GetVersionListByPlatformCode
func GetVersionListByPlatformCode(platform string, code int) ([]*entity.Version, error) {
	if len(platform) <= 0 {
		return nil, errors.New("nil_platform")
	} else if code <= 0 {
		return nil, nil
	}
	platform = strings.ToLower(strings.TrimSpace(platform))
	// mysql
	list, err := mysql.GetVersionListByCode(platform, code)
	return list, err
}
