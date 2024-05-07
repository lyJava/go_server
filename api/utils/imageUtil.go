package utils

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"log"
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

func MergeImages2(imgPath1, imgPath2 string) {
	// 打开第一张图片
	img1, err := imaging.Open(imgPath1)
	if err != nil {
		log.Fatalf("failed to open image1.jpg: %v", err)
	}

	// 打开第二张图片
	img2, err := imaging.Open(imgPath2)
	if err != nil {
		log.Fatalf("failed to open image2.jpg: %v", err)
	}

	// 计算合并后图片的大小
	canvasWidth := img1.Bounds().Dx()
	canvasHeight := img1.Bounds().Dy() + img2.Bounds().Dy() + 1 // 加上间隔的高度
	dst := imaging.New(canvasWidth, canvasHeight, color.Transparent)

	// 将第一张图片绘制到新画布的顶部
	dst = imaging.Paste(dst, img1, image.Pt(0, 0))

	// 将第二张图片绘制到新画布的底部，并留下一段间隔
	dst = imaging.Paste(dst, img2, image.Pt(0, img1.Bounds().Dy()+1))

	// 将第一张图片绘制到新画布的顶部
	dst = imaging.Paste(dst, img1, image.Pt(0, 0))

	// 将第二张图片绘制到新画布的底部，并留下一段间隔
	dst = imaging.Paste(dst, img2, image.Pt(0, img1.Bounds().Dy()+1))

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
