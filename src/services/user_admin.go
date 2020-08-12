package services

import (
	"models/entity"
	"models/mysql"
)

// IsAdminister
func IsAdminister(u *entity.User) bool {
	if u == nil || u.Id <= 0 {
		return false
	}
	// mysql
	a, err := mysql.GetAdminByUser(u.Id)
	if err != nil {
		return false
	} else if a == nil || a.Id <= 0 {
		return false
	} else if a.Permission <= entity.ADMIN_PERMISSION_NO {
		return false
	}
	return true
}
