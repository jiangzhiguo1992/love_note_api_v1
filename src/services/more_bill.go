package services

import (
	"errors"
	"libs/alipay"
	"libs/apple"
	"libs/utils"
	"libs/wxpay"
	"math"
	"models/entity"
	"models/mysql"
	"strconv"
	"strings"
	"time"
)

// GetAliPayOrderInfo 生成阿里订单，并同步bill
func GetAliPayOrderInfo(uid, cid int64, goodsType int) (string, error) {
	if uid <= 0 {
		return "", errors.New("nil_user")
	} else if cid <= 0 {
		return "", errors.New("nil_couple")
	}
	// goods
	goods := getGoodsByType(goodsType)
	if goods == nil {
		return "", errors.New("pay_goods_nil")
	}
	// order
	p, results, err := getAliPayWithTrade(uid, cid, goods)
	if err != nil {
		return "", err
	}
	// bill
	b := &entity.Bill{
		PlatformPay:  entity.BILL_PLATFORM_PAY_ALI,
		PayType:      entity.BILL_PAY_TYPE_APP,
		PayAmount:    goods.Amount,
		TradeNo:      p.OutTradeNo,
		TradeReceipt: "",
		GoodsType:    goodsType,
	}
	b, err = AddBill(uid, cid, b)
	if err != nil {
		return "", err
	} else if b == nil {
		return "", errors.New("nil_bill")
	}
	return results, err
}

// GetWXPayOrderInfo 生成微信订单，并同步bill
func GetWXPayOrderInfo(uid, cid int64, goodsType int) (*WXOrder, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// goods
	goods := getGoodsByType(goodsType)
	if goods == nil {
		return nil, errors.New("pay_goods_nil")
	}
	// order
	wxOrder, tradeNo, err := getWXPayWithTrade(uid, cid, goods)
	if err != nil {
		return wxOrder, err
	}
	// bill
	b := &entity.Bill{
		PlatformPay:  entity.BILL_PLATFORM_PAY_WX,
		PayType:      entity.BILL_PAY_TYPE_APP,
		PayAmount:    goods.Amount,
		TradeNo:      tradeNo,
		TradeReceipt: "",
		GoodsType:    goodsType,
	}
	b, err = AddBill(uid, cid, b)
	if err != nil {
		return wxOrder, err
	} else if b == nil {
		return wxOrder, errors.New("nil_bill")
	}
	return wxOrder, nil
}

// GetApplePayOrderInfo 生成苹果订单，并同步bill
func GetApplePayOrderInfo(uid, cid int64, goodsType int) (*AppleOrder, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	}
	// goods
	goods := getGoodsByType(goodsType)
	if goods == nil {
		return nil, errors.New("pay_goods_nil")
	}
	// order
	order := &AppleOrder{
		ProductId:     strings.TrimSpace(strconv.Itoa(goodsType)),
		TransactionId: "-",
		Receipt:       "",
	}
	// bill
	b := &entity.Bill{
		PlatformPay:  entity.BILL_PLATFORM_PAY_APPLE,
		PayType:      entity.BILL_PAY_TYPE_APP,
		PayAmount:    goods.Amount,
		TradeNo:      order.TransactionId,
		TradeReceipt: order.Receipt,
		GoodsType:    goodsType,
	}
	b, err := AddBill(uid, cid, b)
	if err != nil {
		return nil, err
	} else if b == nil {
		return nil, errors.New("nil_bill")
	}
	return order, nil
}

