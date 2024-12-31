package utils

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 验证码字符集
var charset = "0123456789ABCDEFGHJKPLMNPQRSTUVWXYZ"

// GenerateCode 生成指定长度的随机验证码字符串
//
// 参数
//
//	length (int): 验证码长度
func GenerateCode(length int) string {
	rand.NewSource(time.Now().UnixNano())
	// 生成随机验证码字符串
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// RandomCode 随机数字或者字符串验证码
//
// 参数
//
//	size (int): 验证码个数
//	typeStr (string): 验证码类型(空字符串或者math默认为随机数字，其他为随机字符串)
func RandomCode(size int, typeStr string) string {
	code := ""
	rand.NewSource(time.Now().UnixNano())
	if typeStr == "" || typeStr == "math" {
		for i := 0; i < size; i++ {
			code += fmt.Sprintf("%d", rand.Intn(10))
		}
	}
	if typeStr != "" && typeStr != "math" {
		code = GenerateCode(size)
	}
	return code
}

func RandomCodeStr(length int, codeStr string) string {
	rand.NewSource(time.Now().UnixNano())
	// 生成随机验证码字符串
	code := make([]byte, length)
	for i := range code {
		code[i] = codeStr[rand.Intn(len(codeStr))]
	}
	return string(code)
}

// GenerateImage 生成验证码图片
//
// 参数
//
//	code (string): 验证码
func GenerateImage(code string) image.Image {
	// 图片尺寸
	width := 120
	height := 40

	// 创建空白图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景色设置为灰白色
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Pt(0, 0), draw.Src)

	// 使用默认字体
	face := basicfont.Face7x13

	// 绘制验证码文本
	point := fixed.Point26_6{X: 10, Y: 20}
	// 计算文本总宽度
	totalWidth := fixed.I(len(code)) * fixed.I(20)
	// 计算起始位置
	startX := fixed.I(width)/2 - totalWidth/2
	point.X = startX
	for _, ch := range code {
		drawText(img, face, point, string(ch), color.Black) // 设置文本颜色为黑色
		point.X += 20                                       // 每个字符之间的间距
	}

	return img
}

// drawText 在指定位置绘制文本
func drawText(img draw.Image, face font.Face, point fixed.Point26_6, text string, color color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
}

func CaptchaHandler(w http.ResponseWriter, r *http.Request) {
	// 生成验证码
	code := GenerateCode(4)

	// 生成验证码图片
	img := GenerateImage(code)

	// 设置响应头
	w.Header().Set("Content-Type", "image/png")

	// 将图片编码为 PNG 格式并发送给客户端
	err := png.Encode(w, img)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
	fmt.Println("验证码:", code)
}

// drawString 在指定位置绘制文本
func drawString(img draw.Image, str string, point image.Point, color color.Color) {
	for i := 0; i < len(str); i++ {
		drawRect(img, point, color)
		point.X += 10 // 每个字符之间的间距
	}
}

// drawRect 在指定位置绘制矩形代替字符
func drawRect(img draw.Image, point image.Point, color color.Color) {
	// 矩形宽度和高度
	width := 10
	height := 20

	// 绘制矩形
	rect := image.Rect(point.X, point.Y, point.X+width, point.Y+height)
	draw.Draw(img, rect, &image.Uniform{C: color}, image.Pt(0, 0), draw.Src)
}

//https://blog.csdn.net/m0_46198325/article/details/134913801

func CreateImage(code string) image.Image {

	/* fontPath, err := GetFontPath("/api/expressAPI", "/Arial Unicode.ttf")
	if err != nil {
		log.Println("获取字体文件路径异常")
		return nil
	} */

	//字体设置
	fontFile, err := os.ReadFile(fontAbsolutePath)
	if err != nil {
		log.Printf("open file failed:%+v", err)
		return nil
	}
	parse, err := truetype.Parse(fontFile)
	if err != nil {
		log.Printf("load parse failed:%+v", err)
		return nil
	}

	rand.NewSource(time.Now().UnixNano())
	w := 140
	h := 50
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Pt(0, 0), draw.Src)

	fontSize := 32
	fontDPI := 72.0

	dc := freetype.NewContext()
	dc.SetDPI(fontDPI)
	dc.SetFont(parse)
	dc.SetFontSize(float64(fontSize))
	dc.SetClip(img.Bounds())
	dc.SetDst(img)
	dc.SetSrc(&image.Uniform{C: color.RGBA{A: 255}})

	textWidthOld := getTextWidth(code, parse, fontSize)
	startX := (w-textWidthOld)/2 - fontSize - textWidthOld*6
	pt := freetype.Pt(startX, 35)

	for _, ch := range code {
		// 绘制文本
		textRandomColor(dc)
		_, err := dc.DrawString(string(ch), pt)
		if err != nil {
			log.Printf("Draw string failed:%+v", err)
			return nil
		}
		// 更新下一个字符的位置
		pt.X += dc.PointToFixed(float64(fontSize * 5 / 6))
	}

	/*pt := freetype.Pt(10, 60)
	for _, ch := range code {
		_, err := dc.DrawString(string(ch), pt)
		if err != nil {
			log.Println("Draw string failed")
		}
		pt.X += dc.PointToFixed(float64(fontSize * 5 / 7))
	}*/

	// 计算文本宽度
	/*textWidthOld := getTextWidth(code, parse, fontSize)
	log.Println("textWidth", textWidthOld)
	log.Println("fontSize", fontSize)
	startX := (w-textWidthOld)/2 - fontSize - textWidthOld*2 // 计算起始绘制位置
	log.Println("开始x", startX)
	pt := freetype.Pt(startX, 35)
	for _, ch := range code {
		textRandomColor(dc)
		_, err := dc.DrawString(string(ch), pt)
		if err != nil {
			log.Println("Draw string failed:", err)
			return nil
		}
		pt.X += dc.PointToFixed(float64(fontSize * 5 / 7))
	}*/

	// 添加干扰线
	AddInterferenceLines(4, 10, rand.Intn(2), img)
	// 添加噪点
	AddNoise(100, 200, img)

	/*file, err := os.Create("test" + code + ".png")
	if err != nil {
		log.Println("create png failed")
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Println("encode png file failed")
	}
	log.Println("Successfully!")*/
	return img
}

