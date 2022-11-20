package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"

	"IShare/global"
)

// GenerateToken 生成一个token
func GenerateToken(id uint64) (signedToken string) {
	expiresHours, _ := strconv.ParseInt(global.VP.GetString("jwt.expiresHours"), 10, 64)
	claims := jwt.StandardClaims{
		Issuer:    "?",
		ExpiresAt: expiresHours*60*60 + time.Now().Unix(),
		Audience:  strconv.FormatUint(id, 10),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signature := global.VP.GetString("jwt.signature")
	signedToken, _ = token.SignedString([]byte(signature))
	return signedToken
}

// ParseToken 验证token的正确性，正确则返回id
func ParseToken(signedToken string) (id uint64, err error) {
	signature := global.VP.GetString("jwt.signature")
	token, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signature), nil
		},
	)
	if err != nil || !token.Valid {
		err = errors.New("token isn't valid")
		return
	}
	id, err = strconv.ParseUint(token.Claims.(jwt.MapClaims)["aud"].(string), 10, 64)
	if err != nil {
		err = errors.New("token isn't valid")
	}
	return id, err
}