// AddBill
func AddBill(uid, cid int64, b *entity.Bill) (*entity.Bill, error) {
	if uid <= 0 {
		return nil, errors.New("nil_user")
	} else if cid <= 0 {
		return nil, errors.New("nil_couple")
	} else if b == nil {
		return nil, errors.New("nil_bill")
	} else if b.PlatformPay != entity.BILL_PLATFORM_PAY_ALI &&
		b.PlatformPay != entity.BILL_PLATFORM_PAY_WX &&
		b.PlatformPay != entity.BILL_PLATFORM_PAY_APPLE &&
		b.PlatformPay != entity.BILL_PLATFORM_PAY_GOOGLE {
		return nil, errors.New("pay_platform_err")
	} else if len(strings.TrimSpace(b.TradeNo)) <= 0 {
		return nil, errors.New("bill_trade_no_err")
	} else if b.GoodsType <= 0 {
		return nil, errors.New("pay_goods_nil")
	}
	// platform
	entry, err := mysql.GetEntryLatestByUser(uid)
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, errors.New("nil_entry")
	}
	// mysql
	b.UserId = uid
	b.CoupleId = cid
	b.PlatformOs = entry.Platform
	b.GoodsId = 0
	//b.TradePay = false
	//b.GoodsOut = false
	b, err = mysql.AddBill(b)
	return b, err
}

// GetBillList
func GetBillList(uid, cid int64, tradeNo string, page int) ([]*entity.Bill, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Bill
	offset := page * limit
	list, err := mysql.GetBillList(uid, cid, tradeNo, offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_bill")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性和同步
	return list, err
}

// GetBillListWithNoSync
func GetBillListWithNoSync(page int) ([]*entity.Bill, error) {
	if page < 0 {
		return nil, errors.New("limit_page_err")
	}
	// mysql
	limit := GetPageSizeLimit().Bill
	offset := page * limit
	list, err := mysql.GetBillListWithNoSync(offset, limit)
	if err != nil {
		return nil, err
	} else if list == nil || len(list) <= 0 {
		if page <= 0 {
			return nil, errors.New("no_data_bill")
		} else {
			return nil, nil
		}
	}
	// 没有额外属性和同步
	return list, err
}

// GetBillAmountTotalByCreateWithPay
func GetBillAmountTotalByCreateWithPay(start, end int64, platformOs string, platformPay, payType, goodsType int) float64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetBillAmountTotalByCreateWithPay(start, end, platformOs, platformPay, payType, goodsType)
	return total
}

// GetBillTotalByCreateWithDelApy
func GetBillTotalByCreateWithDelApy(start, end int64, platformOs string, platformPay, payType, goodsType int) int64 {
	if start >= end {
		return 0
	}
	// mysql
	total := mysql.GetBillTotalByCreateWithDelApy(start, end, platformOs, platformPay, payType, goodsType)
	return total
}

// CheckBillsStatus
func CheckBillsStatus(uid, cid int64, order *AppleOrder) error {
	if uid <= 0 {
		return errors.New("nil_user")
	} else if cid <= 0 {
		return errors.New("nil_couple")
	}
	// mysql
	list, err := mysql.GetBillListByUserCoupleWithNoSync(uid, cid, 0, 20)
	if err != nil {
		return err
	} else if list == nil || len(list) <= 0 {
		return errors.New("no_data_bill")
	}
	// 开始遍历订单
	for _, b := range list {
		if b == nil || len(b.TradeNo) <= 0 {
			continue
		}
		if b.PlatformPay == entity.BILL_PLATFORM_PAY_ALI {
			checkAliTrade(b)
		} else if b.PlatformPay == entity.BILL_PLATFORM_PAY_WX {
			checkWXTrade(b)
		} else if b.PlatformPay == entity.BILL_PLATFORM_PAY_APPLE {
			checkAppleTrade(b, order)
		} else if b.PlatformPay == entity.BILL_PLATFORM_PAY_GOOGLE {
			// google
		}
	}
	return err
}

// CheckBillStatus
func CheckBillStatus(bid int64) error {
	if bid <= 0 {
		return errors.New("nil_bill")
	}
	// bill
	b, err := mysql.GetBillById(bid)
	if err != nil {
		return err
	} else if b == nil || len(b.TradeNo) <= 0 {
		return errors.New("nil_bill")
	} else if b.UserId <= 0 {
		return errors.New("nil_user")
	} else if b.CoupleId <= 0 {
		return errors.New("nil_couple")
	}
	// 开始检查
	if b.PlatformPay == entity.BILL_PLATFORM_PAY_ALI {
		checkAliTrade(b)
	} else if b.PlatformPay == entity.BILL_PLATFORM_PAY_WX {
		checkWXTrade(b)
	} else if b.PlatformPay == entity.BILL_PLATFORM_PAY_APPLE {
		// apple 没有(控制台做不了)
	} else if b.PlatformPay == entity.BILL_PLATFORM_PAY_GOOGLE {
		// google
	}
	return err
}

