package controllers

import (
	"libs/utils"
	"models/entity"
	"net/http"
	"services"
	"time"
)

func HandlerMoreHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetMoreHome(w, r)
	} else {
		response405(w, r)
	}
}

// GetMoreHome
func GetMoreHome(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// broadcastList
	broadcastList := make([]*entity.Broadcast, 0)
	if getModelStatus("more_broadcast") {
		broadcastList, _ = services.GetBroadcastListNoEnd()
	}
	// vip
	var vip *entity.Vip
	if getModelStatus("more_vip") && couple != nil && couple.Id > 0 {
		vip, _ = services.GetVipLatest(couple.Id)
	}
	// coin
	var coin *entity.Coin
	if getModelStatus("more_coin") && couple != nil && couple.Id > 0 {
		coin, _ = services.GetCoinLatest(couple.Id)
	}
	// sign
	var sign *entity.Sign
	if getModelStatus("more_sign") && couple != nil && couple.Id > 0 {
		now := utils.GetCSTDateByUnix(time.Now().Unix())
		sign, _ = services.GetSignByCoupleYearMonthDay(couple.Id, now.Year(), int(now.Month()), now.Day())
	}
	// mathPeriod
	var wifePeriod *entity.MatchPeriod
	var letterPeriod *entity.MatchPeriod
	var discussPeriod *entity.MatchPeriod
	if getModelStatus("more_match_period") {
		wifePeriod, _ = services.GetMatchPeriodNow(entity.MATCH_KIND_WIFE_PICTURE)
		letterPeriod, _ = services.GetMatchPeriodNow(entity.MATCH_KIND_LETTER_SHOW)
		discussPeriod, _ = services.GetMatchPeriodNow(entity.MATCH_KIND_DISCUSS_MEET)
	}
	response200Data(w, r, struct {
		BroadcastList []*entity.Broadcast `json:"broadcastList"`
		Vip           *entity.Vip         `json:"vip"`
		Coin          *entity.Coin        `json:"coin"`
		Sign          *entity.Sign        `json:"sign"`
		WifePeriod    *entity.MatchPeriod `json:"wifePeriod"`
		LetterPeriod  *entity.MatchPeriod `json:"letterPeriod"`
		DiscussPeriod *entity.MatchPeriod `json:"discussPeriod"`
	}{broadcastList,
		vip,
		coin,
		sign,
		wifePeriod,
		letterPeriod,
		discussPeriod,
	})
}
