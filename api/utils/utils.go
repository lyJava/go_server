package utils

import (
	"apiProject/api/expressAPI/config"
	configure "apiProject/api/expressAPI/types/config"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	_ "math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// ConvertToInt64 字符串转换为int64
func ConvertToInt64(str string) int64 {

	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Printf("字符串转换int64失败:%+v", err)
	}
	return intVal
}

// ConvertToInt 字符串转换为int
func ConvertToInt(str string) int {
	// 将字符串转换为 int64 类型
	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Printf("字符串转换int64失败:%+v", err)
	}
	return int(intVal)
}

// ConvertInt64ToStr int64转换为字符串
func ConvertInt64ToStr(val int64) string {
	return strconv.Itoa(int(val))
}

// ConvertIntToStr int转换为字符串
func ConvertIntToStr(val int) string {
	return strconv.Itoa(val)
}

// ConvertStrToBool 字符串转换bool
func ConvertStrToBool(val string) (bool, error) {
	parseBool, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("string convert to bool error:%+v", err.Error())
		return false, err
	}
	return parseBool, nil
}

// ConvertBoolToStr bool转换字符串
func ConvertBoolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

type LocalTime time.Time

func (t *LocalTime) String() string {
	// 如果时间 null 那么我们需要把返回的值进行修改
	if t == nil || t.IsZero() {
		return ""
	}
	return time.Time(*t).Format("2006-01-02 15:04:05")
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
	return fmt.Errorf("can not convert %+v to timestamp", v)
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
	/*
		在这个函数中，t 是 LocalTime 类型的接收者，它是一个自定义类型。在 Go 语言中，对于结构体或自定义类型的方法，即使接收者是指针类型，
		如果使用的是非指针类型的值来调用该方法，Go 语言也会自动进行值到指针的转换。因此，在 MarshalJSON 方法中，
		t 的类型是 LocalTime，即使 t 为 nil，也不等于 nil。因此，您无需检查 &t == nil 来判断 t 是否为零值或空值
	*/
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if t.IsZero() {
		return []byte(""), nil
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
	return t.Format("2006-01-02 15:04:05")
}

// HashPassword 返回密码hash
//
// 参数
//
//	pwd (string): 密码
func HashPassword(pwd string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
		return ""
	}
	return string(password)
}

// CreateJWT 创建JWT凭证
//
// 参数
//
//	user (*domain.User): 用户信息
//	days (int64): 有效天数
//	secret ([]byte): 加密字符串
func CreateJWT(user *domain.User, days int64, secret []byte) (string, error) {
	// "userId":     strconv.Itoa(int(userId)),
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"userId":     strconv.Itoa(int(userId)),
	//	"expireTime": time.Now().Add(time.Hour * 24 * time.Duration(days)).Unix(), // 设置有效期为20天
	//})
	//tokenStr, err := token.SignedString(secret)
	claims := configure.MyClaims{
		UserId:   strconv.FormatInt(user.UserId, 10),
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(days)).Unix(), // 过期时间
			Issuer:    "el-admin",                                                  // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		fmt.Printf("创建JWT凭证错误===%+v", err)
		return "", err
	}
	return tokenStr, nil
}

// CreateToken 创建token
// 参数
//
//	userId: 用户信息
//	days: 有效天数
func CreateToken(user *domain.User, days int64) (string, error) {
	secret := []byte(config.EnvConfig.JWTSecret)
	token, err := CreateJWT(user, days, secret)
	if err != nil {
		fmt.Printf("创建凭证错误===%+v", err)
		return "", errors.New("创建token失败")
	}
	return token, nil
}

// CreateAndSetAuthCookie 创建并设置token到cookie
// 参数
//
//	userId: 用户ID
//	days: token有效天数
//	w:	http返回
func CreateAndSetAuthCookie(w http.ResponseWriter, user *domain.User, days int64) (string, error) {
	secret := []byte(config.EnvConfig.JWTSecret)
	token, err := CreateJWT(user, days, secret)
	if err != nil {
		log.Printf("create and set token error === %+v", err)
		return "", errors.New("创建并设置token失败")
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
		fmt.Printf("Error decoding:%+v", err)
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
//
// 参数
//
//	publicKeyByte ([]byte): 公钥byte
//	origData ([]byte): 加密byte
func RsaEncrypt(publicKeyByte []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKeyByte)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface any
	var err error

	switch block.Type {
	case "RSA PUBLIC KEY":
		pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	case "PUBLIC KEY":
		pubInterface, err = x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		log.Printf("未知的公钥类型===%s, %+v", block.Type, err)
		err = fmt.Errorf("未知的公钥类型: %s", block.Type)
	}

	if err != nil {
		log.Printf("RSA公钥解析错误===%+v", err)
		return nil, fmt.Errorf("解析公钥失败: %+v", err)
	}

	publicKey := pubInterface.(*rsa.PublicKey)
	// 加密数据
	encryptDataByte, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, origData)
	if err != nil {
		log.Printf("RSA公钥加密错误===%+v", err)
		return nil, errors.New("加密失败")
	}
	return encryptDataByte, nil
}