// checkAliTrade 支付宝订单检查
func checkAliTrade(b *entity.Bill) {
	if b == nil || len(b.TradeNo) <= 0 {
		return
	}
	// 查询订单
	req := alipay.AliPayTradeQuery{
		OutTradeNo: b.TradeNo,
	}
	response, err := getAliPay().TradeQuery(req)
	if response != nil && response.IsSuccess() && err == nil {
		// 只看成功支付的订单
		query := response.AliPayTradeQuery
		if query.TradeNo != b.TradeNo && query.OutTradeNo != b.TradeNo {
			// 订单号不对
			return
		}
		payAmount, _ := strconv.ParseFloat(query.TotalAmount, 64)
		if payAmount <= 0 || payAmount < b.PayAmount {
			// 金额不对
			return
		}
		if strings.TrimSpace(query.TradeStatus) != alipay.K_TRADE_STATUS_TRADE_SUCCESS &&
			strings.TrimSpace(query.TradeStatus) != alipay.K_TRADE_STATUS_TRADE_FINISHED {
			// 没有支付
			return
		}
		// bill更新
		//if !b.TradePay {
		//	b.TradePay = true
		//	b, err = mysql.UpdateBill(b)
		//	if err != nil {
		//		return
		//	} else if b == nil {
		//		return
		//	}
		//}
		// 发货
		goodsCheck(b)
	}
}

// checkWXTrade 微信订单检查
func checkWXTrade(b *entity.Bill) {
	if b == nil || len(b.TradeNo) <= 0 {
		return
	}
	// 查询订单
	client := getWxPay()
	params := make(wxpay.Params)
	params.
		SetString("out_trade_no", b.TradeNo).
		SetString("nonce_str", createTradeNo())
	query, err := client.OrderQuery(params)
	if query != nil && query["return_code"] == wxpay.Success && query["result_code"] == wxpay.Success && err == nil {
		// 只看成功支付的订单
		if query["out_trade_no"] != b.TradeNo {
			// 订单号不对
			return
		}
		payAmount, _ := strconv.ParseFloat(query["total_fee"], 64)
		if payAmount <= 0 || payAmount < b.PayAmount*100 {
			// 金额不对
			return
		}
		if query["trade_state"] != wxpay.Success {
			// 没有支付
			return
		}
		// bill更新
		//if !b.TradePay {
		//	b.TradePay = true
		//	b, err = mysql.UpdateBill(b)
		//	if err != nil {
		//		return
		//	} else if b == nil {
		//		return
		//	}
		//}
		// 发货
		goodsCheck(b)
	}
}

