package services

import (
	"errors"
	"libs/utils"
	"models/entity"
	"models/mysql"
	"time"
)

// AddMenses2
func AddMenses2(uid, cid int64, m *entity.Menses) (*entity.Menses2, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if m == nil {
		return nil, errors.New("nil_menses")
	} else if m.Year <= 0 || m.MonthOfYear <= 0 || m.DayOfMonth <= 0 {
		return nil, errors.New("limit_happen_err")
	}
	dateAt := time.Date(m.Year, time.Month(m.MonthOfYear), m.DayOfMonth, 0, 0, 0, 0, time.Local)
	timeAt := utils.GetUnixByCSTDate(dateAt)
	// old + limit
	var old *entity.Menses2
	var err error
	if m.IsStart {
		// 开始
		old, err = mysql.GetMenses2ByUserCoupleAtBetween(uid, cid, timeAt)
		if err != nil {
			return nil, err
		} else if old == nil || old.Id <= 0 {
			// total
			monthTotal := mysql.GetMenses2TotalByUserCoupleDateStart(uid, cid, m.Year, m.MonthOfYear)
			if int(monthTotal) >= GetLimit().MensesMaxPerMonth {
				// 超出每月最大数
				return nil, errors.New("limit_total_over")
			}
			// repeat
			repeat, err := mysql.GetMenses2ByUserCoupleDateStart(uid, cid, m.Year, m.MonthOfYear, m.DayOfMonth)
			if err != nil {
				return nil, err
			} else if repeat != nil && repeat.Id > 0 {
				return nil, errors.New("menses_add_repeat")
			}
			// 执行 add/update
		} else {
			// 执行 update
		}
	} else {
		// 结束
		old, err = mysql.GetMenses2ByUserCoupleStartDown(uid, cid, timeAt+16*60*60) // 16*60*60，兼容 android v20
		if err != nil {
			return nil, err
		} else if old == nil || old.Id <= 0 {
			// 没有start的数据
			return nil, errors.New("limit_happen_err")
		} else {
			// durationDay
			var day = (timeAt - old.StartAt) / (60 * 60 * 24)
			if day < 0 || int(day) > GetLimit().MensesMaxDurationDay {
				return nil, errors.New("limit_happen_err")
			}
			// 执行 update/del
		}
	}
	// mysql
	if m.IsStart {
		// 开始
		if old == nil || old.Id <= 0 {
			// length，要-1
			mensesLength, err := GetMensesLengthByUserCouple(uid, cid)
			endDate := dateAt.Add(time.Hour * 24 * time.Duration(mensesLength.DurationDay-1))
			endAt := utils.GetUnixByCSTDate(endDate) + 24*60*60 - 1
			if err != nil {
				return nil, err
			} else if mensesLength == nil {
				return nil, errors.New("nil_menses_length")
			}
			// 查找当日所有的数据，防止重复添加
			old, err = mysql.GetMenses2AllByUserCoupleDateStart(uid, cid, m.Year, m.MonthOfYear, m.DayOfMonth)
			if err != nil {
				return nil, err
			}
			if old == nil || old.Id <= 0 {
				// add 重置start和end
				old = &entity.Menses2{
					BaseObj: entity.BaseObj{
						Status: entity.STATUS_VISIBLE,
					},
					BaseCp: entity.BaseCp{
						UserId:   uid,
						CoupleId: cid,
					},
					StartAt:          timeAt,
					EndAt:            endAt,
					StartYear:        m.Year,
					StartMonthOfYear: m.MonthOfYear,
					StartDayOfMonth:  m.DayOfMonth,
					EndYear:          endDate.Year(),
					EndMonthOfYear:   int(endDate.Month()),
					EndDayOfMonth:    endDate.Day(),
				}
				old, err = mysql.AddMenses2(old)
			} else {
				// update 重置start和end
				old.Status = entity.STATUS_VISIBLE
				old.StartAt = timeAt
				old.EndAt = endAt
				old.StartYear = m.Year
				old.StartMonthOfYear = m.MonthOfYear
				old.StartDayOfMonth = m.DayOfMonth
				old.EndYear = endDate.Year()
				old.EndMonthOfYear = int(endDate.Month())
				old.EndDayOfMonth = endDate.Day()
				old, err = mysql.UpdateMenses2(old)
			}
		} else {
			// update 重置start，不变end
			old.Status = entity.STATUS_VISIBLE
			old.StartAt = timeAt
			old.StartYear = m.Year
			old.StartMonthOfYear = m.MonthOfYear
			old.StartDayOfMonth = m.DayOfMonth
			old, err = mysql.UpdateMenses2(old)
		}
	} else {
		// 结束
		if old == nil || old.Id <= 0 {
			// 没有start的数据
			return nil, errors.New("limit_happen_err")
		} else {
			if old.StartYear != m.Year || old.StartMonthOfYear != m.MonthOfYear || old.StartDayOfMonth != m.DayOfMonth {
				// update 不变start，重置end
				endAt := time.Unix(timeAt-60*60*24, 0).Unix()
				endDate := utils.GetCSTDateByUnix(endAt)
				old.Status = entity.STATUS_VISIBLE
				old.EndAt = endAt + 24*60*60 - 1
				old.EndYear = endDate.Year()
				old.EndMonthOfYear = int(endDate.Month())
				old.EndDayOfMonth = endDate.Day()
				old, err = mysql.UpdateMenses2(old)
			} else {
				// del
				err = mysql.DelMenses2(old)
			}
		}
	}
	if old == nil || err != nil {
		return nil, err
	}
	// 同步
	go func() {
		// trends
		trends := CreateTrends(uid, cid, entity.TRENDS_ACT_TYPE_INSERT, entity.TRENDS_CON_TYPE_MENSES, old.Id)
		_, _ = AddTrends(trends)
		// push
		content := "push_content_menses_come"
		if !m.IsStart {
			content = "push_content_menses_gone"
		}
		AddPushInCouple(uid, m.Id, "push_title_note_update", content, entity.PUSH_TYPE_NOTE_MENSES)
	}()
	return old, err
}

