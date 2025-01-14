package utils

import (
	"bytes"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"mime/multipart"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/fogleman/gg"

	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

const (
	Rgba           = 0.6 //透明度
	Radians        = -30 //旋转度
	IntervalWidth  = 24  //水印间隔宽
	IntervalHeight = 100 //水印间隔高
	FontSize       = 24
)

// 定义一个可以实现 io.WriteSeeker 接口的类型
type pdfWriter struct {
	buf []byte
}

func (w *pdfWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *pdfWriter) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("Seek is not implemented ")
}

func AddTextWaterForPdfFile(file io.ReadSeeker, text string) ([]byte, error) {
	// 文字制作为pdf水印
	// 这里创建的参数是硬编码的，你可以根据需要动态调整
	watermarkText := text // 水印文字
	//pdfFontConf := "Helvetica" // 默认字体配置，可以根据需要调整
	//scale := 0.2 // 水印缩放比例，取值范围是0-1

	fontFile, err := os.Open("/Users/yangge/GolandProjects/apiProject/api/utils/simhei.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to open font file: %w", err)
	}
	defer fontFile.Close()

	fmt.Println("正常打开字体")

	api.EnsureDefaultConfigAt("/Users/yangge/GolandProjects/apiProject/upload")

	//err = api.InstallFonts([]string{"/Users/yangge/GolandProjects/apiProject/upload/pdfcpu/fonts/Simhei.gob"})
	//if err != nil {
	//	return nil, fmt.Errorf("failed to install font error: %w", err)
	//}

	//api.LoadConfiguration()
	//
	// install fonts from path
	//err = api.InstallFonts([]string{"./simhei.ttf"})
	//if err != nil {
	//	return nil, err
	//}

	fonts, err := api.ListFonts()
	if err != nil {
		return nil, fmt.Errorf("failed to list fonts: %w", err)
	}
	fmt.Println("all fonts:", fonts)
	api.LoadConfiguration()
	const PdfFontConfig = "points:42,opacity:0.3,rot:30,fontname:SimHei,scalefactor:1.5,aligntext:c"

	// 创建水印文本的配置
	textWm, err := pdfcpu.ParseTextWatermarkDetails(GetWaterMarkStr(watermarkText), PdfFontConfig, true, types.INCHES)
	//textWm, err := pdfcpu.ParseTextWatermarkDetails(GetWaterMarkStr2(watermarkText, 612, 792, 0.3), PdfFontConfig, true, types.INCHES)
	if err != nil {
		return nil, fmt.Errorf("failed to create text watermark: %w", err)
	}

	// 配置水印的缩放比例和透明度
	textWm.Scale = 0.4
	textWm.Opacity = 0.4 // 设置透明度
	textWm.Dx = 240
	textWm.Dy = 40
	//textWm.BgColor = &color.SimpleColor{R: 1.0, G: 0.6, B: 0.3}
	//textWm.Color = color.SimpleColor{R: 1.0, G: 0.6, B: 0.3}
	//textWm.Width = 120

	// 将水印添加到 PDF 文件
	pdfWriter := &pdfWriter{}
	var writer io.WriteSeeker = pdfWriter
	err = api.AddWatermarks(file, writer, nil, textWm, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to add watermark: %w", err)
	}

	return pdfWriter.buf, nil
}

// GetWaterMarkStr function for multiple lines text
func GetWaterMarkStr(waterMark string) string {
	textLen := len(waterMark)
	// automatically calculate the line count
	count := int(math.Round(425.0 / (float64(textLen) + 45)))
	var sb1 strings.Builder
	for i := 0; i < count; i++ {
		sb1.WriteString(waterMark)
		if i < count-1 {
			sb1.WriteString(strings.Repeat(" ", 5)) //文字间隔
		}
	}
	//single line
	sb1Str := sb1.String()
	//multiple lines
	var sb2 strings.Builder
	lineSpace := "\n \n \n \n \n \n"
	height := float64(len(lineSpace) * 2)
	rows := int(math.Round(475.0 / (height + 5)))
	for i := 0; i < rows; i++ {
		if i%2 == 0 {
			sb2.WriteString(sb1Str)
		} else {
			sb2.WriteString(strings.Repeat(" ", 2) + sb1Str[:len(sb1Str)-10])
		}
		if i < rows-1 {
			sb2.WriteString(lineSpace)
		}
	}
	return sb2.String()
}

