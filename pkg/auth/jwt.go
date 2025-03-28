package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims 自定义JWT载荷
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT配置
var (
	jwtSecret       = []byte("your-jwt-secret-key") // TODO: 从配置文件读取
	accessTokenExp  = time.Hour * 24                // 访问令牌有效期24小时
	refreshTokenExp = time.Hour * 24 * 7            // 刷新令牌有效期7天
)

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (accessToken, refreshToken string, err error) {
	// 创建访问令牌
	accessClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "huiyuanguanli",
			Subject:   "access_token",
		},
	}

	// 创建刷新令牌
	refreshClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "huiyuanguanli",
			Subject:   "refresh_token",
		},
	}

	// 生成访问令牌
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新访问令牌
func RefreshToken(refreshToken string) (string, error) {
	// 解析刷新令牌
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 验证是否是刷新令牌
	if claims.Subject != "refresh_token" {
		return "", errors.New("无效的刷新令牌")
	}

	// 生成新的访问令牌
	accessClaims := JWTClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "huiyuanguanli",
			Subject:   "access_token",
		},
	}

	// 生成新的访问令牌
	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