// checkAppleTrade
func checkAppleTrade(b *entity.Bill, order *AppleOrder) {
	if b == nil || order == nil {
		return
	}
	// 查询订单
	bundleID := utils.GetConfigStr("conf", "app.conf", "phone", "bundle_ios")
	apps, err := apple.GetTradeInfo(bundleID, order.Receipt)
	if apps != nil && len(apps) > 0 && err == nil {
		for _, app := range apps {
			if app == nil {
				continue
			}
			if strings.TrimSpace(strconv.Itoa(b.GoodsType)) != strings.TrimSpace(order.ProductId) {
				// 商品类型不对
				continue
			}
			bill, err := mysql.GetBillByTradeNo(strings.TrimSpace(app.TransactionId), strings.TrimSpace(order.Receipt))
			if bill != nil || err != nil {
				// TradeNo只能加一次
				continue
			}
			// bill更新(这里会更新tradeNo)
			transactionId := strings.TrimSpace(app.TransactionId)
			orderReceipt := strings.TrimSpace(order.Receipt)
			if strings.TrimSpace(b.TradeNo) != transactionId && strings.TrimSpace(b.TradeReceipt) != orderReceipt {
				// bill更新(这里会更新tradeNo，且只能更新一次)
				b.TradeNo = transactionId
				b.TradeReceipt = orderReceipt
				b, err = mysql.UpdateBill(b)
				if err != nil {
					continue
				} else if b == nil {
					continue
				}
			}
			// 发货
			goodsCheck(b)
		}
	}
	// 审核时候开放
	//apps, err = apple.GetTradeInfoByDebug(bundleID, order.Receipt)
	//if apps != nil && len(apps) > 0 && err == nil {
	//	for _, app := range apps {
	//		if app == nil {
	//			continue
	//		}
	//		if strings.TrimSpace(strconv.Itoa(b.GoodsType)) != strings.TrimSpace(order.ProductId) {
	//			// 商品类型不对
	//			return
	//		}
	//		bill, err := mysql.GetBillByTradeNo(strings.TrimSpace(app.TransactionId), strings.TrimSpace(order.Receipt))
	//		if bill != nil || err != nil {
	//			// TradeNo只能加一次
	//			return
	//		}
	//		// bill更新(这里会更新tradeNo)
	//		transactionId := strings.TrimSpace(app.TransactionId)
	//		orderReceipt := strings.TrimSpace(order.Receipt)
	//		if strings.TrimSpace(b.TradeNo) != transactionId && strings.TrimSpace(b.TradeReceipt) != orderReceipt {
	//			// bill更新(这里会更新tradeNo，且只能更新一次)
	//			b.TradeNo = transactionId
	//			b.TradeReceipt = orderReceipt
	//			b, err = mysql.UpdateBill(b)
	//			if err != nil {
	//				continue
	//			} else if b == nil {
	//				continue
	//			}
	//		}
	//		// 发货
	//		goodsCheck(b)
	//	}
	//}
}

// goodsCheck 商品检查
func goodsCheck(b *entity.Bill) error {
	if b == nil {
		return errors.New("nil_bill")
	} else if b.UserId <= 0 {
		return errors.New("nil_user")
	} else if b.CoupleId <= 0 {
		return errors.New("nil_couple")
	}
	if b.CreateAt < 1577030400 {
		// 兼容没有goods_id的旧数据，2019-12-23 00:00:00
		return nil
	}
	// type
	typeVip1 := entity.BILL_AND_GOODS_TYPE_VIP_1
	typeVip2 := entity.BILL_AND_GOODS_TYPE_VIP_2
	typeVip3 := entity.BILL_AND_GOODS_TYPE_VIP_3
	typeCoin1 := entity.BILL_AND_GOODS_TYPE_COIN_1
	typeCoin2 := entity.BILL_AND_GOODS_TYPE_COIN_2
	typeCoin3 := entity.BILL_AND_GOODS_TYPE_COIN_3
	if b.PlatformPay == entity.BILL_PLATFORM_PAY_APPLE {
		typeVip1 = entity.BILL_IOS_GOODS_TYPE_VIP_1
		typeVip2 = entity.BILL_IOS_GOODS_TYPE_VIP_2
		typeVip3 = entity.BILL_IOS_GOODS_TYPE_VIP_3
		typeCoin1 = entity.BILL_IOS_GOODS_TYPE_COIN_1
		typeCoin2 = entity.BILL_IOS_GOODS_TYPE_COIN_2
		typeCoin3 = entity.BILL_IOS_GOODS_TYPE_COIN_3
	}
	// goods检查
	var err error
	if b.GoodsType == typeVip1 || b.GoodsType == typeVip2 || b.GoodsType == typeVip3 {
		// 先检查是否存在(防止重复添加)
		if b.GoodsId != 0 {
			g, err := mysql.GetVipById(b.GoodsId)
			if err != nil {
				return err
			} else if g != nil && g.Id != 0 {
				return nil
			}
		}
		// 获取expireDays
		limit := GetLimit()
		expireDays := 0
		if b.GoodsType == typeVip1 {
			expireDays = limit.PayVipGoods1Days
		} else if b.GoodsType == typeVip2 {
			expireDays = limit.PayVipGoods2Days
		} else if b.GoodsType == typeVip3 {
			expireDays = limit.PayVipGoods3Days
		}
		// 开始添加
		v := &entity.Vip{
			ExpireDays: expireDays,
		}
		v, err = AddVipByPay(b.UserId, b.CoupleId, b.Id, v)
		if v == nil || err != nil {
			return err
		}
		// 更新bill
		b.GoodsId = v.Id
		_, err = mysql.UpdateBill(b)
	} else if b.GoodsType == typeCoin1 || b.GoodsType == typeCoin2 || b.GoodsType == typeCoin3 {
		// 先检查是否存在(防止重复添加)
		if b.GoodsId != 0 {
			g, err := mysql.GetCoinById(b.GoodsId)
			if err != nil {
				return err
			} else if g != nil && g.Id != 0 {
				return nil
			}
		}
		// 获取count
		limit := GetLimit()
		count := 0
		if b.GoodsType == typeCoin1 {
			count = limit.PayCoinGoods1Count
		} else if b.GoodsType == typeCoin2 {
			count = limit.PayCoinGoods2Count
		} else if b.GoodsType == typeCoin3 {
			count = limit.PayCoinGoods3Count
		}
		// 开始添加
		c := &entity.Coin{
			Change: count,
		}
		c, err = AddCoinByPay(b.UserId, b.CoupleId, b.Id, c)
		if c == nil || err != nil {
			return err
		}
		// 更新bill
		b.GoodsId = c.Id
		_, err = mysql.UpdateBill(b)
	}
	return err
}

