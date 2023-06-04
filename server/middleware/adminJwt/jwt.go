package adminJwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     -1,
				"error_code": "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":     -1,
					"error_code": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     -1,
				"error_code": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error  = errors.New("token is expired")
	TokenNotValidYet error  = errors.New("token not active yet")
	TokenMalformed   error  = errors.New("that's not even a token")
	TokenInvalid     error  = errors.New("couldn't handle this token:")
	SignKey          string = "0YETlygku5hB"
)

type CustomClaims struct {
	ID uint `json:"user_id"`
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return SignKey
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func GenerateToken(userID uint) (string, error) {
	j := JWT{
		[]byte(GetSignKey()),
	}
	claims := CustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 128*24*3600), // 过期时间 一小时
		},
	}
	return j.CreateToken(claims)
}
