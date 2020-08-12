package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"libs/utils"

	"github.com/gomodule/redigo/redis"
	"models/entity"
)

// SetCoupleByUser
func SetCoupleByUser(uid int64, couple *entity.Couple) error {
	// 解析user
	if couple == nil || couple.Id <= 0 {
		utils.LogWarn("SetCoupleByUser", "无效的配对: "+fmt.Sprintf("%+v", couple))
		return errors.New("nil_couple")
	}
	// 让你存，反正我app端做了判断了
	//if couple != nil && couple.State != nil && couple.State.State != COUPLE_STATE_520 {
	//	// cp不可见状态就不往redis里存了，正在分手也不存
	//	couple = nil
	//}
	if uid <= 0 {
		utils.LogWarn("SetCoupleByUser", "无效的uid: "+fmt.Sprintf("%+v", couple))
		return errors.New("nil_user")
	}
	bytes, err := json.Marshal(couple)
	if err != nil {
		utils.LogErr("SetCoupleByUser", err)
		// 添加失败就删除
		DelCoupleByUser(uid)
		return errors.New("data_decode_err")
	}
	coupleBody := string(bytes)
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
	// 存储用户
	uidStr := strconv.FormatInt(uid, 10)
	_, err = conn.Do("SET", KEY_UID_COUPLE+uidStr, coupleBody)
	if err != nil {
		utils.LogErr("SetCoupleByUser", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_UID_COUPLE+uidStr, getRedisCoupleExpire())
	utils.LogErr("SetCoupleByUser", err)
	return err
}

// DelCoupleByUser
func DelCoupleByUser(uid int64) error {
	if uid <= 0 {
		utils.LogWarn("DelCoupleByUser", "uid <= 0")
		return errors.New("nil_user")
	}
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
	// 开始删除
	_, err := conn.Do("DEL", KEY_UID_COUPLE+strconv.FormatInt(uid, 10))
	utils.LogErr("DelCoupleByUser", err)
	return err
}

// GetCoupleByUser
func GetCoupleByUser(uid int64) (*entity.Couple, error) {
	if uid <= 0 {
		utils.LogWarn("GetCoupleByUser", "uid <= 0")
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
	reply, err := conn.Do("GET", KEY_UID_COUPLE+strconv.FormatInt(uid, 10))
	if err != nil {
		utils.LogErr("couple_redis", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析user
	couple := &entity.Couple{}
	err = json.Unmarshal(bytes, couple)
	if err != nil {
		utils.LogErr("couple_redis", err)
		return nil, err
	}
	return couple, nil
}
