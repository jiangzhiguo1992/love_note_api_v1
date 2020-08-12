package controllers

import (
	"net/http"
	"strconv"
	"time"

	"libs/utils"
	"models/entity"
	"models/mysql"
	"services"
)

// HandlerCouple 在router中没有进行couple检查  在当下做检查
func HandlerCouple(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		checkModelStatus(w, r, "couple_pair")
		PostCouple(w, r)
	} else if r.Method == http.MethodPut {
		checkModelStatus(w, r, "couple_modify")
		PutCouple(w, r)
	} else if r.Method == http.MethodGet {
		GetCouple(w, r)
	} else {
		response405(w, r)
	}
}

// PostCouple
func PostCouple(w http.ResponseWriter, r *http.Request) {
	me := getTokenUser(r)
	// 接受参数
	ta := &entity.User{}
	checkRequestBody(w, r, ta)
	// 检查ta
	ta, _ = services.GetUserByPhone(ta.Phone)
	if ta == nil {
		response417Dialog(w, r, "user_phone_no_exist")
	} else if ta.Status <= entity.STATUS_DELETE {
		response417Dialog(w, r, "nil_user")
	}
	// 获取双方曾经的cp(两个人的关系检查)
	couple, err := services.GetCoupleBy2User(me.Id, ta.Id)
	response417ErrDialog(w, r, err)
	if couple == nil {
		// 曾经没有邀请
		couple = &entity.Couple{}
		couple.CreatorId = me.Id
		couple.InviteeId = ta.Id
		couple, err = services.AddCouple(couple)
		response417ErrDialog(w, r, err)
	} else {
		// 曾经有过邀请
		couple, err = services.AddCoupleStateInvitee(me.Id, couple)
		response417ErrDialog(w, r, err)
	}
	// 返回
	response200DataToast(w, r, "couple_invite_send", struct {
		Couple *entity.Couple `json:"couple"`
	}{couple})
}

