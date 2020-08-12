package utils

import (
	"math/rand"
	"time"
	//"sync/atomic"
)

func GetRandNum() int {
	t := time.Now().UnixNano()       //当前时间纳秒数
	r := rand.New(rand.NewSource(t)) //随机数源
	return r.Int()
}

//0-num
func GetRandMax(max int) int {
	t := time.Now().UnixNano()       //当前时间纳秒数
	r := rand.New(rand.NewSource(t)) //随机数源
	return r.Intn(max)
}

//0-num
func GetRandMax64(max int64) int64 {
	t := time.Now().UnixNano()       //当前时间纳秒数
	r := rand.New(rand.NewSource(t)) //随机数源
	return r.Int63n(max)
}

//min-max --> (-2,-100)
func GetRandRange(max, min int) int {
	t := time.Now().UnixNano()       //当前时间纳秒数
	r := rand.New(rand.NewSource(t)) //随机数源
	return r.Intn(max-min) + min
}

//func GetUuid(curTime time.Time) string {
//	var clockSeq uint32
//	var hardwareAddr []byte
//	var timeBase = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()
//
//	var u [16]byte
//	utcTime := curTime.In(time.UTC)
//	t := uint64(utcTime.Unix()-timeBase)*10000000 + uint64(utcTime.Nanosecond()/100)
//	u[0], u[1], u[2], u[3] = byte(t>>24), byte(t>>16), byte(t>>8), byte(t)
//	u[4], u[5] = byte(t>>40), byte(t>>32)
//	u[6], u[7] = byte(t>>56)&0x0F, byte(t>>48)
//	clock := atomic.AddUint32(&clockSeq, 1)
//	u[8] = byte(clock >> 8)
//	u[9] = byte(clock)
//	copy(u[10:], hardwareAddr)
//	u[6] |= 0x10 // set version to 1 (time based uuid)
//	u[8] &= 0x3F // clear variant
//	u[8] |= 0x80 // set to IETF variant
//
//	var offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}
//	const hexString = "0123456789abcdef"
//	r := make([]byte, 36)
//	for i, b := range u {
//		r[offsets[i]] = hexString[b>>4]
//		r[offsets[i]+1] = hexString[b&0xF]
//	}
//	r[8] = '-'
//	r[13] = '-'
//	r[18] = '-'
//	r[23] = '-'
//	return string(r)
//}
