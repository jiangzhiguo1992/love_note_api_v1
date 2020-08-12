package services

import (
	"errors"
	"libs/utils"
	"math"
	"models/entity"
	"models/mysql"
	"models/redis"
	"strconv"
	"strings"
	"time"
)

// IsUserBlack用户是否被拉黑
func IsUserBlack(u *entity.User) bool {
	if u == nil {
		return false
	}
	return u.Status <= entity.STATUS_DELETE
}

// IsUserInfoComplete 检查user信息是否完善
func IsUserInfoComplete(u *entity.User) bool {
	if u.Sex != entity.USER_SEX_GIRL && u.Sex != entity.USER_SEX_BOY {
		return false
	}
	// 生日就不检查了，毕竟有人1970年出生
	return true
}

// IsUserInCouple user是否在couple中
func IsUserInCouple(uid int64, c *entity.Couple) bool {
	if c == nil || uid <= 0 {
		return false
	}
	b := (c.CreatorId == uid) || (c.InviteeId == uid)
	return b
}

// GetCoupleIdByUser 获取user的cid
func GetCoupleIdByUser(u *entity.User) int64 {
	if u != nil && u.Couple != nil {
		return u.Couple.Id
	}
	return 0
}

// GetTaId 获取user中的couple的ta的id
func GetTaId(u *entity.User) int64 {
	if u == nil {
		return 0
	}
	couple := u.Couple
	if couple == nil {
		return 0
	}
	var taId int64
	if couple.CreatorId == u.Id {
		taId = couple.InviteeId
	} else {
		taId = couple.CreatorId
	}
	return taId
}

// AddUser 密码注册
func AddUserByPwd(u *entity.User) (*entity.User, error) {
	if u == nil {
		return nil, errors.New("nil_user")
	} else if len(strings.TrimSpace(u.Phone)) != PHONE_LENGTH {
		return nil, errors.New("limit_phone_err")
	} else if len([]rune(u.Password)) <= 0 {
		return nil, errors.New("user_pwd_nil")
	}
	// phone查重
	dbUser, err := mysql.GetUserByPhone(u.Phone)
	if err != nil {
		return nil, err
	} else if dbUser != nil {
		return nil, errors.New("user_phone_exist")
	}
	// mysql
	u.UserToken = createUserToken()
	u, err = mysql.AddUser(u)
	if u == nil || err != nil {
		return nil, err
	}
	// redis
	redis.SetUser(u)
	return u, err
}

// AddUser 验证码注册
func AddUserByVer(u *entity.User) (*entity.User, error) {
	if u == nil {
		return nil, errors.New("nil_user")
	} else if len(strings.TrimSpace(u.Phone)) != PHONE_LENGTH {
		return nil, errors.New("limit_phone_err")
	}
	// phone查重
	dbUser, err := mysql.GetUserByPhone(u.Phone)
	if err != nil {
		return nil, err
	} else if dbUser != nil {
		return nil, errors.New("user_phone_exist")
	}
	// mysql
	u.Password = ""
	u.UserToken = createUserToken()
	u, err = mysql.AddUser(u)
	if u == nil || err != nil {
		return nil, err
	}
	// redis
	redis.SetUser(u)
	return u, err
}