// GetMenses2ListByUserCoupleYearMonth
func GetMenses2ListByUserCoupleYearMonth(uid, cid int64, year, month int) ([]*entity.Menses2, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if year <= 0 {
		return nil, errors.New("limit_year_nil")
	}
	// mysql
	limit := GetPageSizeLimit().Menses
	list, err := mysql.GetMenses2ListByUserCoupleYearMonth(uid, cid, year, month, 0, limit)
	if err != nil {
		return nil, err
	}
	if list != nil && len(list) > 0 {
		// 有记录，先添加真实记录的dayList
		for _, v := range list {
			if v != nil && v.Id > 0 {
				v.MensesDayList, _ = mysql.GetMensesDayListByMenses(v.Id, 0, GetLimit().MensesMaxDurationDay)
			}
		}
		// 只有一个，预测当月记录
		length, _ := GetMensesLengthByUserCouple(uid, cid)
		if length != nil {
			if len(list) == 1 {
				menses2 := list[0]
				// before
				startYear := year
				startMonth := month
				monthStartDate := time.Date(startYear, time.Month(startMonth), 0, 0, 0, 0, 0, time.Local)
				monthStartUnix := utils.GetUnixByCSTDate(monthStartDate)
				menses2Before, _ := mysql.GetMenses2ByUserCoupleEndNear(uid, cid, monthStartUnix)
				if menses2Before == nil || menses2Before.Id <= 0 || int((menses2.StartAt-menses2Before.EndAt)/(60*60*24)) > GetLimit().MensesMaxCycleDay {
					// 两者相距天数不少，可以预测
					beforeEndUnix := menses2.StartAt - int64(length.CycleDay*60*60*24) + int64((length.DurationDay-1)*60*60*24)
					beforeEndDate := utils.GetCSTDateByUnix(beforeEndUnix)
					if int(beforeEndDate.Month()) == month {
						// beforeEndDate
						endAt := beforeEndUnix
						startAt := endAt - int64((length.DurationDay-1)*60*60*24)
						list = append(list, newMenses2NoReal(uid, cid, startAt, endAt))
					}
				}
				// after
				var endYear int
				var endMonth int
				if month < 12 {
					endYear = year
					endMonth = month + 1
				} else {
					endYear = year + 1
					endMonth = 1
				}
				monthEndDate := time.Date(endYear, time.Month(endMonth), 0, 0, 0, 0, 0, time.Local)
				monthEndUnix := utils.GetUnixByCSTDate(monthEndDate)
				menses2After, _ := mysql.GetMenses2ByUserCoupleStartNear(uid, cid, monthEndUnix)
				if menses2After == nil || menses2After.Id <= 0 || int((menses2After.StartAt-menses2.EndAt)/(60*60*24)) > GetLimit().MensesMaxCycleDay {
					// 两者相距天数不少，可以预测
					afterStartUnix := menses2.StartAt + int64(length.CycleDay*60*60*24)
					afterStartDate := utils.GetCSTDateByUnix(afterStartUnix)
					if int(afterStartDate.Month()) == month {
						// afterStartDate
						startAt := afterStartUnix
						endAt := startAt + int64((length.DurationDay-1)*60*60*24)
						list = append(list, newMenses2NoReal(uid, cid, startAt, endAt))
					}
				}
			}
		}
	} else {
		// 无记录，预测
		startYear := year
		startMonth := month
		monthStartDate := time.Date(startYear, time.Month(startMonth), 0, 0, 0, 0, 0, time.Local)
		monthStartUnix := utils.GetUnixByCSTDate(monthStartDate)
		menses2Before, err := mysql.GetMenses2ByUserCoupleEndNear(uid, cid, monthStartUnix)
		if err != nil {
			// 之前的记录出错
			return nil, err
		}
		var endYear int
		var endMonth int
		if month < 12 {
			endYear = year
			endMonth = month + 1
		} else {
			endYear = year + 1
			endMonth = 1
		}
		monthEndDate := time.Date(endYear, time.Month(endMonth), 0, 0, 0, 0, 0, time.Local)
		monthEndUnix := utils.GetUnixByCSTDate(monthEndDate)
		menses2After, err := mysql.GetMenses2ByUserCoupleStartNear(uid, cid, monthEndUnix)
		if err != nil {
			// 之后的记录出错
			return nil, err
		}
		if (menses2Before == nil || menses2Before.Id <= 0) && (menses2After == nil || menses2After.Id <= 0) {
			// 没有历史记录，无法预测
			return nil, nil
		}
		length, _ := GetMensesLengthByUserCouple(uid, cid)
		if length == nil {
			// 设置出错
			return nil, nil
		}
		list = make([]*entity.Menses2, 0)
		if (menses2Before != nil && menses2Before.Id > 0) && (menses2After == nil || menses2After.Id <= 0) {
			// menses2Before
			for timeAt := menses2Before.StartAt; timeAt <= monthEndUnix; timeAt += int64(length.CycleDay * 60 * 60 * 24) {
				startAt := timeAt
				endAt := startAt + int64((length.DurationDay-1)*60*60*24)
				if endAt < monthStartUnix {
					// 结束日期不在当月
					continue
				}
				list = append(list, newMenses2NoReal(uid, cid, startAt, endAt))
			}
		} else if (menses2Before == nil || menses2Before.Id <= 0) && (menses2After != nil && menses2After.Id > 0) {
			// menses2After
			for timeAt := menses2After.StartAt; timeAt+int64((length.DurationDay-1)*60*60*24) >= monthStartUnix; timeAt -= int64(length.CycleDay * 60 * 60 * 24) {
				startAt := timeAt
				endAt := startAt + int64((length.DurationDay-1)*60*60*24)
				if startAt > monthEndUnix {
					// 开始日期不在当月
					continue
				}
				list = append(list, newMenses2NoReal(uid, cid, startAt, endAt))
			}
		} else {
			if int((menses2After.StartAt-menses2Before.EndAt)/(60*60*24)) <= GetLimit().MensesMaxCycleDay {
				// 两者相距天数太少，不预测
				return nil, nil
			}
			// 计算两个哪个更接近
			if monthStartUnix-menses2Before.EndAt < menses2After.StartAt-monthEndUnix {
				// menses2Before
				for timeAt := menses2Before.StartAt; timeAt <= monthEndUnix; timeAt += int64(length.CycleDay * 60 * 60 * 24) {
					startAt := timeAt
					endAt := startAt + int64((length.DurationDay-1)*60*60*24)
					if endAt < monthStartUnix {
						// 结束日期不在当月
						continue
					}
					list = append(list, newMenses2NoReal(uid, cid, startAt, endAt))
				}
			} else {
				// menses2After
				for timeAt := menses2After.StartAt; timeAt+int64((length.DurationDay-1)*60*60*24) >= monthStartUnix; timeAt -= int64(length.CycleDay * 60 * 60 * 24) {
					startAt := timeAt
					endAt := startAt + int64((length.DurationDay-1)*60*60*24)
					if startAt > monthEndUnix {
						// 开始日期不在当月
						continue
					}
					list = append(list, newMenses2NoReal(uid, cid, startAt, endAt))
				}
			}
		}
	}
	if list != nil && len(list) > 0 {
		// 计算安全期，危险期，排卵日
	} else {
		// return nil, errors.New("no_data_menses") 不要了，防止toast
	}
	// 同步
	go func() {
		trends := CreateTrendsByList(uid, cid, entity.TRENDS_ACT_TYPE_QUERY, entity.TRENDS_CON_TYPE_MENSES)
		AddTrends(trends)
	}()
	return list, err
}

func newMenses2NoReal(uid, cid, startAt, endAt int64) *entity.Menses2 {
	startDate := utils.GetCSTDateByUnix(startAt)
	endDate := utils.GetCSTDateByUnix(endAt)
	realMenses2, _ := mysql.GetMenses2ByUserCoupleDateStart(uid, cid, startDate.Year(), int(startDate.Month()), startDate.Day())
	if realMenses2 != nil && realMenses2.Id > 0 {
		// 防止预测到真实的数据，导致预测数据和真实数据不准
		return realMenses2
	}
	return &entity.Menses2{
		BaseCp: entity.BaseCp{
			UserId:   uid,
			CoupleId: cid,
		},
		StartAt:          startAt,
		EndAt:            endAt,
		StartYear:        startDate.Year(),
		StartMonthOfYear: int(startDate.Month()),
		StartDayOfMonth:  startDate.Day(),
		EndYear:          endDate.Year(),
		EndMonthOfYear:   int(endDate.Month()),
		EndDayOfMonth:    endDate.Day(),
		IsReal:           false,
	}
}
