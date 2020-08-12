package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"models/entity"
	"services"
)

// HandlerUser
func HandlerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		checkModelStatus(w, r, "user_register")
		PostUser(w, r)
	} else if r.Method == http.MethodPut {
		checkModelStatus(w, r, "user_modify")
		PutUser(w, r)
	} else if r.Method == http.MethodGet {
		GetUser(w, r)
	} else {
		response405(w, r)
	}
}

// PostUser
func PostUser(w http.ResponseWriter, r *http.Request) {
	// 参数
	values := r.URL.Query()
	code := values.Get("code")
	postUser := &entity.User{}
	checkRequestBody(w, r, postUser)
	// sms检查
	err := services.SmsCheckCode(postUser.Phone, entity.SMS_TYPE_REGISTER, code)
	response417ErrDialog(w, r, err)
	// 开始注册
	addUser, err := services.AddUserByPwd(postUser)
	response417ErrDialog(w, r, err)
	if addUser != nil {
		addUser.Password = ""
	}
	// 返回
	response417NoUserInfo(w, r, addUser)
}

// PutUser
func PutUser(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	values := r.URL.Query()
	paramsType, _ := strconv.Atoi(values.Get("type"))
	paramsUser := &entity.User{}
	checkRequestBody(w, r, paramsUser)
	// 分类修改
	var modifyUser *entity.User
	var err error
	switch paramsType {
	case services.USER_TYPE_FORGET_PWD: // 忘记密码
		// 登录之前做的操作
		modifyUser, err = services.GetUserByPhone(paramsUser.Phone)
		response417ErrDialog(w, r, err)
		if modifyUser == nil {
			response417Dialog(w, r, "user_phone_no_exist")
		} else if modifyUser.Status <= entity.STATUS_DELETE {
			response406(w, r)
		}
		code := values.Get("code")
		err = services.SmsCheckCode(paramsUser.Phone, entity.SMS_TYPE_FORGET, code)
		response417ErrDialog(w, r, err)
		// 开始修改
		modifyUser.Password = paramsUser.Password
		modifyUser, err = services.UpdateUserPwd(modifyUser)
		response417ErrDialog(w, r, err)
	case services.USER_TYPE_UPDATE_PWD: // 修改密码
		// 登录之后做的操作
		modifyUser = checkTokenUser(w, r)
		oldPwd := values.Get("old_pwd")
		if len([]rune(modifyUser.Password)) <= 0 {
			// 没有旧密码，验证码注册没有旧密码
			modifyUser.Password = paramsUser.Password
			modifyUser, err = services.UpdateUserPwd(modifyUser)
			response417ErrDialog(w, r, err)
		} else {
			// 有旧密码，则校验
			if modifyUser.Password != oldPwd { // 有旧密码且错误
				response417Dialog(w, r, "user_pwd_old_wrong")
			}
			if modifyUser.Password == paramsUser.Password { // 密码相同
				response417Dialog(w, r, "user_pwd_same")
			}
			// 开始修改
			modifyUser.Password = paramsUser.Password
			modifyUser, err = services.UpdateUserPwd(modifyUser)
			response417ErrDialog(w, r, err)
		}
	case services.USER_TYPE_UPDATE_PHONE: // 修改手机
		// 登录之后做的操作
		modifyUser = checkTokenUser(w, r)
		code := values.Get("code")
		// phone检验
		if modifyUser.Phone == paramsUser.Phone {
			response417Dialog(w, r, "user_phone_same")
		}
		// sms检查
		err = services.SmsCheckCode(paramsUser.Phone, entity.SMS_TYPE_PHONE, code)
		response417ErrDialog(w, r, err)
		// 开始修改
		oldPhone := modifyUser.Phone
		modifyUser.Phone = paramsUser.Phone
		modifyUser, err = services.UpdateUserPhone(modifyUser, oldPhone)
		response417ErrDialog(w, r, err)
	case services.USER_TYPE_UPDATE_INFO: // 修改信息
		// 登录之后做的操作，但是不能检查info
		token := strings.TrimSpace(r.Header.Get(HEAD_PARAMS_TOKEN1))
		if len(token) <= 0 {
			token = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_TOKEN2))
		}
		if len(token) <= 0 {
			token = strings.TrimSpace(r.Header.Get(HEAD_PARAMS_TOKEN3))
		}
		modifyUser, err = services.GetUserByToken(token)
		response417ErrDialog(w, r, err)
		if modifyUser == nil {
			response417Dialog(w, r, "nil_user")
		} else if modifyUser.Status <= entity.STATUS_DELETE {
			response406(w, r)
		} else if services.IsUserInfoComplete(modifyUser) {
			// 只能修改一次
			response417Dialog(w, r, "user_info_just_one")
		}
		modifyUser.Birthday = paramsUser.Birthday
		modifyUser.Sex = paramsUser.Sex
		modifyUser, err = services.UpdateUserInfo(modifyUser)
		response417ErrDialog(w, r, err)
	case services.USER_TYPE_ADMIN_UPDATE_INFO: // 修改信息(管理员)
		// 登录之后做的操作
		user := checkTokenUser(w, r)
		// admin
		if !services.IsAdminister(user) {
			response417Toast(w, r, "")
		}
		modifyUser, _ = services.GetUserById(paramsUser.Id)
		if modifyUser == nil {
			response417Dialog(w, r, "nil_user")
		}
		modifyUser.Birthday = paramsUser.Birthday
		modifyUser.Sex = paramsUser.Sex
		modifyUser, err = services.UpdateUserInfo(modifyUser)
		response417ErrDialog(w, r, err)
	case services.USER_TYPE_ADMIN_UPDATE_STATUS: // 修改状态(管理员)
		// 登录之后做的操作
		user := checkTokenUser(w, r)
		// admin
		if !services.IsAdminister(user) {
			response417Toast(w, r, "")
		}
		modifyUser, _ = services.GetUserById(paramsUser.Id)
		if modifyUser == nil {
			response417Dialog(w, r, "nil_user")
		}
		modifyUser, err = services.ToggleUserStatus(modifyUser)
		response417ErrDialog(w, r, err)
	default:
		response405(w, r)
	}
	// 装载信息
	if modifyUser != nil {
		modifyUser.Couple, _ = services.GetCoupleVisibleByUser(modifyUser.Id)
		modifyUser.Password = ""
	}
	response200DataToast(w, r, "db_update_success", struct {
		User *entity.User `json:"user"`
	}{modifyUser})
}

