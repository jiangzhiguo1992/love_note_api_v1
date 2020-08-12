package apple

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"libs/utils"
	"net/http"
	"strings"
)

const (
	URL_SANDBOX = "https://sandbox.itunes.apple.com/verifyReceipt"
	URL_VERIFY  = "https://buy.itunes.apple.com/verifyReceipt"
)

type (
	TradeInfo struct {
		Status      int      `json:"status"`
		Environment string   `json:"environment"`
		Receipt     *Receipt `json:"receipt"`
	}
	Receipt struct {
		ReceiptType                string   `json:"receipt_type"`
		AdamId                     int      `json:"adam_id"`
		AppItemId                  int      `json:"app_item_id"`
		BundleId                   string   `json:"bundle_id"`
		ApplicationVersion         string   `json:"application_version"`
		DownloadId                 int      `json:"download_id"`
		VersionExternalIdentifier  int      `json:"version_external_identifier"`
		ReceiptCreationDate        string   `json:"receipt_creation_date"`
		ReceiptCreationDateMs      string   `json:"receipt_creation_date_ms"`
		ReceiptCreationDatePst     string   `json:"receipt_creation_date_pst"`
		RequestDate                string   `json:"request_date"`
		RequestDateMs              string   `json:"request_date_ms"`
		RequestDatePst             string   `json:"request_date_pst"`
		OriginalPurchaseDate       string   `json:"original_purchase_date"`
		OriginalPurchaseDateMs     string   `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePst    string   `json:"original_purchase_date_pst"`
		OriginalApplicationVersion string   `json:"original_application_version"`
		InApp                      []*InApp `json:"in_app"`
	}
	InApp struct {
		Quantity                string `json:"quantity"`
		ProductId               string `json:"product_id"`
		TransactionId           string `json:"transaction_id"`
		OriginalTransactionId   string `json:"original_transaction_id"`
		PurchaseDate            string `json:"purchase_date"`
		PurchaseDateMs          string `json:"purchase_date_ms"`
		PurchaseDatePst         string `json:"purchase_date_pst"`
		OriginalPurchaseDate    string `json:"original_purchase_date"`
		OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
		IsTrialPeriod           string `json:"is_trial_period"`
	}
)

// 获取订单状态
func GetTradeInfo(bundleId, receipt string) ([]*InApp, error) {
	var receiptData = "{\"receipt-data\":\"" + receipt + "\"}"
	resp, err := http.Post(URL_VERIFY, "application/json", bytes.NewBuffer([]byte(receiptData)))
	if err != nil {
		utils.LogErr("applePay", err)
		return make([]*InApp, 0), errors.New("data_err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.LogErr("applePay", err)
		return make([]*InApp, 0), errors.New("data_err")
	}
	// 格式化返回信息
	result := &TradeInfo{}
	if err := json.Unmarshal(body, result); err != nil || result == nil {
		utils.LogErr("applePay", err)
		return make([]*InApp, 0), errors.New("data_decode_err")
	}
	if result.Receipt == nil {
		return make([]*InApp, 0), errors.New("no_data_bill")
	}
	receiptResult := result.Receipt
	if strings.TrimSpace(receiptResult.BundleId) != strings.TrimSpace(bundleId) {
		return make([]*InApp, 0), errors.New("no_data_bill")
	}
	return receiptResult.InApp, nil
}

// 获取订单状态(沙河环境)
func GetTradeInfoByDebug(bundleId, receipt string) ([]*InApp, error) {
	var receiptData = "{\"receipt-data\":\"" + receipt + "\"}"
	resp, err := http.Post(URL_SANDBOX, "application/json", bytes.NewBuffer([]byte(receiptData)))
	if err != nil {
		utils.LogErr("applePay", err)
		return make([]*InApp, 0), errors.New("data_err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.LogErr("applePay", err)
		return make([]*InApp, 0), errors.New("data_err")
	}
	// 格式化返回信息
	result := &TradeInfo{}
	if err := json.Unmarshal(body, result); err != nil || result == nil {
		utils.LogErr("applePay", err)
		return make([]*InApp, 0), errors.New("data_decode_err")
	}
	if result.Receipt == nil {
		return make([]*InApp, 0), errors.New("no_data_bill")
	}
	receiptResult := result.Receipt
	if strings.TrimSpace(receiptResult.BundleId) != strings.TrimSpace(bundleId) {
		return make([]*InApp, 0), errors.New("no_data_bill")
	}
	return receiptResult.InApp, nil
}
