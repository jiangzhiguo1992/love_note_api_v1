package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"libs/utils"
	"models/entity"
	"sort"
	"strconv"
)

// removeRepeatedMatchWork
func removeRepeatedMatchWork(list []*entity.MatchWork) ([]*entity.MatchWork) {
	returnList := make([]*entity.MatchWork, 0)
	if list == nil || len(list) <= 0 {
		return returnList
	}
	// 遍历所有元素
	for i := 0; i < len(list); i++ {
		iWork := list[i]
		if iWork == nil || iWork.Id <= 0 {
			continue
		}
		// 检查重复元素 0=10 1=11 2=12 3=11 4=12 5=10 6=12
		repeat := false
		for j := i + 1; j < len(list); j++ {
			jWork := list[j]
			if jWork == nil || jWork.Id <= 0 {
				continue
			}
			if iWork.Id == jWork.Id {
				// 更新较晚的加上，相同的话 加后面的
				if iWork.UpdateAt <= jWork.UpdateAt {
					repeat = true
				} else if iWork.UpdateAt > jWork.UpdateAt {
					repeat = false
				}
				// 防止多重
				if repeat {
					break
				}
			}
		}
		if !repeat {
			// 并且去掉user信息
			iWork.Mine = false
			iWork.Our = false
			iWork.Report = false
			iWork.Point = false
			iWork.Coin = false
			returnList = append(returnList, iWork)
		}
	}
	return returnList
}

// SetMatchWork
func SetMatchWork(work *entity.MatchWork) error {
	if work == nil || work.Id <= 0 {
		utils.LogWarn("SetMatchWork", "无效的作品: "+fmt.Sprintf("%+v", work))
		return errors.New("nil_work")
	}
	work.Mine = false
	work.Our = false
	work.Report = false
	work.Point = false
	work.Coin = false
	bytes, err := json.Marshal(work)
	if err != nil {
		utils.LogErr("SetMatchWork", err)
		// 添加失败就删除
		DelMatchWork(work)
		return errors.New("data_decode_err")
	}
	workBody := string(bytes)
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
	// 存储work
	if work.Id > 0 {
		_, err = conn.Do("SET", KEY_ID_WORK+strconv.FormatInt(work.Id, 10), workBody)
		if err != nil {
			utils.LogErr("SetMatchWork", err)
			return err
		}
		_, err = conn.Do("EXPIRE", KEY_ID_WORK+strconv.FormatInt(work.Id, 10), getRedisMatchWorkExpire())
		utils.LogErr("SetMatchWork", err)
	} else {
		utils.LogWarn("SetMatchWork", "mwid <=0")
	}
	return err
}