// GetCouple
func GetUser(w http.ResponseWriter, r *http.Request) {
	// 接受参数
	values := r.URL.Query()
	ta, _ := strconv.ParseBool(values.Get("ta"))
	uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
	phone := values.Get("phone")
	list, _ := strconv.ParseBool(values.Get("list"))
	black, _ := strconv.ParseBool(values.Get("black"))
	total, _ := strconv.ParseBool(values.Get("total"))
	birth, _ := strconv.ParseBool(values.Get("birth"))
	if ta {
		me := checkTokenCouple(w, r)
		// ta
		taId := services.GetTaId(me)
		user, err := services.GetUserById(taId)
		response417ErrDialog(w, r, err)
		if user == nil || user.Status <= entity.STATUS_DELETE {
			response417Dialog(w, r, "nil_user")
		} else {
			user.Password = ""
			user.UserToken = ""
		}
		// 返回
		response200Data(w, r, struct {
			User *entity.User `json:"user"`
		}{user})
	} else if uid > 0 {
		me := checkTokenUser(w, r)
		// admin检查
		if !services.IsAdminister(me) {
			response200Toast(w, r, "")
		}
		user, err := services.GetUserById(uid)
		response417ErrDialog(w, r, err)
		if user != nil {
			user.Password = ""
			user.UserToken = ""
		}
		// 返回
		response200Data(w, r, struct {
			User *entity.User `json:"user"`
		}{user})
	} else if len(phone) > 0 {
		me := checkTokenUser(w, r)
		// admin检查
		if !services.IsAdminister(me) {
			response200Toast(w, r, "")
		}
		user, err := services.GetUserByPhone(phone)
		response417ErrDialog(w, r, err)
		if user != nil {
			user.Password = ""
			user.UserToken = ""
		}
		// 返回
		response200Data(w, r, struct {
			User *entity.User `json:"user"`
		}{user})
	} else if list {
		me := checkTokenUser(w, r)
		// admin检查
		if !services.IsAdminister(me) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		userList, err := services.GetUserList(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			UserList []*entity.User `json:"userList"`
		}{userList})
	} else if black {
		me := checkTokenUser(w, r)
		// admin检查
		if !services.IsAdminister(me) {
			response200Toast(w, r, "")
		}
		page, _ := strconv.Atoi(values.Get("page"))
		userList, err := services.GetUserListByBlack(page)
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			UserList []*entity.User `json:"userList"`
		}{userList})
	} else if total {
		me := checkTokenUser(w, r)
		// admin检查
		if !services.IsAdminister(me) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		total := services.GetUserTotalByCreate(start, end)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else if birth {
		me := checkTokenUser(w, r)
		// admin检查
		if !services.IsAdminister(me) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		birth := services.GetUserBirthAvgByCreate(start, end)
		// 返回
		response200Data(w, r, struct {
			Birth float64 `json:"birth"`
		}{birth})
	} else {
		response405(w, r)
	}
}
