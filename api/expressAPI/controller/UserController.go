package controller

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/interceptor"
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type UserServiceInterfaceTest struct {
	userInter service.UserServiceInterface
}

// NewUserServiceInterfaceTest 创建用户服务器请求
func NewUserServiceInterfaceTest(e service.UserServiceInterface) *UserServiceInterfaceTest {
	return &UserServiceInterfaceTest{userInter: e}
}

// RegisterRoutes 注册快递服务请求路由
func (e *UserServiceInterfaceTest) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/{dataId}", e.handlerGetUser).Methods("GET")
	r.HandleFunc("/user/create", e.handlerCreateUser).Methods("POST")
	r.HandleFunc("/user/checkPwd", e.handlerCheckPwd).Methods("POST")
	r.HandleFunc("/user/login", e.handlerUserLogin).Methods("POST")
}

func (e *UserServiceInterfaceTest) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["dataId"]
	if queryId == "" {
		response.WriteJson(w, response.FailMessageResp("用户ID不能为空"))
		return
	}
	t, err := e.userInter.GetUserById(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("获取用户数据失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}

// handlerCreateUser 处理创建用户
func (e *UserServiceInterfaceTest) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var user *types.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("解析用户新增参数失败"))
		return
	}

	defer r.Body.Close()

	password := utils.HashPassword(user.Password)
	//if err != nil {
	//	response.WriteJson(w, response.FailMessageResp("创建用户失败"))
	//	return
	//}

	user.Password = password
	creatUser, err := e.userInter.CreatUser(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("创建用户数据失败"))
		return
	}

	token, err := utils.CreatAndSerAuthCookie(creatUser.UserId, 7, w)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("返回用户token失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(token))
	return
}

// handlerCheckPwd 验证密码
func (e *UserServiceInterfaceTest) handlerCheckPwd(w http.ResponseWriter, r *http.Request) {
	tokenStr := interceptor.GetTokenFromRequest(r)
	token, err := interceptor.ValidateJWT(tokenStr)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("token不正确"))
		return
	}

	if !token.Valid {
		response.WriteJson(w, response.FailMessageResp("token验证失败"))
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)

	user, err := e.userInter.GetUserById(utils.ConvertToInt64(userId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("未查询到用户信息"))
		return
	}

	var mapPwd = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&mapPwd); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("获取密码失败"))
		return
	}
	defer r.Body.Close()

	password := mapPwd["password"]
	encodePwd := mapPwd["encodePwd"]

	log.Printf("原始密码===%s", password)
	log.Printf("加密密码(base64)===%s", encodePwd)

	encodePwdByte, _ := base64.StdEncoding.DecodeString(encodePwd)

	publicKeyName := "public.pem"
	privateKeyName := "private.pem"
	publicKey, _ := os.ReadFile(publicKeyName)
	privateKey, _ := os.ReadFile(privateKeyName)

	encrypt, _ := utils.RsaEncrypt(publicKey, []byte(password))
	encryptStr := base64.StdEncoding.EncodeToString(encrypt)
	log.Printf("RsaEncrypt===加密后的===%s", encryptStr) // 将这个传给前端进行解密
	//encryptByte, _ := base64.StdEncoding.DecodeString(encryptStr)
	//decrypt2, _ := utils.RsaDecrypt(privateKey, encryptByte)
	//log.Printf("RsaDecrypt解密后的1===%s", string(decrypt2))

	decrypt, _ := utils.RsaDecrypt(privateKey, encodePwdByte)
	log.Printf("RsaDecrypt解密后的2===%s", string(decrypt))

	//decodeStr := utils.RSADecode(encodePwd)
	//log.Printf("解密后的===%s", decodeStr)

	// var passwordHash = "$2a$10$3biXyR88nvb/yLtKQ9Ro0OCQJza2prdlGSqduDpBxfikVTNzGJdZ6"
	isSame := utils.ComparePassword(user.Password, string(decrypt))
	log.Printf("比较密码hash===%v", isSame)

	if !isSame {
		response.WriteJson(w, response.FailMessageResp("密码验证失败"))
		return
	}
	response.WriteJson(w, response.FailMessageResp("密码验证成功"))
	return
}

// handlerUserLogin 用户登录
func (e *UserServiceInterfaceTest) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	var user *types.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("解析用户登录参数失败"))
		return
	}

	defer r.Body.Close()

	loginUser, err := e.userInter.UserLogin(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("查询用户信息失败"))
		return
	}

	if loginUser == nil {
		response.WriteJson(w, response.FailMessageResp("当前用户信息不存在"))
		return
	}

	encodePwdByte, _ := base64.StdEncoding.DecodeString(user.Password)
	decryptPwdByte, _ := utils.RsaDecrypt([]byte(config.EnvConfig.PrivateKey), encodePwdByte)
	decryptPassword := string(decryptPwdByte)
	log.Printf("用户登录解析的密码===%s", decryptPassword)

	// 验证密码hash是否通过
	isSame := utils.ComparePassword(loginUser.Password, decryptPassword)
	log.Printf("用户登录密码hash===%v", isSame)

	if !isSame {
		response.WriteJson(w, response.FailMessageResp("密码输入错误"))
		return
	}

	token, err := utils.CreatAndSerAuthCookie(loginUser.UserId, 10, w)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("返回用户token失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(token))
	return
}
