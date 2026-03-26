package jwt

import (
	"go-demo-server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义 Token 中携带的数据结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func GetJWTSecret() []byte {
	// var secret = []byte("secret")
	return []byte(config.Conf.JWT.Secret)
}

func GenerateToken(userID uint, username string) string {
	duration := time.Duration(config.Conf.JWT.ExpireTime) * time.Hour

	claims := Claims{
		UserID:   userID,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return ""
	}
	return tokenString
}

// parseToken 解析和验证 Token 的辅助函数
func ParseToken(tokenString string) (*jwt.Token, error) {
	// 定义你的密钥（实际项目中最好放在环境变量里）
	// var jwtKey = []byte("your_secret_key")
	var jwtKey = GetJWTSecret()

	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return jwtKey, nil
	})
}
