package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
	"tree-hollow/config"
)

// 自己简单地造了个轮子
// 签名算法使用的是 HMAC SHA256

const (
	HeaderPlain = "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"
)

// Claims Payload
type Claims struct {
	InterArrivalTime int64  `json:"iat"` // 到达时间
	ExpirationDate   int64  `json:"exp"` // 认证时间
	UserName         string `json:"user_name"`
}

func GenerateTokenPairWithUserName(userName string) (token, refreshToken string, err error) {
	now := time.Now().Unix()
	tokenClaims := Claims{
		InterArrivalTime: now,
		ExpirationDate:   config.Config.JWTTimeOut,
		UserName:         userName,
	}

	// refreshToken时间会比token要长
	refreshTokenClaim := Claims{
		InterArrivalTime: now,
		ExpirationDate:   config.Config.JWTTimeOut * 10,
		UserName:         userName,
	}

	token, err = generateByClaims(tokenClaims)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = generateByClaims(refreshTokenClaim)
	if err != nil {
		return "", "", err
	}
	return
}

func generateByClaims(claims Claims) (string, error) {
	bytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := base64.StdEncoding.EncodeToString([]byte(HeaderPlain))
	payload := base64.StdEncoding.EncodeToString(bytes)

	return SignJWT(header, payload), nil
}

func SignJWT(header, payload string) string {

	hash := hmac.New(sha256.New, []byte(config.Config.JWTKey))
	hash.Write([]byte(header + payload))

	signed := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))

	return header + "." + payload + "." + signed
}

func AuthorizeJWT(jwtStr string) (bool, error) {
	claims := Claims{}

	parts := strings.Split(jwtStr, ".")

	payload, _ := base64.StdEncoding.DecodeString(parts[1])
	signed := parts[2]
	dSigned := strings.Split(SignJWT(parts[0], parts[1]), ".")[2]
	if signed != dSigned {
		return false, nil
	}
	err := json.Unmarshal(payload, &claims)
	if err != nil {
		return false, err
	}
	now := time.Now().Unix()

	// 如果现在的时间减去上一次登录时间大于认证时间
	if claims.ExpirationDate < now-claims.InterArrivalTime {
		return false, nil
	}
	return true, nil
}
