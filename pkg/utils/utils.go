package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"dnevnik-rg.ru/internal/models"
	requests "dnevnik-rg.ru/internal/models/request"
	"github.com/golang-jwt/jwt"
)

func NewPassword() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(10000))))
	checkSum := hex.EncodeToString(sum[:])
	return checkSum[:7]
}

func GetKey() int64 {
	return time.Now().Unix()
}

func SetLongJwt(key int, checksum, role string) (string, error) {
	claims := requests.JwtPayload{
		Key:      key,
		CheckSum: checksum,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 5184000,
			Issuer:    os.Getenv("DEPLOY"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return ss, err
}

func SetShortJwt(key int, checksum, role string) (string, error) {
	claims := requests.JwtPayload{
		Key:      key,
		CheckSum: checksum,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(key + 3600),
			Issuer:    os.Getenv("DEPLOY"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return ss, err
}

func GetJwtPayload(token string) (*requests.JwtPayload, error) {
	parsedToken, errParseToken := jwt.ParseWithClaims(token, &requests.JwtPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := parsedToken.Claims.(*requests.JwtPayload); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return &requests.JwtPayload{}, errParseToken
	}
}

func GenerateUpdateSql(table string, key int, newParams []string, newValues []string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("UPDATE %s SET", table))
	for i, param := range newParams {
		buf.WriteString(fmt.Sprintf(" %s='%s'", param, newValues[i]))
		if i < len(newValues)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(fmt.Sprintf(" WHERE key=%d", key))
	return buf.String()
}

func GetNearestBdays(bList []requests.BirthDayList) []requests.BirthDayList {
	var (
		bdayList []struct {
			index int
			bday  time.Time
		}
		finalBdayList []requests.BirthDayList
	)
	for i, e := range bList {
		bday, _ := time.Parse(time.RFC3339, e.Birthday)
		bdayList = append(bdayList, struct {
			index int
			bday  time.Time
		}{index: i, bday: bday})
	}
	timeNow := time.Now()
	today := time.Date(
		timeNow.Year(),
		timeNow.Month(),
		timeNow.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	weekAfterTime := time.Date(
		timeNow.Year(),
		timeNow.Month(),
		timeNow.Day()+7,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	for _, e := range bdayList {
		birthday := time.Date(
			timeNow.Year(),
			e.bday.Month(),
			e.bday.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)
		if birthday.After(today) && birthday.Before(weekAfterTime) {
			finalBdayList = append(finalBdayList, requests.BirthDayList{
				Key:      bList[e.index].Key,
				Fio:      bList[e.index].Fio,
				Birthday: bList[e.index].Birthday,
			})
		}
	}
	return finalBdayList
}

type schedule map[string]*struct {
	General         bool `json:"general"`
	HalfHourFree    bool `json:"half_hour_free"`
	HourFree        bool `json:"hour_free"`
	OneHalfHourFree bool `json:"one_half_hour_free"`
	TwoHourFree     bool `json:"two_hour_free"`
	TwoHalfHourFree bool `json:"two_half_hour_free"`
}

func GetAvailClassesTimesAlgo(classes []models.ClassMainInfo) schedule {
	timeTable := genTimeTable()
	timeLayout := "15:04"
	for _, class := range classes {
		timeTable[class.ClassTime].General = false
		classTime, err := time.Parse(timeLayout, class.ClassTime)
		if err != nil {
			log.Printf("err parsing time: %v", err)
			continue
		}
		class.ClassDuration = strings.Replace(class.ClassDuration, ":", "h", 1) + "m"
		classDur, err := time.ParseDuration(class.ClassDuration)
		if err != nil {
			log.Printf("err parsing time: %v", err)
			continue
		}
		for i := classTime; i.Before(
			classTime.Add(classDur + 30*time.Minute),
		); i = i.Add(30 * time.Minute) {
			if i.Hour() >= 20 {
				continue
			}
			classTimeString := i.Format(timeLayout)
			timeTable[classTimeString].General = false
		}
		classTime = classTime.Add(-150 * time.Minute)
		timeTable.setTimes(classTime)
	}
	return timeTable
}

func genTimeTable() schedule {
	var (
		timeTable  = make(schedule)
		lastHour   = 9
		timeString string
	)
	for i := 0; i < 23; i += 1 {
		if i%2 == 0 {
			if i < 2 {
				timeString = fmt.Sprintf("0%d:00", lastHour)
			} else {
				timeString = fmt.Sprintf("%d:00", lastHour)
			}
		} else {
			if i < 2 {
				timeString = fmt.Sprintf("0%d:30", lastHour)
			} else {
				timeString = fmt.Sprintf("%d:30", lastHour)
			}
			lastHour += 1
		}
		timeTable[timeString] = &struct {
			General         bool `json:"general"`
			HalfHourFree    bool `json:"half_hour_free"`
			HourFree        bool `json:"hour_free"`
			OneHalfHourFree bool `json:"one_half_hour_free"`
			TwoHourFree     bool `json:"two_hour_free"`
			TwoHalfHourFree bool `json:"two_half_hour_free"`
		}{General: true, HalfHourFree: true, HourFree: true, OneHalfHourFree: true, TwoHourFree: true, TwoHalfHourFree: true}
	}
	return timeTable
}

func (s *schedule) setTimes(classTime time.Time) {
	for i := 1; i < 6; i++ {
		ct := classTime.Add(30 * time.Minute * time.Duration(i))
		if ct.Hour() < 9 {
			continue
		}
		result := (*s)[ct.Format("15:04")]
		for j := 1; j < i+1; j++ {
			reflect.ValueOf(result).Elem().FieldByIndex([]int{6 - j}).SetBool(false)
		}
		(*s)[ct.Format("15:04")] = result
	}
}

func HashSumGen(key int, checksum string) string {
	secret := key * (key % 5)
	sum := md5.Sum(
		[]byte(checksum + strconv.Itoa(secret)),
	)
	return hex.EncodeToString(sum[:])
}
