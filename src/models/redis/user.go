package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"libs/utils"

	"github.com/gomodule/redigo/redis"
	"models/entity"
)

// SetUser
func SetUser(user *entity.User) error {
	if user == nil {
		utils.LogWarn("SetUser", "无效的用户: "+fmt.Sprintf("%+v", user))
		return errors.New("nil_user")
	}
	// 分别存三份
	uid := user.Id
	phone := strings.TrimSpace(user.Phone)
	token := strings.TrimSpace(user.UserToken)
	bytes, err := json.Marshal(user)
	if err != nil {
		utils.LogErr("SetUser", err)
		// 添加失败就删除
		DelUser(user)
		return errors.New("data_decode_err")
	}
	userBody := string(bytes)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 存储用户-id
	if uid > 0 {
		_, err = conn.Do("SET", KEY_ID_USER+strconv.FormatInt(uid, 10), userBody)
		if err != nil {
			utils.LogErr("SetUser", err)
			return err
		}
		// 设置过期
		_, err = conn.Do("EXPIRE", KEY_ID_USER+strconv.FormatInt(uid, 10), getRedisUserExpire())
		utils.LogErr("SetUser", err)
	} else {
		utils.LogWarn("SetUser", "uid <=0")
	}
	// 存储用户-phone
	if len(phone) > 0 {
		_, err = conn.Do("SET", KEY_PHONE_USER+phone, userBody)
		if err != nil {
			utils.LogErr("SetUser", err)
			return err
		}
		// 设置过期
		_, err = conn.Do("EXPIRE", KEY_PHONE_USER+phone, getRedisUserExpire())
		utils.LogErr("SetUser", err)
	} else {
		utils.LogWarn("SetUser", "len(phone) <=0")
	}
	// 存储用户-token
	if len(token) > 0 {
		_, err = conn.Do("SET", KEY_TOKEN_USER+token, userBody)
		if err != nil {
			utils.LogErr("SetUser", err)
			return err
		}
		// 设置过期
		_, err = conn.Do("EXPIRE", KEY_TOKEN_USER+token, getRedisUserExpire())
		utils.LogErr("SetUser", err)
	} else {
		//utils.LogWarn("RedisSetUser", "len(token) <=0")
	}
	return err
}

// DelUser
func DelUser(user *entity.User) error {
	if user == nil {
		utils.LogWarn("DelUser", "无效的用户: "+fmt.Sprintf("%+v", user))
		return errors.New("nil_user")
	}
	// 分别删三份
	uid := user.Id
	phone := strings.TrimSpace(user.Phone)
	token := strings.TrimSpace(user.UserToken)
	// 开始连接
	if pool == nil {
		return errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis_conn_nil")
	}
	auth(conn)
	// 开始删除-id
	var err error
	if uid > 0 {
		_, err = conn.Do("DEL", KEY_ID_USER+strconv.FormatInt(uid, 10))
		utils.LogErr("DelUser", err)
	} else {
		utils.LogWarn("DelUser", "uid <=0")
	}
	// 开始删除-phone
	if len(phone) > 0 {
		_, err = conn.Do("DEL", KEY_PHONE_USER+phone)
		utils.LogErr("DelUser", err)
	} else {
		utils.LogWarn("DelUser", "len(phone) <=0")
	}
	// 开始删除-token
	if len(token) > 0 {
		_, err := conn.Do("DEL", KEY_TOKEN_USER+token)
		utils.LogErr("DelUser", err)
	} else {
		//utils.LogWarn("RedisDelUser", "len(token) <=0")
	}
	return err
}

// GetUserById
func GetUserById(id int64) (*entity.User, error) {
	if id <= 0 {
		utils.LogWarn("GetUserById", "id <= 0")
		return nil, errors.New("nil_user")
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_ID_USER+strconv.FormatInt(id, 10))
	if err != nil {
		utils.LogErr("GetUserById", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析user
	user := &entity.User{}
	err = json.Unmarshal(bytes, user)
	if err != nil {
		utils.LogErr("GetUserById", err)
		return nil, err
	}
	return user, nil
}

// GetUserByPhone
func GetUserByPhone(phone string) (*entity.User, error) {
	if len(phone) <= 0 {
		utils.LogWarn("GetUserByPhone", "phone <= 0")
		return nil, errors.New("limit_phone_err")
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_PHONE_USER+phone)
	if err != nil {
		utils.LogErr("GetUserByPhone", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析user
	user := &entity.User{}
	err = json.Unmarshal(bytes, user)
	if err != nil {
		utils.LogErr("GetUserByPhone", err)
		return nil, err
	}
	return user, nil
}

// GetUserByToken
func GetUserByToken(token string) (*entity.User, error) {
	token = strings.TrimSpace(token)
	if len(token) <= 0 {
		utils.LogWarn("GetUserByToken", "token <= 0")
		return nil, errors.New("user_token_nil")
	}
	// 开始连接
	if pool == nil {
		return nil, errors.New("redis_nil")
	}
	conn := pool.Get()
	defer conn.Close()
	if conn == nil {
		return nil, errors.New("redis_conn_nil")
	}
	auth(conn)
	// 获取用户
	reply, err := conn.Do("GET", KEY_TOKEN_USER+token)
	if err != nil {
		utils.LogErr("GetUserByToken", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析user
	user := &entity.User{}
	err = json.Unmarshal(bytes, user)
	if err != nil {
		utils.LogErr("GetUserByToken", err)
		return nil, err
	}
	return user, nil
}
