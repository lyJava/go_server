package utils

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/types"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql/driver"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ConvertToInt64(str string) int64 {
	// 将字符串转换为 int64 类型
	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("字符串转换int64失败:", err)
	}
	return intVal
}

func ConvertToStr(val int64) string {
	return strconv.Itoa(int(val))
}

type LocalTime time.Time

func (t *LocalTime) String() string {
	// 如果时间 null 那么我们需要把返回的值进行修改
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s", time.Time(*t).Format("2006-01-02 15:04:05"))
}

func (t *LocalTime) IsZero() bool {
	return time.Time(*t).IsZero()
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	*t = LocalTime(t1)
	return err
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if &t == nil || t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil

}

//https://blog.csdn.net/tongweizhen/article/details/124190235

// FormatTime
//
// 参数
//
//	t (time.Time): 时间参数
//
// 返回
//
//	string yyyy-MM-dd HH:mm:ss
func FormatTime(t time.Time) string {
	return fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05"))
}

// HashPassword 返回密码hash
func HashPassword(pwd string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
		return ""
	}
	return string(password)
}

func CreateJWT(user *types.User, days int64, secret []byte) (string, error) {
	// "userId":     strconv.Itoa(int(userId)),
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"userId":     strconv.Itoa(int(userId)),
	//	"expireTime": time.Now().Add(time.Hour * 24 * time.Duration(days)).Unix(), // 设置有效期为20天
	//})
	//tokenStr, err := token.SignedString(secret)
	c := types.MyClaims{
		UserId:   strconv.FormatInt(user.UserId, 10),
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(days)).Unix(), // 过期时间
			Issuer:    "el-admin",                                                  // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// CreatAndSerAuthCookie 创建并设置token到cookie
// 参数
//
//		userId: 用户ID
//	 days: token有效天数
//	 w:	http返回
func CreatAndSerAuthCookie(user *types.User, days int64, w http.ResponseWriter) (string, error) {
	secret := []byte(config.EnvConfig.JWTSecret)
	token, err := CreateJWT(user, days, secret)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}

// ComparePassword 比较解密后的密码和哈希密码是否匹配
// 参数
//
//	hashedPassword(string):	用户的密码(hash过的)
//	passwordFrontend(string): 前端解密后的密码
func ComparePassword(hashedPassword, passwordFrontend string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordFrontend))
	if err != nil {
		log.Panicln(err.Error())
		return false
	}
	return err == nil
}

// RSAEncode 加密
func RSAEncode(val string) string {
	encodeByte := []byte(val)
	publicKey := config.EnvConfig.PublicKey
	//log.Println("公钥字符串===" + publicKey)
	publicBlock, _ := pem.Decode([]byte(publicKey))

	publicKeyVal, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		log.Panicln("读取公钥失败", err)
		return ""
	}

	pub := publicKeyVal.(*rsa.PublicKey)
	encryptTextByte, err := rsa.EncryptPKCS1v15(rand.Reader, pub, encodeByte)
	if err != nil {
		log.Panicln("读取公钥失败", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(encryptTextByte)
}

// RSADecode 解密
func RSADecode(val string) string {
	// 解码
	cipherText, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		fmt.Println("Error decoding:", err)
	}
	log.Printf("base64解码===%s", string(cipherText))

	privateKeyStr := config.EnvConfig.PrivateKey
	privateBlock, _ := pem.Decode([]byte(privateKeyStr))

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		log.Panicln("读取私钥失败1", err.Error())
		return ""
	}
	decryptText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		log.Panicln("解密失败", err.Error())
		return ""
	}
	return string(decryptText)
}

// RsaEncrypt 加密
func RsaEncrypt(publicKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RsaDecrypt 解密
func RsaDecrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
