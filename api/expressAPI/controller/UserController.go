package controller

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/interceptor"
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserServiceInterface // 用户服务接口
}

// UserControllerInit 用户控制器初始化
func UserControllerInit(u service.UserServiceInterface) *UserController {
	return &UserController{userService: u}
}

// RegisterRoutes 注册快递服务请求路由
func (u *UserController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/{dataId}", u.handlerGetUser).Methods("GET")
	r.HandleFunc("/user/create", u.handlerCreateUser).Methods("POST")
	r.HandleFunc("/user/checkPwd", u.handlerCheckPwd).Methods("POST")
	r.HandleFunc("/user/login", u.handlerUserLogin).Methods("POST")
}

func (u *UserController) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["dataId"]
	if queryId == "" {
		response.WriteJson(w, response.FailMessageResp("用户ID不能为空"))
		return
	}
	t, err := u.userService.GetUserById(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("获取用户数据失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}

// handlerCreateUser 处理创建用户
func (u *UserController) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var user *domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		response.WriteJson(w, response.FailMessageResp("解析用户新增参数失败"))
		return
	}

	utils.CloseBodyError("用户新增失败", w, r)

	user.Password = utils.HashPassword(user.Password)
	createUser, err := u.userService.CreateUser(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	// 将json格式化输出
	marshal, _ := json.MarshalIndent(createUser, "", "    ")
	log.Printf("用户新增===\n%s", marshal)

	token, err := utils.CreateToken(createUser, 7)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	// 创建map返回数据
	dataMap := map[string]interface{}{
		"token": token,
		"user":  createUser,
	}
	response.WriteJson(w, response.OkDataResp(dataMap))
}

// handlerCheckPwd 验证密码
func (u *UserController) handlerCheckPwd(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := interceptor.GetTokenFromRequest(r)
	if err != nil {
		zap.L().Sugar().Errorf("验证密码获取token错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}
	token, err := interceptor.ValidateJWT(tokenStr)
	if err != nil {
		zap.L().Sugar().Errorf("验证密码验证token有效性错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if !token.Valid {
		response.WriteJson(w, response.FailMessageResp("token验证失败"))
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)

	user, err := u.userService.GetUserById(utils.ConvertToInt64(userId))
	if err != nil {
		zap.L().Sugar().Errorf("验证密码验证获取用户信息错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("未查询到用户信息"))
		return
	}

	var pwdMap = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&pwdMap); err != nil {
		log.Printf("验证密码验证获取密码参数错误===%+v", err.Error())
		zap.L().Sugar().Errorf("验证密码验证获取密码参数错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("获取密码失败"))
		return
	}

	defer utils.CloseBodyError("验证密码失败", w, r)

	if len(pwdMap) != 2 {
		response.WriteJson(w, response.FailMessageResp("参数不正确"))
		return
	}
	password := pwdMap["password"]
	if password == "" {
		// 实际使用情况原始密码一般不需要传递
		response.WriteJson(w, response.FailMessageResp("原始密码不能为空"))
		return
	}

	encodePwd := pwdMap["encodePwd"]
	if encodePwd == "" {
		response.WriteJson(w, response.FailMessageResp("加密密码不能为空"))
		return
	}

	log.Printf("前端原始密码===%s", password)
	log.Printf("前端加密密码(base64)===%s", encodePwd)

	encodePwdByte, _ := base64.StdEncoding.DecodeString(encodePwd)

	/*publicKey, err := utils.GetKeyByteByPath("public.pem")
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}*/
	privateKey, err := utils.GetKeyByteByPath("private.pem")
	if err != nil {
		zap.L().Sugar().Errorf("验证密码验证获取私钥错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	/*encrypt, _ := utils.RsaEncrypt(publicKey, []byte(password))
	encryptStr := base64.StdEncoding.EncodeToString(encrypt)
	log.Printf("RsaEncrypt===加密后的===%s", encryptStr) */ // 将这个传给前端进行解密

	//encryptByte, _ := base64.StdEncoding.DecodeString(encryptStr)
	//decrypt2, _ := utils.RsaDecrypt(privateKey, encryptByte)
	//log.Printf("RsaDecrypt解密后的1===%s", string(decrypt2))

	decrypt, err := utils.RsaDecrypt(privateKey, encodePwdByte)
	log.Printf("RsaDecrypt解密后的2===%s", string(decrypt))
	if err != nil {
		zap.L().Sugar().Errorf("验证密码验证rsa解密错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	//decodeStr := utils.RSADecode(encodePwd)
	//log.Printf("解密后的===%s", decodeStr)

	// var passwordHash = "$2a$10$3biXyR88nvb/yLtKQ9Ro0OCQJza2prdlGSqduDpBxfikVTNzGJdZ6"
	isSame := utils.ComparePassword(user.Password, string(decrypt))
	log.Printf("比较密码hash===%v", isSame)

	if !isSame {
		response.WriteJson(w, response.FailMessageResp("密码验证失败"))
		return
	}
	
	response.WriteJson(w, response.OkMessageResp("密码验证成功"))
}

// handlerUserLogin 用户登录
func (u *UserController) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	var user *domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("解析用户登录参数错误===%+v", err)
		zap.L().Sugar().Errorf("解析用户登录参数错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("解析用户登录参数失败"))
		return
	}
	utils.CloseBodyError("用户登录失败", w, r)

	if user.Username == "" {
		response.WriteJson(w, response.FailMessageResp("用户名不能为空"))
		return
	}

	if user.Password == "" {
		response.WriteJson(w, response.FailMessageResp("密码不能为空"))
		return
	}

	loginUser, err := u.userService.UserLogin(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("查询用户信息失败"))
		return
	}

	if loginUser == nil {
		response.WriteJson(w, response.FailMessageResp("当前用户信息不存在"))
		return
	}

	encodePwdByte, _ := base64.StdEncoding.DecodeString(user.Password)
	decryptPwdByte, err := utils.RsaDecrypt([]byte(config.EnvConfig.PrivateKey), encodePwdByte)
	if err != nil {
		zap.L().Sugar().Errorf("解析用户登录使用rsa解密错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("解密失败"))
		return
	}
	decryptPassword := string(decryptPwdByte)
	log.Printf("用户登录解析的密码===%s", decryptPassword)

	// 验证密码hash是否通过
	isSame := utils.ComparePassword(loginUser.Password, decryptPassword)
	log.Printf("用户登录密码hash===%v", isSame)

	if !isSame {
		response.WriteJson(w, response.FailMessageResp("密码输入错误"))
		return
	}

	token, err := utils.CreateAndSetAuthCookie(w, loginUser, 10)
	if err != nil {
		zap.L().Sugar().Errorf("用户凭证及cookie设置错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkDataResp(token))
}