// getGoodsByType
func getGoodsByType(goodsType int) *Goods {
	limit := GetLimit()
	if goodsType == entity.BILL_AND_GOODS_TYPE_VIP_1 || goodsType == entity.BILL_IOS_GOODS_TYPE_VIP_1 {
		g := &Goods{}
		g.Type = goodsType
		g.Title = limit.PayVipGoods1Title
		g.Amount = limit.PayVipGoods1Amount
		return g
	} else if goodsType == entity.BILL_AND_GOODS_TYPE_VIP_2 || goodsType == entity.BILL_IOS_GOODS_TYPE_VIP_2 {
		g := &Goods{}
		g.Type = goodsType
		g.Title = limit.PayVipGoods2Title
		g.Amount = limit.PayVipGoods2Amount
		return g
	} else if goodsType == entity.BILL_AND_GOODS_TYPE_VIP_3 || goodsType == entity.BILL_IOS_GOODS_TYPE_VIP_3 {
		g := &Goods{}
		g.Type = goodsType
		g.Title = limit.PayVipGoods3Title
		g.Amount = limit.PayVipGoods3Amount
		return g
	} else if goodsType == entity.BILL_AND_GOODS_TYPE_COIN_1 || goodsType == entity.BILL_IOS_GOODS_TYPE_COIN_1 {
		g := &Goods{}
		g.Type = goodsType
		g.Title = limit.PayCoinGoods1Title
		g.Amount = limit.PayCoinGoods1Amount
		return g
	} else if goodsType == entity.BILL_AND_GOODS_TYPE_COIN_2 || goodsType == entity.BILL_IOS_GOODS_TYPE_COIN_2 {
		g := &Goods{}
		g.Type = goodsType
		g.Title = limit.PayCoinGoods2Title
		g.Amount = limit.PayCoinGoods2Amount
		return g
	} else if goodsType == entity.BILL_AND_GOODS_TYPE_COIN_3 || goodsType == entity.BILL_IOS_GOODS_TYPE_COIN_3 {
		g := &Goods{}
		g.Type = goodsType
		g.Title = limit.PayCoinGoods3Title
		g.Amount = limit.PayCoinGoods3Amount
		return g
	}
	return nil
}

