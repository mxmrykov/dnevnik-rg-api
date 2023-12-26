package utils

import (
	"bytes"
	"crypto/md5"
	requests "dnevnik-rg.ru/internal/models/request"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"os"
	"strconv"
	"time"
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