// distortPoint 对给定的点进行扭曲变形
func distortPoint(pt fixed.Point26_6, intensity, frequency float64) fixed.Point26_6 {
	// 计算扭曲的偏移量
	offsetX := int(intensity * float64(pt.Y) * math.Sin(2*math.Pi*float64(pt.X)/frequency))
	// 应用偏移量
	pt.X += fixed.Int26_6(offsetX)
	return pt
}

func distortText(code string, intensity, frequency float64) string {
	var distortedText string
	for _, ch := range code {
		// 将每个字符的坐标进行扭曲
		pt := freetype.Pt(0, 0)
		pt.X += fixed.I(rand.Intn(20))
		pt = distortPoint(pt, intensity, frequency)
		distortedText += string(ch)
	}
	return distortedText
}

func textRandomColorCount(dc *freetype.Context, count int) {
	// 随机颜色，确保颜色不太浅
	var fontColor color.RGBA
	for {
		r := rand.Intn(128) + 128
		g := rand.Intn(128) + 128
		b := rand.Intn(128) + 128
		sum := r + g + b
		if sum < count {
			fontColor = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
			break
		}
	}
	dc.SetSrc(&image.Uniform{C: fontColor})
}

func textRandomColor(dc *freetype.Context) {
	// 随机颜色，确保颜色不太浅
	var fontColor color.RGBA
	for {
		r := rand.Intn(128) + 128
		g := rand.Intn(128) + 128
		b := rand.Intn(128) + 128
		sum := r + g + b
		if sum < 600 {
			fontColor = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
			break
		}
	}
	dc.SetSrc(&image.Uniform{C: fontColor})
}

// AddInterferenceLines 添加干扰线
//
// 参数
//
//	min (int): 最小条数
//	max (int): 最大条数
//	lineWidth (int): 干扰线宽度
//	img (image.Image): 图片对象
func AddInterferenceLines(min, max, lineWidth int, img image.Image) {
	// 在图像上绘制随机直线来添加干扰线
	count := rand.Intn(max-min) + min
	for i := 0; i < count; i++ {
		x1 := rand.Intn(img.Bounds().Dx())
		y1 := rand.Intn(img.Bounds().Dy())
		x2 := rand.Intn(img.Bounds().Dx())
		y2 := rand.Intn(img.Bounds().Dy())
		lineColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}
		//lineWidth := rand.Intn(0)
		drawLine(img, x1, y1, x2, y2, lineColor, lineWidth)
	}
}

// drawLine 在图像上绘制直线
func drawLine(img image.Image, x1, y1, x2, y2 int, clr color.Color, width int) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		// 如果类型断言失败，说明 img 不是 *image.RGBA 类型
		log.Println("img is not *image.RGBA")
		return
	}

	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))
	if steps == 0 {
		return
	}
	dX := dx / float64(steps)
	dY := dy / float64(steps)
	x, y := float64(x1), float64(y1)
	if width == 0 {
		for i := 0; i <= steps; i++ {
			rgba.Set(int(x+0.5), int(y+0.5), clr)
			x += dX
			y += dY
		}
	} else {
		for i := 0; i <= steps; i++ {
			for w := -width / 2; w <= width/2; w++ {
				for h := -width / 2; h <= width/2; h++ {
					rgba.Set(int(x+float64(w)+rand.Float64()*0.5), int(y+float64(h)+rand.Float64()*0.8), clr)
				}
			}
			x += dX
			y += dY
		}
	}
}

