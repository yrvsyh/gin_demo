package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/yrvsyh/gin_demo/database"
	"github.com/yrvsyh/gin_demo/utils"
)

type (
	t struct{}

	JWTClaims struct {
		jwt.StandardClaims
		Name string `json:"name,omitempty"`
		Role int    `json:"role,omitempty"`
	}
)

var (
	JWTAuth = t{}

	tokenSecretKey = []byte("secret_key")
	// token过期时间
	tokenExpireDuration = time.Hour
	// token可刷新时间
	tokenRefreshDuration = time.Hour * 24 * 30
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := JWTAuth.GetToken(c)
		claims, err := JWTAuth.ParseToken(tokenString)
		if err != nil {
			utils.Error(c, utils.ERR_TOKEN_PARSE_FAILD)
			c.Abort()
		} else {
			if err := claims.Valid(); err == nil {
				log.Debug().Str("claimsId", claims.Id).Str("name", claims.Name).Msg("TOKEN VALID")
				// token黑名单判断
				if _, err := database.RDB.Get(claims.StandardClaims.Id).Result(); err == redis.Nil {
					c.Next()
					return
				}
			}
			utils.Error(c, utils.ERR_TOKEN_INVALID)
			c.Abort()
		}

	}
}

// 从http头获取token
func (t) GetToken(c *gin.Context) string {
	tokenString := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		authHeaders := strings.Split(authHeader, " ")
		if len(authHeaders) == 2 && authHeaders[0] == "Bearer" {
			tokenString = authHeaders[1]
		}
	}
	return tokenString
}

// 生成token
func (t) GenToken(name string, role int) (string, error) {
	claims := &JWTClaims{
		Name: name,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSecretKey)
}

// 解析token
func (t) ParseToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenSecretKey, nil
	})
	return claims, err
}

// 刷新token
func (t) RefreshToken(tokenString string) (string, error) {
	claims, err := JWTAuth.ParseToken(tokenString)
	if err == nil {
		if err = claims.Valid(); err == nil {
			if time.Unix(claims.ExpiresAt, 0).Add(tokenRefreshDuration).After(time.Now()) {
				claims.StandardClaims.Id = uuid.NewV4().String()
				// 重置过期时间
				claims.ExpiresAt = time.Now().Add(tokenExpireDuration).Unix()
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				return token.SignedString(tokenSecretKey)
			}
		}
	}
	return "", err
}

// 使token失效, 添加至redis黑名单
func (t) DelToken(tokenString string) error {
	claims, err := JWTAuth.ParseToken(tokenString)
	if err == nil {
		if JWTAuth.VerifyToken(tokenString) {
			tokenId := claims.StandardClaims.Id
			expireTime := time.Until(time.Unix(claims.ExpiresAt, 0))
			err = database.RDB.Set(tokenId, 1, expireTime).Err()
		}
	}
	return err
}

// 验证token是否过期
func (t) VerifyToken(tokenString string) bool {
	claims, err := JWTAuth.ParseToken(tokenString)
	if err == nil {
		if err = claims.Valid(); err == nil {
			return true
		}
	}
	return false
}