// DelMatchWork
func DelMatchWork(work *entity.MatchWork) error {
	if work == nil || work.Id <= 0 {
		utils.LogWarn("DelMatchWork", "无效的作品: "+fmt.Sprintf("%+v", work))
		return errors.New("nil_work")
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
	// 开始删除-id
	var err error
	if work.Id > 0 {
		_, err = conn.Do("DEL", KEY_ID_WORK+strconv.FormatInt(work.Id, 10))
		utils.LogErr("DelMatchWork", err)
	} else {
		utils.LogWarn("DelMatchWork", "mwid <=0")
	}
	return err
}

// GetMatchWorkById
func GetMatchWorkById(mwid int64) (*entity.MatchWork, error) {
	if mwid <= 0 {
		utils.LogWarn("GetMatchWorkById", "mwid <= 0")
		return nil, errors.New("nil_work")
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
	reply, err := conn.Do("GET", KEY_ID_WORK+strconv.FormatInt(mwid, 10))
	if err != nil {
		utils.LogErr("GetMatchWorkById", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析work
	work := &entity.MatchWork{}
	err = json.Unmarshal(bytes, work)
	if err != nil {
		utils.LogErr("GetMatchWorkById", err)
		return nil, err
	}
	return work, nil
}

//------------------------------------CoinList------------------------------------

// DelMatchWorkInCoinList
func DelMatchWorkInCoinList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("DelMatchWorkInCoinList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧数据
	workList, _ := GetMatchWorkListCoinByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.MatchWork, 0)
	for _, v := range workList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != mw.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetMatchWorkListCoinByKind(period, mw.Kind, saveList, false)
	return err
}

// UpdateMatchWorkInCoinList
func UpdateMatchWorkInCoinList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("UpdateMatchWorkInCoinList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧列表
	workList, _ := GetMatchWorkListCoinByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		return nil
	}
	// 添加新元素
	workList = append(workList, mw)
	// 重新保存
	err := SetMatchWorkListCoinByKind(period, mw.Kind, workList, false)
	return err
}

// SetMatchWorkListCoinByKind
// merge 是否合并已有的数据
func SetMatchWorkListCoinByKind(period, kind int, list []*entity.MatchWork, merge bool) error {
	if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.MatchWork, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisWorkList, _ := GetMatchWorkListCoinByKind(period, kind, -1, -1)
		if redisWorkList != nil && len(redisWorkList) > 0 {
			for _, v := range redisWorkList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedMatchWork(returnList)
	// 开始降序排序
	sort.Sort(workCoinSlice(returnList))
	// 构造数据
	saveList := &matchWorkList{
		MatchWorkList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetMatchWorkListCoinByKind", err)
		// list添加失败不删除
		if merge {

		}
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
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
	// 存储work
	_, err = conn.Do("SET", KEY_PERIOD_KIND_COIN_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind), listBody)
	if err != nil {
		utils.LogErr("SetMatchWorkListCoinByKind", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_PERIOD_KIND_COIN_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind), getRedisMatchWorkExpire())
	utils.LogErr("SetMatchWorkListCoinByKind", err)
	return nil
}

// GetMatchWorkListCoinByKind
func GetMatchWorkListCoinByKind(period, kind, offset, limit int) ([]*entity.MatchWork, error) {
	if limit == 0 {
		utils.LogWarn("GetMatchWorkListCoinByKind", "limit == 0")
		return nil, nil
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
	reply, err := conn.Do("GET", KEY_PERIOD_KIND_COIN_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind))
	if err != nil {
		utils.LogErr("GetMatchWorkListCoinByKind", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析work
	list := &matchWorkList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetMatchWorkListCoinByKind", err)
		return nil, err
	}
	workList := list.MatchWorkList
	// 数量不够
	if workList == nil || len(workList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return workList, nil
	}
	// 正常获取，数量够
	returnList := workList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------PointList------------------------------------

// DelMatchWorkInPointList
func DelMatchWorkInPointList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("DelMatchWorkInPointList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧数据
	workList, _ := GetMatchWorkListPointByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.MatchWork, 0)
	for _, v := range workList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != mw.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetMatchWorkListPointByKind(period, mw.Kind, saveList, false)
	return err
}

// UpdateMatchWorkInPointList
func UpdateMatchWorkInPointList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("UpdateMatchWorkInPointList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧列表
	workList, _ := GetMatchWorkListPointByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		return nil
	}
	// 添加新元素
	workList = append(workList, mw)
	// 重新保存
	err := SetMatchWorkListPointByKind(period, mw.Kind, workList, false)
	return err
}

// SetMatchWorkListPointByKind
// merge 是否合并已有的数据
func SetMatchWorkListPointByKind(period, kind int, list []*entity.MatchWork, merge bool) error {
	if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.MatchWork, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisWorkList, _ := GetMatchWorkListPointByKind(period, kind, -1, -1)
		if redisWorkList != nil && len(redisWorkList) > 0 {
			for _, v := range redisWorkList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedMatchWork(returnList)
	// 开始降序排序
	sort.Sort(workPointSlice(returnList))
	// 构造数据
	saveList := &matchWorkList{
		MatchWorkList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetMatchWorkListPointByKind", err)
		// list添加失败不删除
		if merge {

		}
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
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
	// 存储work
	_, err = conn.Do("SET", KEY_PERIOD_KIND_POINT_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind), listBody)
	if err != nil {
		utils.LogErr("SetMatchWorkListPointByKind", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_PERIOD_KIND_POINT_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind), getRedisMatchWorkExpire())
	utils.LogErr("SetMatchWorkListPointByKind", err)
	return nil
}

// GetMatchWorkListPointByKind
func GetMatchWorkListPointByKind(period, kind, offset, limit int) ([]*entity.MatchWork, error) {
	if limit == 0 {
		utils.LogWarn("GetMatchWorkListPointByKind", "limit == 0")
		return nil, nil
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
	reply, err := conn.Do("GET", KEY_PERIOD_KIND_POINT_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind))
	if err != nil {
		utils.LogErr("GetMatchWorkListPointByKind", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析work
	list := &matchWorkList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetMatchWorkListPointByKind", err)
		return nil, err
	}
	workList := list.MatchWorkList
	// 数量不够
	if workList == nil || len(workList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return workList, nil
	}
	// 正常获取，数量够
	returnList := workList[offset : offset+limit]
	return returnList, nil
}

//------------------------------------CreateList------------------------------------

// AddMatchWorkInCreateList
func AddMatchWorkInCreateList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("AddMatchWorkInCreateList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧数据
	workList, _ := GetMatchWorkListCreateByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		workList = make([]*entity.MatchWork, 0)
	}
	// 添加指定元素
	workList = append(workList, mw)
	// 重新保存
	err := SetMatchWorkListCreateByKind(period, mw.Kind, workList, false)
	return err
}

// DelMatchWorkInCreateList
func DelMatchWorkInCreateList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("DelMatchWorkInCreateList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧数据
	workList, _ := GetMatchWorkListCreateByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		return nil
	}
	// 删除指定元素
	saveList := make([]*entity.MatchWork, 0)
	for _, v := range workList {
		if v == nil || v.Id <= 0 {
			continue
		}
		if v.Id != mw.Id {
			saveList = append(saveList, v)
		}
	}
	// 重新保存
	err := SetMatchWorkListCreateByKind(period, mw.Kind, saveList, false)
	return err
}

// UpdateMatchWorkInCreateList
func UpdateMatchWorkInCreateList(period int, mw *entity.MatchWork) error {
	if mw == nil || mw.Id <= 0 {
		utils.LogWarn("UpdateMatchWorkInCreateList", "work == nil")
		return errors.New("nil_work")
	}
	// 获取旧列表
	workList, _ := GetMatchWorkListCreateByKind(period, mw.Kind, -1, -1)
	if workList == nil || len(workList) <= 0 {
		return nil
	}
	// 添加新元素
	workList = append(workList, mw)
	// 重新保存
	err := SetMatchWorkListCreateByKind(period, mw.Kind, workList, false)
	return err
}

// SetMatchWorkListCreateByKind
// merge 是否合并已有的数据
func SetMatchWorkListCreateByKind(period, kind int, list []*entity.MatchWork, merge bool) error {
	if list == nil || len(list) <= 0 {
		return nil
	}
	// 合并list
	returnList := make([]*entity.MatchWork, 0)
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		returnList = append(returnList, v)
	}
	// 合并redisList，并过滤重复元素
	if merge {
		redisWorkList, _ := GetMatchWorkListCreateByKind(period, kind, -1, -1)
		if redisWorkList != nil && len(redisWorkList) > 0 {
			for _, v := range redisWorkList {
				if v == nil || v.Id <= 0 {
					continue
				}
				repeat := false
				for _, n := range returnList {
					if n == nil || n.Id <= 0 {
						continue
					}
					if v.Id == n.Id && v.UpdateAt == n.UpdateAt {
						repeat = true
					}
				}
				if !repeat {
					returnList = append(returnList, v)
				}
			}
		}
	}
	// 去重，并剔除user数据
	returnList = removeRepeatedMatchWork(returnList)
	// 开始降序排序
	sort.Sort(workCreateSlice(returnList))
	// 构造数据
	saveList := &matchWorkList{
		MatchWorkList: returnList,
	}
	bytes, err := json.Marshal(saveList)
	if err != nil {
		utils.LogErr("SetMatchWorkListCreateByKind", err)
		// list添加失败不删除
		if merge {

		}
		return errors.New("data_decode_err")
	}
	listBody := string(bytes)
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
	// 存储work
	_, err = conn.Do("SET", KEY_PERIOD_KIND_CREATE_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind), listBody)
	if err != nil {
		utils.LogErr("SetMatchWorkListCreateByKind", err)
		return err
	}
	// 设置过期
	_, err = conn.Do("EXPIRE", KEY_PERIOD_KIND_CREATE_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind), getRedisMatchWorkExpire())
	utils.LogErr("SetMatchWorkListCreateByKind", err)
	return nil
}

// GetMatchWorkListCreateByKind
func GetMatchWorkListCreateByKind(period, kind, offset, limit int) ([]*entity.MatchWork, error) {
	if limit == 0 {
		utils.LogWarn("GetMatchWorkListCreateByKind", "limit == 0")
		return nil, nil
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
	reply, err := conn.Do("GET", KEY_PERIOD_KIND_CREATE_WORK_LIST+strconv.Itoa(period)+":"+strconv.Itoa(kind))
	if err != nil {
		utils.LogErr("GetMatchWorkListCreateByKind", err)
		return nil, err
	}
	bytes, err := redis.Bytes(reply, err)
	if err != nil {
		// 无用户，不打印
		return nil, err
	}
	// 解析work
	list := &matchWorkList{}
	err = json.Unmarshal(bytes, list)
	if err != nil {
		utils.LogErr("GetMatchWorkListCreateByKind", err)
		return nil, err
	}
	workList := list.MatchWorkList
	// 数量不够
	if workList == nil || len(workList) < (offset+limit) {
		return nil, nil
	}
	// 获取全部
	if offset < 0 && limit < 0 {
		return workList, nil
	}
	// 正常获取，数量够
	returnList := workList[offset : offset+limit]
	return returnList, nil
}