// RsaDecrypt 解密
//
// 参数
//
//	privateKeyByte ([]byte): 密钥byte
//	cipherText ([]byte): 解密byte
func RsaDecrypt(privateKeyByte []byte, cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKeyByte)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var privateInterface interface{}
	var err error

	if block.Type == "RSA PRIVATE KEY" {
		privateInterface, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		privateInterface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	} else {
		log.Printf("未知的私钥类型===%+v", err)
		err = errors.New("解密失败")
	}

	if err != nil {
		log.Printf("RSA私钥解析错误===%+v", err)
		return nil, err
	}

	privateKey := privateInterface.(*rsa.PrivateKey)
	decryptDataByte, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		log.Printf("RSA私钥解密错误===%+v", err)
		return nil, errors.New("解密失败")
	}
	return decryptDataByte, nil
}

// GetKeyByteByPath 从文件中读取密钥
//
// 参数
//
//	pemPath (string): 密钥路径
func GetKeyByteByPath(pemPath string) ([]byte, error) {
	fileByte, err := os.ReadFile(pemPath)
	if err != nil {
		log.Printf("read key byte from %s failed %+v", pemPath, err)
		return nil, errors.New("读取密钥文件错误")
	}
	log.Printf("read key byte from %s success", pemPath)
	return fileByte, nil
}

// GetPrefix 获取前缀
//
// 参数
//
//	str (string): 文件名
func GetPrefix(str string) string {
	var prefix = ""
	lastDotIndex := strings.LastIndex(str, ".")
	if lastDotIndex != -1 {
		fileName := str[:lastDotIndex]
		prefix = url.QueryEscape(fileName)
	} else {
		return str
	}
	return prefix
}

// GetSuffix 获取后缀
//
// 参数
//
//	str (string): 文件名
func GetSuffix(str string) string {
	var suffix = ""
	lastDotIndex := strings.LastIndex(str, ".")
	if lastDotIndex != -1 {
		fileName := str[lastDotIndex+1:]
		suffix = url.QueryEscape(fileName)
	} else {
		return str
	}
	return suffix
}

// CreateZipBuffer 创建ZIP文件到缓冲区
func CreateZipBuffer(files []string) (*bytes.Buffer, error) {
	// 创建 ZIP 文件
	zipBuffer := new(bytes.Buffer)
	zw := zip.NewWriter(zipBuffer)

	// 逐个文件添加到 ZIP 文件中
	for _, filePath := range files {
		// 打开要添加到 ZIP 文件的文件
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("打开zip文件错误===%+v", err)
			return nil, err
		}
		defer file.Close()

		// 将文件添加到 ZIP 文件中
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			fmt.Printf("获取zip文件头部信息错误===%+v", err)
			return nil, err
		}

		// 设置zip文件名
		header.Name = filepath.Base(filePath)

		// 创建zip文件头
		writer, err := zw.CreateHeader(header)
		if err != nil {
			fmt.Printf("创建zip头错误===%+v", err)
			return nil, err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			fmt.Printf("项zip文件中复制错误===%+v", err)
			return nil, err
		}
	}
	// 关闭 ZIP Writer
	err := zw.Close()
	if err != nil {
		return nil, err
	}

	return zipBuffer, nil
}

