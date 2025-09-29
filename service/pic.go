package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func Pic() {
	// 配置参数
	//imageFolder := "C:/Users/73547/Desktop/py/allPic" // 图片文件夹路径
	imageFolder := "C:/Users/73547/Desktop/py/rename/20250928/wuxing"
	excelFile := imageFolder + "/image_names.xlsx" // 输出的Excel文件名

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置表头
	headers := []string{"原始文件名", "新文件名"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	// 获取图片文件列表
	var imageFiles []string
	err := filepath.Walk(imageFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(path) {
			imageFiles = append(imageFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("遍历文件夹出错: %v\n", err)
		return
	}

	if len(imageFiles) == 0 {
		fmt.Printf("在文件夹 %s 中没有找到图片文件\n", imageFolder)
		return
	}

	// 处理每个图片文件
	for i, oldPath := range imageFiles {
		// 获取文件名（不含路径和后缀）
		oldNameWithExt := filepath.Base(oldPath)
		oldName := strings.TrimSuffix(oldNameWithExt, filepath.Ext(oldNameWithExt))
		ext := filepath.Ext(oldPath)

		// 生成随机文件名：纳秒级时间戳 + 随机数
		newName := generateNanoTimestampRandomName() + ext

		// 新文件路径
		newPath := filepath.Join(filepath.Dir(oldPath), newName)

		// 记录到Excel（第i+2行，因为第一行是表头）
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), oldName) // 只记录文件名，不含后缀
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), newName) // 新文件名包含后缀

		// 重命名文件
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("重命名文件 %s 出错: %v\n", oldNameWithExt, err)
			continue
		}

		fmt.Printf("已处理: %s -> %s\n", oldNameWithExt, newName)
	}

	// 保存Excel文件
	if err := f.SaveAs(excelFile); err != nil {
		fmt.Printf("保存Excel文件出错: %v\n", err)
		return
	}

	fmt.Printf("\n处理完成! 共处理 %d 个文件\n", len(imageFiles))
	fmt.Printf("Excel文件已保存: %s\n", excelFile)
	fmt.Printf("请妥善保管Excel文件，这是文件名对应的唯一记录！\n")
}

// 生成纳秒级时间戳 + 随机数的文件名
func generateNanoTimestampRandomName() string {
	// 获取当前时间的纳秒级时间戳
	timestamp := time.Now().UnixNano()

	// 生成6位随机字符串
	randomPart := generateRandomString(6)

	// 格式：纳秒时间戳_随机字符串
	return fmt.Sprintf("%d_%s", timestamp, randomPart)
}

// 生成指定长度的随机字符串（包含数字和字母）
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		// 使用加密安全的随机数生成器
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果加密随机失败，使用当前时间的纳秒作为备选
			result[i] = charset[time.Now().Nanosecond()%len(charset)]
		} else {
			result[i] = charset[num.Int64()]
		}
	}
	return string(result)
}

// 判断是否为图片文件
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".svg", ".heic", ".heif"}

	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}