// UpdateUserOnLogin 登录
// 1.user全属性
func UpdateUserOnLogin(u *entity.User) (*entity.User, error) {
	if u == nil || u.Id <= 0 {
		return u, errors.New("nil_user")
	}
	// redis-del
	redis.DelUser(u)
	// mysql
	u.UserToken = createUserToken()
	u, err := mysql.UpdateUser(u)
	if u == nil || err != nil {
		return nil, errors.New("user_login_fail")
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// UpdateUserPhone 修改手机
// 1.user全属性
func UpdateUserPhone(u *entity.User, oldPhone string) (*entity.User, error) {
	if u == nil || u.Id <= 0 {
		return u, errors.New("nil_user")
	} else if len(strings.TrimSpace(u.Phone)) != PHONE_LENGTH {
		return u, errors.New("limit_phone_err")
	}
	// phone查重
	dbUser, err := mysql.GetUserByPhone(u.Phone)
	if err != nil {
		return dbUser, err
	} else if dbUser != nil {
		return dbUser, errors.New("user_phone_exist")
	}
	// redis-del
	old := &entity.User{
		BaseObj: entity.BaseObj{
			Id: u.Id,
		},
		Phone:     oldPhone,
		UserToken: u.UserToken,
	}
	redis.DelUser(old)
	// data
	u.UserToken = createUserToken()
	// mysql
	u, err = mysql.UpdateUser(u)
	if u == nil || err != nil {
		return nil, err
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// UpdateUserPwd 修改密码
// 1.user全属性
func UpdateUserPwd(u *entity.User) (*entity.User, error) {
	if u == nil || u.Id <= 0 {
		return u, errors.New("nil_user")
	} else if len([]rune(u.Password)) <= 0 {
		return u, errors.New("user_pwd_nil")
	}
	// redis-del
	redis.DelUser(u)
	// data
	u.UserToken = createUserToken()
	// mysql
	u, err := mysql.UpdateUser(u)
	if u == nil || err != nil {
		return nil, err
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// UpdateUserInfo 修改信息
// 1.user全属性
func UpdateUserInfo(u *entity.User) (*entity.User, error) {
	if u == nil || u.Id <= 0 {
		return u, errors.New("nil_user")
	} else if u.Sex != entity.USER_SEX_GIRL && u.Sex != entity.USER_SEX_BOY {
		return u, errors.New("user_sex_nil")
	} else if u.Birthday == 0 {
		// 1970年之前的小于0
		return u, errors.New("user_birth_nil")
	} else if u.Birthday > time.Now().Unix() {
		// 不能超过现在
		return u, errors.New("user_birth_over")
	} else if time.Unix(u.Birthday, 0).Add(time.Hour * 24 * 365 * 120).Unix() < time.Now().Unix() {
		// 不能大于120岁
		return u, errors.New("user_birth_over")
	}
	// redis-del
	redis.DelUser(u)
	// mysql
	u, err := mysql.UpdateUser(u)
	if u == nil || err != nil {
		return nil, err
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// ToggleUserStatus
func ToggleUserStatus(u *entity.User) (*entity.User, error) {
	if u == nil || u.Id <= 0 {
		return u, errors.New("nil_user")
	}
	// redis-del
	redis.DelUser(u)
	// mysql
	if u.Status == entity.STATUS_DELETE {
		u.Status = entity.STATUS_VISIBLE
	} else {
		u.Status = entity.STATUS_DELETE
	}
	u.UserToken = createUserToken()
	u, err := mysql.UpdateUserStatus(u)
	if u == nil || err != nil {
		return nil, errors.New("user_login_fail")
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// GetUserByToken 应用于用户验证校对
// 1.返回全status + 全信息
func GetUserByToken(token string) (*entity.User, error) {
	if len(token) <= 0 {
		return nil, errors.New("user_token_nil")
	}
	// redis-get
	u, _ := redis.GetUserByToken(token)
	if u != nil && u.Id > 0 {
		return u, nil
	}
	// mysql
	u, err := mysql.GetUserByToken(token)
	if u == nil || err != nil {
		return nil, err
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// GetUserByPhone 登录调用、注册检查、修改检查
// 1.返回全status + 全信息
func GetUserByPhone(phone string) (*entity.User, error) {
	if len(phone) <= 0 {
		// 只检查长度
		return nil, errors.New("limit_phone_err")
	}
	// redis-get
	u, _ := redis.GetUserByPhone(phone)
	if u != nil && u.Id > 0 {
		return u, nil
	}
	// mysql
	u, err := mysql.GetUserByPhone(phone)
	if u == nil || err != nil {
		return nil, err
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// GetUserById 所有用户
// 1.返回全status + 全信息
func GetUserById(uid int64) (*entity.User, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	}
	// redis-get
	u, _ := redis.GetUserById(uid)
	if u != nil && u.Id > 0 {
		return u, nil
	}
	// mysql
	u, err := mysql.GetUserById(uid)
	if u == nil || err != nil {
		return nil, err
	}
	// redis-set
	redis.SetUser(u)
	return u, err
}

// GetUserList
func GetUserList(page int) ([]*entity.User, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().User
	offset := page * limit
	list, err := mysql.GetUserList(offset, limit)
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

// GetUserListByBlack
func GetUserListByBlack(page int) ([]*entity.User, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().User
	offset := page * limit
	list, err := mysql.GetUserListByBlack(offset, limit)
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

// GetUserTotalByCreate
func GetUserTotalByCreate(start, end int64) (int64) {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetUserTotalByCreate(start, end)
	return total
}

// GetUserBirthAvgByCreate
func GetUserBirthAvgByCreate(start, end int64) (float64) {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetUserBirthAvgByCreate(start, end)
	return total
}

// createUserToken timestamp + 随机数
func createUserToken() string {
	// 前半截时间戳
	unixNa := time.Now().UnixNano()
	unix := strconv.FormatInt(unixNa, 16)
	// 后半截随机数
	length := 16
	max := math.Pow10(length) - 1
	min := math.Pow10(length - 1)
	rand16 := utils.GetRandRange(int(max), int(min))
	//rand16 := utils.GetRandRange(9999999999999999, 1000000000000000)
	return unix + "-" + strconv.Itoa(rand16)
}
