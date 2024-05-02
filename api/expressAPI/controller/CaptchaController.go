package controller

import (
	"apiProject/api/response"
	"apiProject/api/utils"
	"github.com/gorilla/mux"
	"image/png"
	"net/http"
)

type CaptchaController struct {
}

func CaptchaControllerTest() *CaptchaController {
	return &CaptchaController{}
}

func (*CaptchaController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/captcha", CaptchaHandler).Methods("GET")
}

func CaptchaHandler(w http.ResponseWriter, _ *http.Request) {
	code := utils.RandomCode(4, "")
	// cache.Set("captcha_"+code, code, 30)
	img := utils.CreateImage(code)
	err := png.Encode(w, img)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("验证码生成失败"))
		return
	}
}