// getAliPay 阿里支付类
func getAliPay() (client *alipay.AliPay) {
	appId := utils.GetConfigStr("conf", "third.conf", "ali-pay", "app_id")
	//notifyUrl :=
	aliPublicKey := utils.GetConfigStr("conf", "third.conf", "ali-pay", "ali_public_key")
	appPrivateKey := utils.GetConfigStr("conf", "third.conf", "ali-pay", "app_private_key")
	return alipay.New(appId, "", aliPublicKey, appPrivateKey, true)
}

// getAliPayWithTrade 阿里订单生成
func getAliPayWithTrade(uid, cid int64, goods *Goods) (alipay.AliPayTradeAppPay, string, error) {
	var p = alipay.AliPayTradeAppPay{}
	p.Body = goods.Title
	p.Subject = goods.Title
	p.TotalAmount = strconv.FormatFloat(goods.Amount, 'g', -1, 64)
	p.OutTradeNo = createTradeNo()
	p.NotifyURL = utils.GetConfigStr("conf", "third.conf", "ali-pay", "notify_url")
	//p.ProductCode = "p_1010101"
	results, err := getAliPay().TradeAppPay(p)
	return p, results, err
}

// getWxPay 微信支付类
func getWxPay() (client *wxpay.Client) {
	appId := utils.GetConfigStr("conf", "third.conf", "wx-pay", "app_id")
	partnerId := utils.GetConfigStr("conf", "third.conf", "wx-pay", "partner_id")
	appKey := utils.GetConfigStr("conf", "third.conf", "wx-pay", "app_key")
	account := wxpay.NewAccount(appId, partnerId, appKey, false)
	return wxpay.NewClient(account)
}

// getWXPayWithTrade 微信订单生成
func getWXPayWithTrade(uid, cid int64, goods *Goods) (*WXOrder, string, error) {
	appName := utils.GetConfigStr("conf", "app.conf", "common", "app_name")
	body := appName + "-" + goods.Title
	tradeNo := createTradeNo()
	nonce := createTradeNo()
	amount := int64(goods.Amount * 100)
	notifyUrl := utils.GetConfigStr("conf", "third.conf", "wx-pay", "notify_url")
	billIp := "127.0.0.1"
	// 开始请求
	client := getWxPay()
	params := make(wxpay.Params)
	params.SetString("body", body).
		SetString("out_trade_no", tradeNo).
		SetString("nonce_str", nonce).
		SetInt64("total_fee", amount).
		SetString("spbill_create_ip", billIp).
		SetString("notify_url", notifyUrl).
		SetString("trade_type", "APP")
	order, err := client.UnifiedOrder(params)
	if order == nil || err != nil {
		return nil, "", err
	} else if order.GetString("return_code") == wxpay.Fail {
		return nil, "", errors.New(order.GetString("return_msg"))
	}
	// 开始赋值
	var wxOrder = &WXOrder{}
	wxOrder.AppId = order["appid"]
	wxOrder.PartnerId = order["mch_id"]
	wxOrder.PrepayId = order["prepay_id"]
	wxOrder.PackageValue = "Sign=WXPay"
	wxOrder.NonceStr = createTradeNo()
	wxOrder.TimeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	// app端签名
	paramsApp := make(wxpay.Params)
	paramsApp.
		SetString("appid", wxOrder.AppId).
		SetString("partnerid", wxOrder.PartnerId).
		SetString("prepayid", wxOrder.PrepayId).
		SetString("package", wxOrder.PackageValue).
		SetString("noncestr", wxOrder.NonceStr).
		SetString("timestamp", wxOrder.TimeStamp)
	// sign
	wxOrder.Sign = client.Sign(paramsApp)
	return wxOrder, tradeNo, err
}

// createTradeNo 创建订单号
func createTradeNo() string {
	// 前半截时间戳
	unixNa := time.Now().UnixNano()
	unix := strconv.FormatInt(unixNa, 16)
	// 后半截随机数
	length := 13
	max := math.Pow10(length) - 1
	min := math.Pow10(length - 1)
	rand13 := utils.GetRandRange(int(max), int(min))
	return unix + "_" + strconv.Itoa(rand13)
}
