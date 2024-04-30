package interceptor

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/service"
	"apiProject/api/response"
	"apiProject/api/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

type UserServiceInterfaceAuth struct {
	userInter service.UserServiceInterface
}

// NewUserServiceInterfaceAuth 创建用户服务器请求
func NewUserServiceInterfaceAuth(e service.UserServiceInterface) *UserServiceInterfaceAuth {
	return &UserServiceInterfaceAuth{userInter: e}
}

// WithJWTAuthorization 验证jwt凭证
//
// 参数
//
//	handlerFunc (http.HandlerFunc): http处理函数
//	u (service.UserServiceInterface): 用户服务接口
//
// 返回
//
//	http.HandlerFunc
func WithJWTAuthorization(handlerFunc http.HandlerFunc, u service.UserServiceInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tokenStr := GetTokenFromRequest(request)
		token, err := ValidateJWT(tokenStr)
		if err != nil {
			response.WriteJson(writer, response.FailMessageResp("token不正确"))
			return
		}

		if !token.Valid {
			response.WriteJson(writer, response.FailMessageResp("token验证失败"))
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userId"].(string)

		user, err := u.GetUserById(utils.ConvertToInt64(userId))
		if err != nil {
			response.WriteJson(writer, response.FailMessageResp("未查询到用户信息"))
			return
		}

		log.Println(user)
		handlerFunc(writer, request)
	}
}

// GetTokenFromRequest 从请求中获取凭证
func GetTokenFromRequest(request *http.Request) string {
	token := request.Header.Get("Authorization")
	if token != "" {
		return token
	}
	return ""
}

// ValidateJWT 验证JWT凭证有效性
func ValidateJWT(t string) (*jwt.Token, error) {
	secret := config.EnvConfig.JWTSecret
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("验证token失败 %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
