package controllers

import (
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerCoin(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_coin")
	if r.Method == http.MethodPost {
		PostCoin(w, r)
	} else if r.Method == http.MethodGet {
		GetCoin(w, r)
	} else {
		response405(w, r)
	}
}

// PostCoin
func PostCoin(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	coin := &entity.Coin{}
	checkRequestBody(w, r, coin)
	// 开始插入
	var err error
	if coin.Kind == entity.COIN_KIND_ADD_BY_AD_WATCH || coin.Kind == entity.COIN_KIND_ADD_BY_AD_CLICK {
		coin, err = services.AddCoinByAd(user.Id, couple.Id, coin.Kind)
	} else if coin.Kind == entity.COIN_KIND_ADD_BY_SYS {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		coin, err = services.AddCoinByAdmin(coin)
	}
	response417ErrToast(w, r, err)
	// 返回
	response200DataToast(w, r, "coin_in_account", struct {
		Coin *entity.Coin `json:"coin"`
	}{coin})
}

// GetCoin
func GetCoin(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	home, _ := strconv.ParseBool(values.Get("home"))
	list, _ := strconv.ParseBool(values.Get("list"))
	change, _ := strconv.ParseBool(values.Get("change"))
	total, _ := strconv.ParseBool(values.Get("total"))
	if home {
		coin, err := services.GetCoinLatest(couple.Id)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			Coin *entity.Coin `json:"coin"`
		}{coin})
	} else if list {
		page, _ := strconv.Atoi(values.Get("page"))
		admin, _ := strconv.ParseBool(values.Get("admin"))
		var coinList []*entity.Coin
		var err error
		if admin && services.IsAdminister(user) {
			uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
			cid, _ := strconv.ParseInt(values.Get("cid"), 10, 64)
			bid, _ := strconv.ParseInt(values.Get("bid"), 10, 64)
			kind, _ := strconv.Atoi(values.Get("kind"))
			coinList, err = services.GetCoinList(uid, cid, bid, kind, page)
		} else {
			coinList, err = services.GetCoinListByCouple(couple.Id, page)
		}
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			CoinList []*entity.Coin `json:"coinList"`
		}{coinList})
	} else if change {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		infoList, err := services.GetCoinChangeListByCreateWithPay(start, end)
		response417ErrDialog(w, r, err)
		// 返回
		response200Data(w, r, struct {
			InfoList []*entity.FiledInfo `json:"infoList"`
		}{infoList})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		kind, _ := strconv.Atoi(values.Get("kind"))
		total := services.GetCoinTotalByCreateKindWithDel(start, end, kind)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