// AddNoise 添加噪点
//
// 参数
//
//	min (int): 最小数量
//	max (int): 最大数量
//	img (image.Image): 图片对象
func AddNoise(min, max int, img image.Image) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		// 如果类型断言失败，说明 img 不是 *image.RGBA 类型
		log.Println("img is not *image.RGBA")
		return
	}
	count := rand.Intn(max-min) + min
	bounds := img.Bounds()
	for i := 0; i < count; i++ {
		x := rand.Intn(bounds.Max.X)
		y := rand.Intn(bounds.Max.Y)
		rgba.Set(x, y, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255})
	}
}

// getTextWidth 计算文本宽度
func getTextWidth(text string, font *truetype.Font, fontSize int) int {
	dc := freetype.NewContext()
	dc.SetFont(font)
	dc.SetFontSize(float64(fontSize))
	var width fixed.Int26_6
	for _, ch := range text {
		adv := font.HMetric(fixed.Int26_6(truetype.Index(ch)), truetype.Index(fontSize))

		width += adv.AdvanceWidth
	}
	return width.Ceil()
}

// GenerateMathCode 数学运算验证码
//
// 参数
//
//	width (int): 宽度
//	height (int): 高度
func GenerateMathCode(width, height int) (string, string, image.Image) {
	rand.NewSource(time.Now().UnixNano())

	/* fontPath, err := GetFontPath("/api/expressAPI", "/Arial Unicode.ttf")
	if err != nil {
		log.Println("获取字体文件路径异常")
		return "", "", nil
	} */
	//字体设置
	fontFile, err := os.ReadFile(fontAbsolutePath)
	if err != nil {
		log.Printf("open font file failed ===%+v", err)
		return "", "", nil
	}
	parse, err := truetype.Parse(fontFile)
	if err != nil {
		log.Printf("load font failed===%+v", err)
		return "", "", nil
	}

	// 选择运算符
	operators := []string{"+", "-", "*"}
	operator := operators[rand.Intn(len(operators))]

	num1 := rand.Intn(10) // 第一个随机数
	num2 := rand.Intn(10) // 第二个随机数

	// 根据运算符调整数值以保证结果的合理性
	switch operator {
	case "-":
		if num1 < num2 {
			num1, num2 = num2, num1
		}
	case "*":
		// 限制乘积的大小
		for num1*num2 > 100 {
			num2 = rand.Intn(10)
		}
	}

	captchaText := strconv.Itoa(num1) + operator + strconv.Itoa(num2) + "=?"
	var result string
	switch operator {
	case "+":
		result = strconv.Itoa(num1 + num2)
	case "-":
		result = strconv.Itoa(num1 - num2)
	case "*":
		result = strconv.Itoa(num1 * num2)
	}

	// 创建一个RGBA图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Pt(0, 0), draw.Src)

	dc := freetype.NewContext()
	dc.SetFont(parse)
	dc.SetFontSize(float64(27))
	dc.SetDPI(72)
	dc.SetDst(img)
	dc.SetClip(img.Bounds())
	dc.SetSrc(&image.Uniform{C: color.RGBA{A: 255}})

	textWidthOld := getTextWidth(captchaText, parse, 27)
	startX := (width-textWidthOld)/2 - height + 10 - textWidthOld*5
	pt := freetype.Pt(startX, 35)

	for _, ch := range captchaText {
		code := string(ch)
		if code == "+" || code == "-" || code == "*" || code == "=" || code == "?" {
			textRandomColorCount(dc, 400)
		} else {
			textRandomColor(dc)
		}
		_, err := dc.DrawString(code, pt)
		if err != nil {
			log.Printf("Draw string failed:%+v", err)
			return "", "", nil
		}
		// 更新下一个字符的位置
		pt.X += dc.PointToFixed(float64(32 * 5 / 6))
	}

	return captchaText, result, img
}

// GetFontPath 获取字体文件绝对路径
//
// 参数
//
//	currentRelativePath (string): 当前路径
//	fontPath (string): 字体路径
//
// 返回
//
//	string: 绝对路径
func GetFontPath(currentRelativePath, fontPath string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录错误===%+v", err)
		return "", err
	}
	log.Println("当前绝对路径", path)
	fontAbsolutePath := filepath.Join(strings.TrimSuffix(path, currentRelativePath), fontPath)
	log.Println("字体文件路径", fontAbsolutePath)
	return fontAbsolutePath, nil
}

// fontAbsolutePath 字体文件绝对路径
var fontAbsolutePath string

// init 初始化的时候获取字体文件绝对路径
// func init() {
// 	fontPath, err := GetFontPath("/api/expressAPI", "/Arial Unicode.ttf")
// 	if err != nil {
// 		// fmt.Printf("%+v", err)或者log.Printf("%+v", err)记录详细的错误信息，包括堆栈跟踪，非常适合用于生产环境中的错误记录和调试
// 		log.Printf("获取字体文件路径异常===%+v", err)
// 	}
// 	fontAbsolutePath = fontPath
// }
