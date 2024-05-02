package controller

import (
	"apiProject/api/response"
	"apiProject/api/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"image/png"
	"log"
	"net/http"
	"strconv"
)

type CaptchaController struct {
}

type CaptchaData struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type CaptchaValidate struct {
	Id   string
	Code string
}

func CaptchaControllerTest() *CaptchaController {
	return &CaptchaController{}
}

func (*CaptchaController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/captcha", CaptchaCustomerHandler).Methods("GET")
	r.HandleFunc("/captcha/create", CaptchaHandlerCreate).Methods("GET")
	r.HandleFunc("/captcha/validate", CaptchaHandlerValidate).Methods("POST")
}

// CaptchaCustomerHandler 自定义验证码
func CaptchaCustomerHandler(w http.ResponseWriter, _ *http.Request) {
	code := utils.RandomCode(4, "")
	// cache.Set("captcha_"+code, code, 30)
	img := utils.CreateImage(code)
	// 将图像编码为字节切片
	var imgBuffer bytes.Buffer
	err := png.Encode(&imgBuffer, img)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 将字节切片转换为 base64 编码的字符串
	imgBase64Str := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBuffer.Bytes())
	log.Println(imgBase64Str)
	// 这里设置Content-Length长度，前端下载才能获取进度
	w.Header().Set("Content-Length", strconv.Itoa(imgBuffer.Len()))

	// 设置响应头
	//w.Header().Set("Content-Type", "image/png")
	// 将 base64 编码的字符串写入响应体
	//http.ServeContent(w, r, code+".png", time.Time{}, bytes.NewReader(imgBuffer.Bytes()))

	err = png.Encode(w, img)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("验证码生成失败"))
		return
	}
}

// CaptchaHandlerCreate 使用github.com/dchest/captcha生成验证码
func CaptchaHandlerCreate(w http.ResponseWriter, r *http.Request) {
	//current, _ := user.Current()
	//log.Println("当前用户id：", current.Uid)
	//log.Println("当前用户的username：", current.Username)
	//log.Println("当前用户的homedir：", current.HomeDir)
	//log.Println("当前用户的gid：", current.Gid)
	imageId := captcha.NewLen(4)
	log.Println("验证码ID:", imageId)
	var content bytes.Buffer
	err := captcha.WriteImage(&content, imageId, 120, 60)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}
	// 这里是直接返回图片
	//http.ServeContent(w, r, codeId+".png", time.Time{}, bytes.NewReader(content.Bytes()))
	// 这里使用map方式
	//imageData := make(map[string]string)
	//imageData["id"] = imageId
	//imageData["data"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(content.Bytes())
	imageData := &CaptchaData{
		Id:   imageId,
		Data: "data:image/png;base64," + base64.StdEncoding.EncodeToString(content.Bytes()),
	}
	response.WriteJson(w, response.OkDataResp(imageData))
	return
}

// CaptchaHandlerValidate 校验验证码
func CaptchaHandlerValidate(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Query().Get("id")
	//code := r.URL.Query().Get("code")

	//imageData := make(map[string]string)
	var imageData CaptchaValidate
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		log.Println(err.Error())
		response.WriteJson(w, response.FailMessageResp("解析用户新增参数失败"))
		return
	}

	// 读取请求的 Body 数据
	/*body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// 定义用于存储解析后数据的 map
	var dataMap map[string]interface{}

	// 使用 mapstructure 解析 Body 数据到一个 map 中
	if err := json.Unmarshal(body, &dataMap); err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// 使用 mapstructure 包的 Decode 函数将 map 解析为结构体
	if err := mapstructure.Decode(dataMap, &imageData); err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}*/

	defer r.Body.Close()

	if imageData.Id == "" || imageData.Code == "" {
		response.WriteJson(w, response.FailMessageResp("参数不能为空"))
		return
	}

	//result := captcha.VerifyString(imageData, code)
	// result := captcha.VerifyString(imageData["id"], imageData["code"])
	// 三元表达式写法???
	/*response.WriteJson(w, func() response.Response {
		if result {
			return response.OkMessageResp("验证成功")
		}
		return response.FailMessageResp("验证失败")
	}())*/
	// result := captcha.VerifyString(imageData["id"], imageData["code"])
	result := captcha.VerifyString(imageData.Id, imageData.Code)
	if !result {
		response.WriteJson(w, response.FailMessageResp("验证失败"))
		return
	}
	response.WriteJson(w, response.OkMessageResp("验证成功"))
}