// GetWaterMarkStr2 function for multiple lines text with dynamic distribution
func GetWaterMarkStr2(waterMark string, pageWidth, pageHeight float64, scale float64) string {
	textLen := len(waterMark)
	// Calculate the width of one watermarked text block
	textWidth := float64(textLen) * scale * 10.0 // assuming 10px width per character
	// Calculate the number of columns and rows based on page size
	cols := int(pageWidth / (textWidth))   // Increase factor to reduce columns
	rows := int(pageHeight / (scale * 25)) // Increase factor to reduce rows

	// Build the watermark string with larger spacing
	var sb strings.Builder
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			sb.WriteString(waterMark)
			if j < cols-1 {
				sb.WriteString(strings.Repeat(" ", 10)) // Increased space between words
			}
		}
		if i < rows-1 {
			sb.WriteString("\n") // new line after each row
		}
	}
	return sb.String()
}

// AddMultipleWatermarksToPdf 向PDF文件添加多个水印
func AddMultipleWatermarksToPdf(filePath string, texts []string) ([]byte, error) {
	// 打开PDF文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	//defer file.Close()

	// 读取PDF文件并解析
	ctx, err := api.ReadContext(file, nil)
	if err != nil {
		fmt.Printf("failed to read PDF context: %v\n", err)
	}

	// 输出文件的基本信息
	fmt.Printf("PDF 文件：%s\n", filePath)
	fmt.Printf("页数：%d\n", ctx.PageCount)

	// 输出PDF的元数据（如有）
	fmt.Printf("标题：%s\n", ctx.Title)
	fmt.Printf("作者：%s\n", ctx.Author)
	fmt.Printf("主题：%s\n", ctx.Subject)
	fmt.Printf("关键词：%s\n", ctx.Keywords)

	// 确认PDF的页数
	//if ctx.PageCount == 0 {
	//	return nil, fmt.Errorf("no pages found in the PDF file")
	//}

	fmt.Printf("pdf的页数：%d", ctx.PageCount)

	// 将水印添加到每一页
	pdfWriter := &pdfWriter{}
	var writer io.WriteSeeker = pdfWriter
	// 对每一页添加多个水印
	for page := 1; page <= ctx.PageCount; page++ {
		for _, text := range texts {
			// 创建水印
			watermark := &model.Watermark{
				OnTop:      true,
				TextString: text,
				InpUnit:    types.POINTS,
				Scale:      0.5,                                 // 缩放比例
				Opacity:    0.5,                                 // 透明度
				Color:      color.SimpleColor{R: 0, G: 0, B: 1}, // 设置水印颜色为蓝色
				FontName:   "Helvetica",                         // 设置字体
				FontSize:   24,                                  // 设置字体大小
			}

			// 对每一页添加一个水印
			err := api.AddWatermarks(file, writer, nil, watermark, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to add watermark to page: %w", err)
			}
		}
	}
	return pdfWriter.buf, nil
}

// AddTextWaterForPdf1 添加多个水印到PDF
func AddTextWaterForPdf1(filePath string, text string) ([]byte, error) {
	// 打开PDF文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 创建水印文本配置
	textWm, err := pdfcpu.ParseTextWatermarkDetails(text, "", true, types.POINTS)
	if err != nil {
		return nil, fmt.Errorf("failed to create text watermark: %w", err)
	}

	// 配置水印的缩放比例和透明度
	textWm.Scale = 0.5
	textWm.Opacity = 0.5 // 设置透明度

	// 获取PDF的页面数量
	ctx, err := api.ReadContext(file, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF context: %w", err)
	}

	pageCount := ctx.PageCount

	// 创建一个新的buffer来保存输出的PDF
	var buf []byte

	for page := 1; page <= pageCount; page++ {
		pdfWriter := &pdfWriter{}
		var writer io.WriteSeeker = pdfWriter
		// 每页添加水印
		err = api.AddWatermarks(file, writer, nil, textWm, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to add watermark to page %d: %w", page, err)
		}
	}

	return buf, nil
}

