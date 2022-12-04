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
	//从config.yml中获取服务器设定的过期间隔
	expiresHours, _ := strconv.ParseInt(global.VP.GetString("jwt.expiresHours"), 10, 64)
	//config.yml中expiresHours: 4800
	claims := jwt.StandardClaims{
		//指定token的发行人
		Issuer: "?",
		//token过期时间=服务器设定的过期间隔 + 当前时间
		ExpiresAt: expiresHours*60*60 + time.Now().Unix(),
		Audience:  strconv.FormatUint(id, 10),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	//从config.yml获取签名
	signature := global.VP.GetString("jwt.signature")
	//生成签名字符串
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
