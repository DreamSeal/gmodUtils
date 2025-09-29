package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

func CheckPic() {
	excelFileCheck := "C:/Users/73547/Desktop/py/rename/全部图片/全部/全部图片.xlsx"
	f, err := excelize.OpenFile(excelFileCheck)
	if err != nil {
		fmt.Printf("打开Excel失败: %v\n", err)
		return
	}
	defer f.Close()

	excelDir := filepath.Dir(excelFileCheck)
	if excelDir == "." {
		excelDir, _ = os.Getwd()
	}

	rows, _ := f.GetRows(f.GetSheetName(0))

	var existCount, notExistCount int

	fmt.Printf("在目录 %s 中检测文件...\n\n", excelDir)

	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 2 || rows[i][1] == "" {
			continue
		}

		fileName := rows[i][1]
		filePath := filepath.Join(excelDir, fileName)

		if fileExists(filePath) {
			fmt.Printf("✅ 第%d行: %s\n", i+1, fileName)
			existCount++
		} else {
			fmt.Printf("❌ 第%d行: %s\n", i+1, fileName)
			notExistCount++
		}
	}

	fmt.Printf("\n=== 检测结果 ===\n")
	fmt.Printf("存在的文件: %d 个\n", existCount)
	fmt.Printf("不存在的文件: %d 个\n", notExistCount)
	fmt.Printf("总计检测: %d 个文件\n", existCount+notExistCount)
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}