// PdfFontConfDesc PdfFontConfDesc比例因子:1, 字体大小points:12, 透明度opacity:0.6, 旋转角度rotation:30
const PdfFontConfDesc = "sc:1,points:12,opacity:0.6,rot:30"

func AddTextWaterForImage(file multipart.File, fileName string, text string) ([]byte, error) {

	defer file.Close()

	//文件转换为image对象
	var imageFile image.Image
	if strings.HasSuffix(fileName, "png") {
		i, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		imageFile = i
	}
	if strings.HasSuffix(fileName, "jpeg") || strings.HasSuffix(fileName, "jpg") {
		i, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		imageFile = i
	}
	if strings.HasSuffix(fileName, "bmp") {
		i, err := bmp.Decode(file)
		if err != nil {
			return nil, err
		}
		imageFile = i
	}
	if strings.HasSuffix(fileName, "webp") {
		i, err := webp.Decode(file)
		if err != nil {
			return nil, err
		}
		imageFile = i
	}
	if imageFile == nil {
		return nil, errors.New(fmt.Sprintf("invalid file suffix:%s", fileName))
	}

	//将源文件作为底层画布
	dc := gg.NewContextForImage(imageFile)

	//加载字体对象，本身加载16进制字体，无需预先加载
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	//默认字体大小格式等，不设置会有默认值
	face := truetype.NewFace(font, &truetype.Options{Size: FontSize})

	//设置字体
	dc.SetFontFace(face)
	//设置颜色，透明度
	dc.SetRGBA(Rgba, Rgba, Rgba, Rgba)

	//画布的 x/y轴大小
	maxX := dc.Width()
	maxY := dc.Height()

	//设置旋转度，后两个字段代表以 画布中心为 旋转点
	dc.RotateAbout(gg.Radians(Radians), float64(maxX/2), float64(maxY/2))

	//连续水印
	textW, _ := dc.MeasureString(text) //文字宽度
	width := int(math.Ceil(textW))
	for i := -maxX / 2; i <= maxX+maxX/2; i += width + IntervalWidth {
		for j := -maxY / 2; j <= maxY+maxY/2; j += IntervalHeight {
			dc.DrawString(text, float64(i), float64(j))
		}
	}

	// 将生成的图片转换成buffer
	buffer := bytes.NewBuffer(make([]byte, 0, 512))
	err = jpeg.Encode(buffer, dc.Image(), &jpeg.Options{
		Quality: jpeg.DefaultQuality,
	})
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil

}

func AddTextWaterForPdf(file io.Reader, text string) ([]byte, error) {
	// 将 file 转换为 *bytes.Buffer
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// 将 *bytes.Buffer 转换为 *bytes.Reader，后者实现了 io.ReadSeeker
	reader := bytes.NewReader(buf.Bytes())
	//文字制作为pdf水印
	//该工具会保证单行文字全部显示，所以会根据文字长度等比例强行缩放文字大小，所以需要根据传入的文字长度来设置每行的文字次数；
	count := 4
	textLen := len(text)
	switch {
	case textLen > 25 && textLen <= 40:
		count = 3
	case textLen > 40 && textLen <= 60:
		count = 2
	case textLen > 60:
		count = 1
	}
	var sb1 strings.Builder
	for i := 0; i < count; i++ {
		sb1.WriteString(text)
		if i < count-1 {
			sb1.WriteString("     ") //文字间隔
		}
	}
	//单行文字制作完成
	sb1Str := sb1.String()
	//拼接成多行文字
	var sb2 strings.Builder
	for i := 0; i < 10; i++ { //最多10行文字
		sb2.WriteString(sb1Str)
		if i < 9 {
			sb2.WriteString("\n \n \n \n \n \n \n \n \n \n")
		}
	}
	//制作成文字水印
	textWm, err := pdfcpu.ParseTextWatermarkDetails(sb2.String(), PdfFontConfDesc, true, types.POINTS)
	if err != nil {
		return nil, err
	}
	//将水印添加到pdf文件并生成新文件
	var resultBuf bytes.Buffer

	// 创建一个内存缓冲区来保存生成的 PDF
	err = api.AddWatermarks(reader, &resultBuf, nil, textWm, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