// CreateZipTemp 创建zip到临时文件夹
//
// 参数
//
//	fileNames 文件路径加名称
//	tempPathPattern 临时文件路径形式
func CreateZipTemp(fileNames []string, tempPathPattern string) (*os.File, error) {
	// 创建临时文件, tempPathPattern为"temp-zip-*.zip"表示临时文件为temp-zip-xxxxx.zip格式
	tmpFile, err := os.CreateTemp("", tempPathPattern)
	if err != nil {
		log.Printf("创建临时文件错误===%+v", err)
		return &os.File{}, err
	}

	log.Printf("临时文件:%s 创建成功", tmpFile.Name())
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("临时文件删除失败", err.Error())
		}
		log.Printf("临时文件:%s 删除成功", tmpFile.Name())
	}(tmpFile.Name()) // 在函数退出时删除临时文件

	// 创建 ZIP 编写器
	zw := zip.NewWriter(tmpFile)

	// 遍历文件列表
	for _, fileName := range fileNames {
		// 打开文件
		file, err := OpenFile(fileName)
		if err != nil {
			fmt.Printf("遍历文件错误===%+v", err)
			return &os.File{}, err
		}
		defer file.Close()

		fileInfo, err := GetFileInfo(file)
		if err != nil {
			fmt.Printf("获取文件信息错误===%+v", err)
			return &os.File{}, err
		}

		log.Printf("文件路径:%s", fileName)

		// 将文件添加到 ZIP 文件中
		fileInZip, err := zw.Create(fileInfo.Name())
		if err != nil {
			log.Printf("创建ZIP文件内部错误===%+v", err)
			return &os.File{}, err
		}

		// 将文件内容复制到 ZIP 文件中，该方式适合小文件与非高并发情况
		/*_, err = io.Copy(fileInZip, file)
		if err != nil {
			log.Println("复制文件内容到ZIP文件内部异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("复制文件内容到ZIP文件内部失败"))
			return
		}*/

		// 将文件内容逐块复制到 ZIP 文件中
		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				fmt.Printf("读取文件错误===%+v", err)
				break
			}
			if err != nil {
				log.Printf("读取文件异常===%+v", err)
				return &os.File{}, err
			}

			_, err = fileInZip.Write(buf[:n])
			if err != nil {
				log.Printf("写入ZIP文件内部异常==%+v", err)
				return &os.File{}, err
			}
		}
	}

	// 关闭 ZIP 编写器
	if err := zw.Close(); err != nil {
		log.Printf("关闭ZIP编写器异常===%+v", err)
		return &os.File{}, err
	}
	return tmpFile, nil
}

// OpenFile 打开文件
//
// 参数
//
//	fileName (string): 文件名
func OpenFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("打开文件异常===%+v", err)
		return &os.File{}, err
	}
	return file, nil
}

// GetFileInfo 获取文件信息
//
// 参数
//
//	file (*os.File): 文件对象
func GetFileInfo(file *os.File) (os.FileInfo, error) {
	// 获取临时文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("获取文件信息异常===%+v", err)
		return nil, err
	}

	contentLength := strconv.FormatInt(fileInfo.Size(), 10)
	log.Println("文件长度", contentLength)
	zipDownloadName := filepath.Base(fileInfo.Name())
	log.Println("文件名称", zipDownloadName)
	return fileInfo, nil
}

// CloseBodyError 处理请求体关闭可能存在的异常
//
// 参数
//
//	message (string): 消息
//	w (http.ResponseWriter): 返回
//	r (*http.Request): 请求
func CloseBodyError(message string, w http.ResponseWriter, r *http.Request) {
	if err := r.Body.Close(); err != nil {
		log.Printf("%s关闭出现错误===%+v", message, err)
		response.WriteJson(w, response.FailMessageResp(message+"关闭失败"))
		return
	}
	log.Println(message + "关闭成功")
}

// CamelToSnakeCase 将驼峰形式的字符串转换为下划线形式
func CamelToSnakeCase(s string) string {
	var buffer bytes.Buffer
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// ConvertOrder 将ant design 表格排序转换为sql排序关键字
func ConvertOrder(key string) string {
	orderMap := map[string]string{
		"ascend":  "ASC",
		"descend": "DESC",
	}
	return orderMap[key]
}

// HandlerColumnOrder 处理列排序
func HandlerColumnOrder(column, order string) (string, string) {
	if column == "" && order == "" {
		return "", ""
	}
	return CamelToSnakeCase(column), ConvertOrder(order)
}

// DownloadFile 下载文件
//
// 参数
//
//	fileNamePath (string): 文件名称路径
//	w (http.ResponseWriter): 返回
//	r (*http.Request): 请求
func DownloadFile(fileNamePath string, w http.ResponseWriter, r *http.Request) {
	file, err := OpenFile(fileNamePath)
	if err != nil {
		log.Printf("下载文件出现错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	defer file.Close()

	fileInfo, err := GetFileInfo(file)
	if err != nil {
		log.Printf("获取文件信息异常===%+v", err)
		response.WriteJson(w, response.FailMessageResp("获取文件信息失败"))
		return
	}

	var fileMineType string

	// 获取文件扩展名
	suffix := strings.ToLower(filepath.Ext(fileInfo.Name()))
	// 检查文件扩展名是否为.xlsx 或.xls
	if suffix == ".xlsx" {
		// excel 2007格式
		fileMineType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	} else if suffix == ".xls" {
		fileMineType = "application/vnd.ms-excel"
	} else {
		// 读取文件的前512字节
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			log.Printf("读取文件异常===%+v", err)
			response.WriteJson(w, response.FailMessageResp("读取文件失败"))
			return
		}
		fileMineType = http.DetectContentType(buffer)
	}
	log.Println("当前文件MIME:", fileMineType)

	w.Header().Set("Content-Type", fileMineType)
	// 对文件名进行编码处理，避免在firefox或postman中请求下载变成response.bin
	w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.QueryEscape(fileNamePath))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	http.ServeFile(w, r, fileNamePath)
}

// GetStrFromMap 从map中获取value
//
// 参数
//
//	m (map[string]interface{}): map
//	key (string): map的key
func GetStrFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// ConvertInt64ListToStrList int64数组转换字符串数组
func ConvertInt64ListToStrList(list []int64) []string {
	if len(list) == 0 {
		return nil
	}
	strList := make([]string, 0, len(list))
	for _, i := range list {
		strList = append(strList, ConvertInt64ToStr(i))
	}
	return strList
}

// ConvertStrToIntList 字符串转换int64切片
func ConvertStrToIntList(str string) []int64 {
	if str == "" || len(str) == 0 {
		return nil
	}
	strList := strings.Split(str, ",")
	int64List := make([]int64, 0, len(strList))
	for _, s := range strList {
		i := strings.TrimSpace(s)
		int64List = append(int64List, ConvertToInt64(i))
	}
	return int64List
}

// ConvertStrListToIntList 字符串切片转换int64切片
//
// 参数
//
//	list ([]string): 字符串切片
//
// 返回
//
//	[]int64: int64切片
func ConvertStrListToIntList(list []string) []int64 {
	if len(list) == 0 {
		return nil
	}
	int64List := make([]int64, 0, len(list))
	for _, s := range list {
		i := strings.TrimSpace(s)
		int64List = append(int64List, ConvertToInt64(i))
	}
	return int64List
}

// GeneratePlaceholders 生成占位符字符串，例如: $1, $2, $3
// 参数
//
//	n (int): 生成的数量
func GeneratePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	var placeholders []string
	for i := 1; i <= n; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}
	return strings.Join(placeholders, ", ")
}

