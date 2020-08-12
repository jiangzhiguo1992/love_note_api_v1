package controllers

import (
	"encoding/json"
	"models/entity"
	"net/http"
	"services"
	"strconv"
)

func HandlerBill(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "more_bill")
	if r.Method == http.MethodPost {
		PostBill(w, r)
	} else if r.Method == http.MethodGet {
		GetBill(w, r)
	} else {
		response405(w, r)
	}
}

// PostBill
func PostBill(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	check, _ := strconv.ParseBool(values.Get("check"))
	if check {
		// 接受参数(这块order可能为nil)
		bid, _ := strconv.ParseInt(values.Get("bid"), 10, 64)
		order := &services.AppleOrder{}
		bytes := getRequestBody(r)
		json.Unmarshal(bytes, order)
		var err error
		if bid > 0 {
			err = services.CheckBillStatus(bid)
		} else {
			err = services.CheckBillsStatus(user.Id, couple.Id, order)
		}
		response417ErrDialog(w, r, err)
		// 返回
		response200Toast(w, r, "db_update_success")
	} else {
		// 回调url最好不到带参数，有可能回调失败
		// 这里最好不要做，统一自己回调
		response200Data(w, r, nil)
	}
}

// GetBill
func GetBill(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	values := r.URL.Query()
	before, _ := strconv.ParseBool(values.Get("before"))
	list, _ := strconv.ParseBool(values.Get("list"))
	amount, _ := strconv.ParseBool(values.Get("amount"))
	total, _ := strconv.ParseBool(values.Get("total"))
	if before {
		platformPay, _ := strconv.Atoi(values.Get("pay_platform"))
		goods, _ := strconv.Atoi(values.Get("goods"))
		if platformPay == entity.BILL_PLATFORM_PAY_ALI {
			// 检查支付宝
			checkModelStatus(w, r, "more_pay_ali")
			// orderBefore
			orderBefore := &services.OrderBefore{Platform: platformPay}
			var err error
			orderBefore.AliOrder, err = services.GetAliPayOrderInfo(user.Id, couple.Id, goods)
			response417ErrDialog(w, r, err)
			// 返回
			response200Data(w, r, struct {
				OrderBefore *services.OrderBefore `json:"orderBefore"`
			}{orderBefore})
		} else if platformPay == entity.BILL_PLATFORM_PAY_WX {
			// 检查微信
			checkModelStatus(w, r, "more_pay_wx")
			// orderBefore
			orderBefore := &services.OrderBefore{Platform: platformPay}
			var err error
			orderBefore.WXOrder, err = services.GetWXPayOrderInfo(user.Id, couple.Id, goods)
			response417ErrDialog(w, r, err)
			// 返回
			response200Data(w, r, struct {
				OrderBefore *services.OrderBefore `json:"orderBefore"`
			}{orderBefore})
		} else if platformPay == entity.BILL_PLATFORM_PAY_APPLE {
			// 检查苹果
			checkModelStatus(w, r, "more_pay_apple")
			// orderBefore
			orderBefore := &services.OrderBefore{Platform: platformPay}
			var err error
			orderBefore.AppleOrder, err = services.GetApplePayOrderInfo(user.Id, couple.Id, goods)
			response417ErrDialog(w, r, err)
			// 返回
			response200Data(w, r, struct {
				OrderBefore *services.OrderBefore `json:"orderBefore"`
			}{orderBefore})
		} else if platformPay == entity.BILL_PLATFORM_PAY_GOOGLE {
			// google
		} else {
			response405(w, r)
		}
	} else if list {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		sync, _ := strconv.ParseBool(values.Get("sync"))
		page, _ := strconv.Atoi(values.Get("page"))
		var billList []*entity.Bill
		var err error
		if sync {
			billList, err = services.GetBillListWithNoSync(page)
		} else {
			uid, _ := strconv.ParseInt(values.Get("uid"), 10, 64)
			cid, _ := strconv.ParseInt(values.Get("cid"), 10, 64)
			tradeNo := values.Get("trade_no")
			billList, err = services.GetBillList(uid, cid, tradeNo, page)
		}
		response200ErrShow(w, r, err)
		// 返回
		response200Data(w, r, struct {
			BillList []*entity.Bill `json:"billList"`
		}{billList})
	} else if amount {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		platformOs := values.Get("platform_os")
		platformPay, _ := strconv.Atoi(values.Get("platform_pay"))
		payType, _ := strconv.Atoi(values.Get("pay_type"))
		goodsType, _ := strconv.Atoi(values.Get("goods_type"))
		amount := services.GetBillAmountTotalByCreateWithPay(start, end, platformOs, platformPay, payType, goodsType)
		// 返回
		response200Data(w, r, struct {
			Amount float64 `json:"amount"`
		}{amount})
	} else if total {
		// admin检查
		if !services.IsAdminister(user) {
			response200Toast(w, r, "")
		}
		start, _ := strconv.ParseInt(values.Get("start"), 10, 64)
		end, _ := strconv.ParseInt(values.Get("end"), 10, 64)
		platformOs := values.Get("platform_os")
		platformPay, _ := strconv.Atoi(values.Get("platform_pay"))
		payType, _ := strconv.Atoi(values.Get("pay_type"))
		goodsType, _ := strconv.Atoi(values.Get("goods_type"))
		total := services.GetBillTotalByCreateWithDelApy(start, end, platformOs, platformPay, payType, goodsType)
		// 返回
		response200Data(w, r, struct {
			Total int64 `json:"total"`
		}{total})
	} else {
		response405(w, r)
	}
}
