package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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
			ExpiresAt: int64(key + 5184000),
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

func GetAvailClassesTimesAlgo(classes []models.ClassMainInfo) map[string]struct {
	General         bool `json:"general"`
	HalfHourFree    bool `json:"half_hour_free"`
	HourFree        bool `json:"hour_free"`
	OneHalfHourFree bool `json:"one_half_hour_free"`
	TwoHourFree     bool `json:"two_hour_free"`
	TwoHalfHourFree bool `json:"two_half_hour_free"`
} {
	timeTable := genTimeTable()
	timeLayout := "15:04"
	for _, class := range classes {
		entry := timeTable[class.ClassTime]
		entry.General = false
		timeTable[class.ClassTime] = entry
		classTime, err := time.Parse(timeLayout, class.ClassTime)
		if err != nil {
			log.Printf("err parsing time: %v", err)
			continue
		}
		// todo: придумать нормальный алгос для перебора и установки досутпного времени
		classTime = classTime.Add(-30 * time.Minute)
		result := classTime.Format(timeLayout)
		entry = timeTable[result]
		entry.HalfHourFree = false
		timeTable[result] = entry

		classTime, _ = time.Parse(timeLayout, result)
		classTime = classTime.Add(-30 * time.Minute)
		result = classTime.Format(timeLayout)
		entry = timeTable[result]
		entry.HourFree = false
		timeTable[result] = entry

		classTime, _ = time.Parse(timeLayout, result)
		classTime = classTime.Add(-30 * time.Minute)
		result = classTime.Format(timeLayout)
		entry = timeTable[result]
		entry.OneHalfHourFree = false
		timeTable[result] = entry

		classTime, _ = time.Parse(timeLayout, result)
		classTime = classTime.Add(-30 * time.Minute)
		result = classTime.Format(timeLayout)
		entry = timeTable[result]
		entry.TwoHourFree = false
		timeTable[result] = entry

		classTime, _ = time.Parse(timeLayout, result)
		classTime = classTime.Add(-30 * time.Minute)
		result = classTime.Format(timeLayout)
		entry = timeTable[result]
		entry.TwoHalfHourFree = false
		timeTable[result] = entry
	}
	return timeTable
}

func genTimeTable() map[string]struct {
	General         bool `json:"general"`
	HalfHourFree    bool `json:"half_hour_free"`
	HourFree        bool `json:"hour_free"`
	OneHalfHourFree bool `json:"one_half_hour_free"`
	TwoHourFree     bool `json:"two_hour_free"`
	TwoHalfHourFree bool `json:"two_half_hour_free"`
} {
	var (
		timeTable = make(map[string]struct {
			General         bool `json:"general"`
			HalfHourFree    bool `json:"half_hour_free"`
			HourFree        bool `json:"hour_free"`
			OneHalfHourFree bool `json:"one_half_hour_free"`
			TwoHourFree     bool `json:"two_hour_free"`
			TwoHalfHourFree bool `json:"two_half_hour_free"`
		})
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
		timeTable[timeString] = struct {
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

func SendTgTechPgPingAlert(token, shard, attempts string) {
	buf, err := os.Open("pkg/storage/pg_bad_ping.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer buf.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	chatID := "1242802644"
	_ = writer.WriteField("chat_id", chatID)
	_ = writer.WriteField("parse_mode", "markdown")
	_ = writer.WriteField("caption",
		"*ALERT*\r\nOne of the shard was currently unavailable for _"+
			attempts+"_ attempts.\r\n\r\n*Shard name*: "+shard+"\r\n*Time*: "+
			time.Now().Format(time.RFC3339))

	fw, err := writer.CreateFormFile("photo", "pg_bad_ping.jpg")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(fw, buf)
	if err != nil {
		log.Fatal(err)
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "https://api.telegram.org/bot"+token+"/sendPhoto", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	rsp, _ := client.Do(req)

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}

	cmd := exec.Command("docker start", shard)
	err = cmd.Run()
	log.Printf("Command finished with error: %v", err)
}