// GenerateArgs 构建arg切片
// 参数
//
//	ids ([]interface{}): 泛型切片
//
// 返回
//
//	[]interface{}: 泛型切片
func GenerateArgs(ids []interface{}) []interface{} {
	args := make([]interface{}, len(ids))
	/* for i, id := range ids {
		args[i] = id
	} */
	copy(args, ids)
	return args
}

// BuildArgsWithBrackets 将参数使用括号构建
//
// 参数
//
//	args ([]any): 参数
func BuildArgsWithBrackets(args []any) string {
	/*output := "("
	for i, arg := range args {
		output += fmt.Sprintf("%v", arg)
		if i < len(args)-1 {
			output += ", "
		}
	}
	output += ")"
	return output*/
	var builder strings.Builder
	builder.WriteString("(")
	for i, arg := range args {
		builder.WriteString(fmt.Sprintf("%v", arg))
		if i < len(args)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")
	return builder.String()
}

// TimeForHuman 格式化日期
//
// 参数
//
//	timeValue (int64): 10位时间戳
func TimeForHuman(timeValue int64) string {
	SECOND := int64(1)
	MINUTE := SECOND * 60
	HOUR := MINUTE * 60
	DAY := HOUR * 24
	DAY8 := DAY * 8

	nowTime := time.Now().Unix()
	diffTime := nowTime - timeValue

	log.Printf("传入==%v,当前==%v,时间差==%v", timeValue, nowTime, diffTime)
	if diffTime <= MINUTE {
		return "刚刚"
	} else if diffTime < HOUR {
		return fmt.Sprintf("%d分钟前", int(diffTime/MINUTE))
	} else if diffTime <= DAY {
		return fmt.Sprintf("%d小时前", int(diffTime/HOUR))
	} else if diffTime <= DAY8 {
		return fmt.Sprintf("%d天前", int(diffTime/DAY))
	} else {
		return time.Unix(timeValue, 0).Format("2006-01-02")
	}
}

// ShowJsonFormat 返回格式化json
func ToJsonFormat(data any) string {
	if data == "" {
		return ""
	}
	marshal, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		zap.L().Sugar().Errorf("数据格式化json错误===%+v", err)
		return ""
	}
	return string(marshal)
}

// ForcedToJsonFormat 强制json格式化输出
func ForcedToJsonFormat(b []byte) string {

	// 检查输入数据是否为空
	if len(b) == 0 {
		log.Printf("Error: input is empty")
		return ""
	}

	var result map[string]interface{}
	// 将JSON字符串解码到result变量中
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Printf("Error occurred during unmarshaling. Error: %+v", err)
	}

	// 格式化输出JSON字符串
	formattedJson, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error occurred during marshaling. Error:%+v", err)
	}
	return string(formattedJson)
}

// QueryParamToJson 将url的query参数转换json
func QueryParamToJson(queryParam string) string {
	// 解析查询字符串
	values, _ := url.ParseQuery(queryParam)
	if len(values) == 0 {
		return ""
	}

	result := make(map[string]any)

	for key, value := range values {
		log.Printf("key===%s,value===%v", key, value)
		if len(value) > 0 {
			result[key] = value[0]
		}
	}

	return ToJsonFormat(result)
}
