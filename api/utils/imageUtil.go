package utils

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func MergeImage() {
	// 打开第一张图片
	img1, err := imaging.Open("./image/test15.png")
	if err != nil {
		log.Fatalf("failed to open image1.jpg: %v", err)
	}

	// 打开第二张图片
	img2, err := imaging.Open("./image/test16.png")
	if err != nil {
		log.Fatalf("failed to open image2.jpg: %v", err)
	}

	// 将第二张图片调整大小以适应第一张图片的大小
	img2 = imaging.Resize(img2, img1.Bounds().Dx(), img1.Bounds().Dy(), imaging.Lanczos)

	// 创建一个新的画布，大小为两张图片的宽度之和，高度取两张图片中最大的高度
	canvasWidth := img1.Bounds().Dx() + img2.Bounds().Dx()
	canvasHeight := img1.Bounds().Dy()
	if img2.Bounds().Dy() > canvasHeight {
		canvasHeight = img2.Bounds().Dy()
	}
	dst := imaging.New(canvasWidth, canvasHeight, color.Transparent)

	// 将第一张图片绘制到新画布的左侧
	dst = imaging.Paste(dst, img1, image.Pt(0, 0))

	// 将第二张图片绘制到新画布的右侧
	dst = imaging.Paste(dst, img2, image.Pt(img1.Bounds().Dx(), 0))

	rand.NewSource(time.Now().UnixNano())
	round := rand.Intn(100)
	idStr := strconv.Itoa(round)

	// 保存合并后的图片
	err = imaging.Save(dst, "./image/merge/combined_image"+idStr+".jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	log.Println("Images combined successfully")
}

func MergeImages2(img1Path, img2Path string, padding int) {
	// 打开第一张图片
	img1, err := imaging.Open(img1Path)
	if err != nil {
		log.Fatalf("failed to open image1.jpg: %v", err)
	}

	// 打开第二张图片
	img2, err := imaging.Open(img2Path)
	if err != nil {
		log.Fatalf("failed to open image2.jpg: %v", err)
	}

	// 计算合成后图片的高度
	canvasHeight := img1.Bounds().Dy() + img2.Bounds().Dy() + padding

	// 画布的宽度
	canvasWidth := 0

	img1Width := img1.Bounds().Dx()
	img2Width := img2.Bounds().Dx()
	/*if img1Width >= img2Width {
		canvasWidth = img1Width
	} else {
		canvasWidth = img2Width
	}*/
	// 以最宽的图片宽度为准
	canvasWidth = int(math.Max(float64(img1Width), float64(img2Width)))

	//width, height := checkWidthAndHeight(canvasWidth, canvasHeight)

	dst := imaging.New(canvasWidth, canvasHeight, color.Transparent)

	// 将第一张图片绘制到新画布的顶部
	dst = imaging.Paste(dst, img1, image.Pt(0, 0))

	// 将第二张图片绘制到新画布的底部，并留下一段间隔
	dst = imaging.Paste(dst, img2, image.Pt(0, img1.Bounds().Dy()+padding))

	rand.NewSource(time.Now().UnixNano())
	round := rand.Intn(100)
	idStr := strconv.Itoa(round)

	// 保存合并后的图片
	err = imaging.Save(dst, "./image/merge/image_"+idStr+".jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	log.Println("Images combined successfully")

}

// checkWidthAndHeight 验证图片宽度与高度
func checkWidthAndHeight(width, height int) (newWidth, newHeight int) {
	// 如果宽度超过A4纸张的宽度（2480像素），则进行裁剪
	a4Width := 2480
	if width > a4Width {
		width = a4Width
	}
	// 如果高度超过A4纸张的高度（3508像素），则进行裁剪
	a4Height := 3508
	if height > a4Height {
		height = a4Height
	}
	return width, height
}