// PutCouple
func PutCouple(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 参数
	values := r.URL.Query()
	paramsType, _ := strconv.Atoi(values.Get("type"))
	paramsCouple := &entity.Couple{}
	checkRequestBody(w, r, paramsCouple)
	// 按类型处理
	var toast string
	if paramsType == services.COUPLE_UPDATE_GOOD || paramsType == services.COUPLE_UPDATE_BAD {
		// 修改状态
		user.Couple, _ = services.GetCoupleById(paramsCouple.Id)
		couple := user.Couple
		if couple == nil || couple.Id <= 0 {
			response417Dialog(w, r, "nil_couple")
		} else if !services.IsUserInCouple(user.Id, couple) {
			// 不是自己的cp
			response417Dialog(w, r, "db_update_refuse")
		}
		// 检查自己已存在的couple，防止多个可操作的couple生成
		self, err := services.GetCoupleSelfByUser(user.Id)
		response417ErrDialog(w, r, err)
		if self == nil {
			response417Toast(w, r, "")
		} else if self.Id != couple.Id {
			response417Dialog(w, r, "couple_state_error")
		}
		// 开始变更状态
		state := couple.State
		state.CoupleId = couple.Id
		if paramsType == services.COUPLE_UPDATE_GOOD {
			// --->更好的状态
			if state.State == entity.COUPLE_STATE_INVITE {
				// 邀请状态
				if state.UserId == user.Id {
					// 邀请人同意
					response417Dialog(w, r, "couple_status_together_limit")
				} else {
					// 被邀请人同意
					state.UserId = user.Id
					state.State = entity.COUPLE_STATE_520
					state, err = services.AddCoupleState(state)
					response417ErrDialog(w, r, err)
					couple.State = state
					toast = "couple_status_together"
				}
			} else if state.State == entity.COUPLE_STATE_BREAK || state.State == entity.COUPLE_STATE_BREAK_ACCEPT {
				// 分手状态
				breakTime := services.CoupleNoBreakTime()
				if state.CreateAt < breakTime || state.State == entity.COUPLE_STATE_BREAK_ACCEPT {
					// 已分手
					response417Dialog(w, r, "couple_status_break_always")
				} else {
					// 正在分手
					if state.UserId == user.Id {
						// 分手者取消
						state.UserId = user.Id
						state.State = entity.COUPLE_STATE_520
						state, err = services.AddCoupleState(state)
						response417ErrDialog(w, r, err)
						couple.State = state
						toast = "couple_status_cancel_break"
					} else {
						// 被分手者取消
						response417Dialog(w, r, "couple_status_cancel_break_limit")
					}
				}
			} else if state.State == entity.COUPLE_STATE_520 {
				// 还在一起
				response417Dialog(w, r, "couple_status_together_always")
			} else {
				response417Dialog(w, r, "couple_state_error")
			}
		} else {
			// --->更坏的状态
			if state.State == entity.COUPLE_STATE_INVITE {
				// 邀请状态
				if state.UserId == user.Id {
					// 邀请人取消
					state.State = entity.COUPLE_STATE_INVITE_CANCEL
					state, err = services.AddCoupleState(state)
					response417ErrDialog(w, r, err)
					couple.State = state
					toast = ""
				} else {
					// 被邀请人拒绝
					state.UserId = user.Id
					state.State = entity.COUPLE_STATE_INVITE_REJECT
					state, err = services.AddCoupleState(state)
					response417ErrDialog(w, r, err)
					couple.State = state
					toast = ""
				}
			} else if state.State == entity.COUPLE_STATE_BREAK || state.State == entity.COUPLE_STATE_BREAK_ACCEPT {
				// 分手状态
				breakTime := services.CoupleNoBreakTime()
				if state.CreateAt < breakTime || state.State == entity.COUPLE_STATE_BREAK_ACCEPT {
					// 已分手
					response417Dialog(w, r, "couple_status_break_always")
				} else {
					// 正在分手
					if state.UserId == user.Id {
						// 分手者同意
						response417Dialog(w, r, "couple_status_break_repeat")
					} else {
						// 被分手者同意
						state.UserId = user.Id
						state.State = entity.COUPLE_STATE_BREAK_ACCEPT
						state, err = services.AddCoupleState(state)
						response417ErrDialog(w, r, err)
						couple.State = state
						toast = "couple_status_break"
					}
				}
			} else if state.State == entity.COUPLE_STATE_520 {
				// 还在一起
				breakNeedSec := services.GetLimit().CoupleBreakNeedSec
				if state.CreateAt+int64(breakNeedSec) > time.Now().Unix() {
					// 在一起时间太短，直接分手
					state.UserId = user.Id
					state.State = entity.COUPLE_STATE_BREAK_ACCEPT
					state, err = services.AddCoupleState(state)
					response417ErrDialog(w, r, err)
					couple.State = state
					toast = "couple_status_break"
				} else {
					// 在一起一定天数后，分手才有倒计时
					state.UserId = user.Id
					state.State = entity.COUPLE_STATE_BREAK
					state, err = services.AddCoupleState(state)
					response417ErrDialog(w, r, err)
					couple.State = state
					toast = ""
				}
			} else {
				response417Dialog(w, r, "couple_state_error")
			}
		}
	} else if paramsType == services.COUPLE_UPDATE_INFO {
		// 修改信息
		user = checkTokenCouple(w, r)
		couple := user.Couple
		// 开始修改，无效信息过滤掉，只能修改对方的
		if paramsCouple.TogetherAt != 0 {
			couple.TogetherAt = paramsCouple.TogetherAt
		}
		if couple.CreatorId == user.Id {
			if len(paramsCouple.InviteeName) > 0 || paramsCouple.InviteeName != "" {
				couple.InviteeName = paramsCouple.InviteeName
			}
			if len(paramsCouple.InviteeAvatar) > 0 || paramsCouple.InviteeAvatar != "" {
				couple.InviteeAvatar = paramsCouple.InviteeAvatar
			}
		} else {
			if len(paramsCouple.CreatorName) > 0 || paramsCouple.CreatorName != "" {
				couple.CreatorName = paramsCouple.CreatorName
			}
			if len(paramsCouple.CreatorAvatar) > 0 || paramsCouple.CreatorAvatar != "" {
				couple.CreatorAvatar = paramsCouple.CreatorAvatar
			}
		}
		var err error
		user.Couple, err = services.UpdateCouple(couple)
		response417ErrDialog(w, r, err)
		toast = "db_update_success"
	} else {
		response405(w, r)
	}
	// 返回
	response200DataToast(w, r, toast, struct {
		Couple *entity.Couple `json:"couple"`
	}{user.Couple})
}

