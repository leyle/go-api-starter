package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type CurTime struct {
	Second    int64  `json:"second" bson:"second"`
	HumanTime string `json:"humanTime" bson:"humanTime"`
}

func GetCurTime() *CurTime {
	curT := time.Now()

	t := &CurTime{
		Second:    curT.Unix(),
		HumanTime: curT.Format("2006-01-02 15:04:05"),
	}

	return t
}

func CurUnixTime() int64 {
	return time.Now().Unix()
}

func CurMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func CurHumanTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurNoSpaceTime() string {
	return time.Now().Format("20060102150405")
}

func FmtTimestampTime(sec int64) string {
	tm := time.Unix(sec, 0)
	return tm.Format("2006-01-02 15:04:05")
}

type BaseStruct struct {
	Id      string   `json:"id" bson:"_id"`
	CreateT *CurTime `json:"createT" bson:"createT"`
	UpdateT *CurTime `json:"updateT" bson:"updateT"`
}

func Sha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5(data string) string {
	m := md5.New()
	m.Write([]byte(data))
	return hex.EncodeToString(m.Sum(nil))
}

func GenerateDataId() string {
	id := uuid.New().String()
	return id
}

func GenerateHashPasswd(loginId, rawPasswd string) string {
	d := strings.ToLower(loginId) + rawPasswd
	h := Sha256(d)
	return h
}

func GenerateToken(userId string) string {
	d := fmt.Sprintf("%s%d", userId, time.Now().Nanosecond())
	h := Md5(d)
	h = strings.ToUpper(h)
	return h
}
