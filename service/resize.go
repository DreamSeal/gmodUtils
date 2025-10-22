package service

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func Resize() {
	// 获取用户输入的目录
	//sourceDir := "C:/Users/73547/Desktop/py/rename/全部图片/全部头像/全部" // 当前目录
	sourceDir := "C:/Users/73547/Desktop/py/rename/test"

	// 检查目录是否存在
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		fmt.Printf("目录不存在: %s\n", sourceDir)
		return
	}

	// 输出目录为源目录下的 yasuo 文件夹
	outputDir := filepath.Join(sourceDir, "yasuo")

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 支持的图片格式
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}

	fmt.Println("\n开始压缩图片...")
	fmt.Printf("源目录: %s\n", sourceDir)
	fmt.Printf("目标尺寸: 254x254\n")
	fmt.Printf("目标大小: < 30KB\n")
	fmt.Printf("输出目录: %s\n\n", outputDir)

	// 遍历指定目录
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		return
	}

	// 先列出所有图片文件
	var imageFiles []string
	for _, file := range files {
		if file.IsDir() || file.Name() == "yasuo" {
			continue
		}

		if resizeIsImageFile(file.Name(), imageExtensions) {
			fullPath := filepath.Join(sourceDir, file.Name())
			imageFiles = append(imageFiles, fullPath)
		}
	}

	fmt.Printf("找到 %d 个图片文件:\n", len(imageFiles))
	for _, filepath := range imageFiles {
		fmt.Printf("  - %s\n", filepath)
	}
	fmt.Println()

	processedCount := 0
	for _, filepath := range imageFiles {
		// 检查文件是否存在
		fileInfo, err := os.Stat(filepath)
		if err != nil {
			fmt.Printf("获取文件信息失败 %s: %v\n", filepath, err)
			continue
		}

		// 检查文件大小
		fileSizeKB := float64(fileInfo.Size()) / 1024
		if fileSizeKB <= 30 {
			fmt.Printf("跳过小文件: %s (%.2f KB)\n", filepath, fileSizeKB)
			continue
		}

		fmt.Printf("处理文件: %s (%.2f KB)\n", filepath, fileSizeKB)

		// 压缩图片
		if err := compressImage(filepath, outputDir); err != nil {
			fmt.Printf("压缩失败 %s: %v\n", filepath, err)
		} else {
			processedCount++
		}
	}

	fmt.Printf("\n处理完成！共压缩 %d 个文件到 %s 目录\n", processedCount, outputDir)
}

// 检查是否为图片文件
func resizeIsImageFile(filename string, extensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// 简单的图片缩放函数
func resizeImage(img image.Image, width, height int) image.Image {
	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// 创建目标图片
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// 计算缩放比例
	xRatio := float64(srcW) / float64(width)
	yRatio := float64(srcH) / float64(height)

	// 简单的最近邻缩放算法
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)

			if srcX >= srcW {
				srcX = srcW - 1
			}
			if srcY >= srcH {
				srcY = srcH - 1
			}

			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

// 压缩图片
func compressImage(inputPath, outputDir string) error {
	// 再次检查文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", inputPath)
	}

	// 打开图片文件
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 解码图片
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("解码图片失败: %v", err)
	}

	// 获取原始尺寸
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	fmt.Printf("  原始尺寸: %dx%d\n", originalWidth, originalHeight)

	// 调整尺寸为 254x254
	var resizedImg image.Image
	if originalWidth != 254 || originalHeight != 254 {
		resizedImg = resizeImage(img, 254, 254)
		fmt.Printf("  调整尺寸: 254x254\n")
	} else {
		resizedImg = img
		fmt.Printf("  尺寸正确，无需调整\n")
	}

	// 构建输出路径（保持原文件名）
	outputPath := filepath.Join(outputDir, filepath.Base(inputPath))

	// 根据格式保存图片
	if strings.ToLower(format) == "png" {
		// PNG格式
		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("创建输出文件失败: %v", err)
		}
		defer outFile.Close()

		if err := png.Encode(outFile, resizedImg); err != nil {
			return fmt.Errorf("PNG编码失败: %v", err)
		}
	} else {
		// JPEG格式 - 使用质量调整来压缩
		quality := 85 // 初始质量

		for quality >= 20 {
			outFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("创建输出文件失败: %v", err)
			}

			// 使用当前质量保存JPEG
			if err := jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: quality}); err != nil {
				outFile.Close()
				return fmt.Errorf("JPEG编码失败: %v", err)
			}
			outFile.Close()

			// 获取文件大小
			fileInfo, _ := os.Stat(outputPath)
			finalSize := fileInfo.Size()
			finalSizeKB := float64(finalSize) / 1024

			fmt.Printf("  质量 %d%%: %.2f KB\n", quality, finalSizeKB)

			if finalSize <= 30*1024 {
				break // 达到目标大小
			}

			quality -= 15         // 降低质量继续尝试
			os.Remove(outputPath) // 删除过大的文件
		}

		// 如果仍然太大，使用最低质量
		if quality < 20 {
			outFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("创建输出文件失败: %v", err)
			}
			defer outFile.Close()

			if err := jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 20}); err != nil {
				return fmt.Errorf("最终压缩失败: %v", err)
			}
		}
	}

	// 获取最终文件大小
	outputInfo, _ := os.Stat(outputPath)
	finalSizeKB := float64(outputInfo.Size()) / 1024
	fmt.Printf("  压缩完成: %.2f KB\n\n", finalSizeKB)

	return nil
}