// GetCouple
func GetCouple(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	values := r.URL.Query()
	self, _ := strconv.ParseBool(values.Get("self"))
	list, _ := strconv.ParseBool(values.Get("list"))
	total, _ := strconv.ParseBool(values.Get("total"))
	state, _ := strconv.ParseBool(values.Get("state"))
	uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
	cid, _ := strconv.ParseInt(values.Get("cid"), 10, 64)
	all, _ := strconv.ParseBool(values.Get("all"))
	if self {
		// 一般为配对界面获取self可见的couple
		couple, _ := services.GetCoupleSelfByUser(user.Id);
		if couple == nil {
			response200Data(w, r, nil)
		}
		// 配对数据封装
		pairCard := &services.PairCard{}
		// ta检查
		user.Couple = couple
		taId := services.GetTaId(user)
		if ta, _ := services.GetUserById(taId); ta != nil {
			pairCard.TaPhone = ta.Phone
		}
		// language
		language := "zh-cn"
		entry, err := mysql.GetEntryLatestByUser(user.Id)
		if err == nil && entry != nil {
			language = entry.Language
		}
		// 开始检查state
		state := couple.State
		if state.State == entity.COUPLE_STATE_INVITE {
			if state.UserId == user.Id {
				// 邀请者显示
				pairCard.Title = utils.GetLanguage(language, "couple_response_status_invite_send_title")
				pairCard.Desc = utils.GetLanguage(language, "couple_response_status_invite_send_message")
				pairCard.BtnBad = utils.GetLanguage(language, "couple_response_status_invite_send_bad")
				pairCard.BtnGood = ""
			} else {
				// 被邀请者显示
				pairCard.Title = utils.GetLanguage(language, "couple_response_status_invite_get_title")
				pairCard.Desc = utils.GetLanguage(language, "couple_response_status_invite_get_message")
				pairCard.BtnBad = utils.GetLanguage(language, "couple_response_status_invite_get_bad")
				pairCard.BtnGood = utils.GetLanguage(language, "couple_response_status_invite_get_good")
			}
		} else {
			response200Show(w, r, "couple_state_error")
		}
		// 返回
		response200Data(w, r, struct {
			PairCard interface{}    `json:"pairCard"`
			Couple   *entity.Couple `json:"couple"`
		}{
			pairCard,
			couple,
		})
	} else if list {
		page, _ := strconv.Atoi(values.Get("page"))
		if !state {
			// admin检查
			if !services.IsAdminister(user) {
				response200Toast(w, r, "")
			}
			coupleList, err := services.GetCoupleList(uid, page)
			response200ErrShow(w, r, err)
			// 返回
			response200Data(w, r, struct {
				CoupleList []*entity.Couple `json:"coupleList"`
			}{coupleList})
		} else {
			coupleStateList, err := services.GetCoupleStateListByCouple(cid, page)
			response200ErrShow(w, r, err)
			// 返回
			response200Data(w, r, struct {
				CoupleStateList []*entity.CoupleState `json:"coupleStateList"`
			}{coupleStateList})
		}
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		if !state {
			total := services.GetCoupleTotalByCreateWithDel(start, end)
			// 返回
			response200Data(w, r, struct {
				Total int64 `json:"total"`
			}{total})
		} else {
			infoList, err := services.GetCoupleStateStateListByCreate(start, end)
			response417ErrDialog(w, r, err)
			// 返回
			response200Data(w, r, struct {
				InfoList []*entity.FiledInfo `json:"infoList"`
			}{infoList})
		}
	} else if uid > 0 {
		// 一般用于获取可见cp
		couple, err := services.GetCoupleVisibleByUser(uid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Couple *entity.Couple `json:"couple"`
		}{couple})
	} else if uid > 0 && all {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		couple, err := services.GetCoupleSelfByUser(uid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Couple *entity.Couple `json:"couple"`
		}{couple})
	} else if cid > 0 {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		couple, err := services.GetCoupleById(cid)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Couple *entity.Couple `json:"couple"`
		}{couple})
	} else {
		response405(w, r)
	}
}
