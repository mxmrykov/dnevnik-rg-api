package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func NewPassword() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(10000))))
	checkSum := hex.EncodeToString(sum[:])
	return checkSum[:7]
}

func GetUdid() int64 {
	return time.Now().Unix()
}
