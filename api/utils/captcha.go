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
	"time"
)

// 验证码字符集
var charset = "0123456789ABCDEFGHJKPLMNPQRSTUVWXYZ"

// GenerateCode 生成指定长度的随机验证码字符串
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
//		size 验证码个数
//	 typeStr 验证码类型(空字符串或者math默认为随机数字，其他为随机字符串)
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
func GenerateImage(code string) image.Image {
	// 图片尺寸
	width := 120
	height := 40

	// 创建空白图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景色设置为灰白色
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.ZP, draw.Src)

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
	draw.Draw(img, rect, &image.Uniform{color}, image.ZP, draw.Src)
}

//https://blog.csdn.net/m0_46198325/article/details/134913801

func CreateImage(code string) image.Image {
	rand.NewSource(time.Now().UnixNano())
	w := 140
	h := 50
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.ZP, draw.Src)

	//字体设置
	fontFile, err := os.ReadFile("Arial Unicode.ttf")
	if err != nil {
		log.Println("open file failed")
	}
	parse, err := truetype.Parse(fontFile)
	if err != nil {
		log.Println("load parse failed")
	}

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
			log.Println("Draw string failed:", err)
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
	addInterferenceLines(4, 10, rand.Intn(2), img)
	// 添加噪点
	addNoise(100, 200, img)

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

// addInterferenceLines 添加干扰线
func addInterferenceLines(min, max, lineWidth int, img *image.RGBA) {
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
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, clr color.Color, width int) {
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
			img.Set(int(x+0.5), int(y+0.5), clr)
			x += dX
			y += dY
		}
	} else {
		for i := 0; i <= steps; i++ {
			for w := -width / 2; w <= width/2; w++ {
				for h := -width / 2; h <= width/2; h++ {
					img.Set(int(x+float64(w)+rand.Float64()*0.5), int(y+float64(h)+rand.Float64()*0.8), clr)
				}
			}
			x += dX
			y += dY
		}
	}
}

// addNoise 添加噪点
func addNoise(min, max int, img *image.RGBA) {
	count := rand.Intn(max-min) + min
	bounds := img.Bounds()
	for i := 0; i < count; i++ {
		x := rand.Intn(bounds.Max.X)
		y := rand.Intn(bounds.Max.Y)
		img.Set(x, y, color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255})
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
