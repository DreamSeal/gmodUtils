package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
)

// 配置参数
var (
	renameImageFolder = "C:/Users/73547/Desktop/py/rename/vc" // 图片文件夹路径
	// Excel文件路径
	excelFile = renameImageFolder + "/名字图片对照表.xlsx" // 输出的Excel文件名
	newPath   = "newNames"
)

//重命名
func Rename() {
	// 打开Excel文件
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Printf("打开Excel文件失败: %v\n", err)
		return
	}
	defer f.Close()

	// 获取第一个工作表名
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		fmt.Println("未找到工作表")
		return
	}

	// 获取所有行数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Printf("读取Excel数据失败: %v\n", err)
		return
	}

	if len(rows) < 2 {
		fmt.Println("Excel文件至少需要2行数据（表头+数据）")
		return
	}

	// 获取当前目录
	currentDir := renameImageFolder

	fmt.Printf("检测目录: %s\n", currentDir)
	fmt.Printf("开始处理文件...\n\n")

	if err := createFolderIfNotExists(currentDir + "/" + newPath); err != nil {
		fmt.Println(err)
	}
	// 从第二行开始处理（索引1）
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		if len(rows[rowIndex]) < 2 {
			continue // 跳过没有第二列数据的行
		}

		oldFileName := rows[rowIndex][1] // 第二列数据（索引1）
		if oldFileName == "" {
			continue // 跳过空文件名
		}

		// 构建完整的文件路径
		oldFilePath := filepath.Join(currentDir, oldFileName)

		// 检查文件是否存在
		if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
			fmt.Printf("❌ 文件不存在: %s\n", oldFileName)
			continue
		}

		// 生成新文件名（这里可以根据需要自定义命名规则）
		newFileName := generateNewFileName(oldFileName)
		newFilePath := filepath.Join(currentDir, newPath+"/"+newFileName)

		// 重命名文件
		err = os.Rename(oldFilePath, newFilePath)
		if err != nil {
			fmt.Printf("❌ 重命名失败 %s -> %s: %v\n", oldFileName, newFileName, err)
			continue
		}

		// 更新Excel中的文件名（第二列，B列）
		cellName, _ := excelize.CoordinatesToCellName(2, rowIndex+1) // 第二列，行号从1开始
		f.SetCellValue(sheetName, cellName, newFileName)

		fmt.Printf("✅ 成功重命名: %s -> %s\n", oldFileName, newFileName)
	}

	// 保存修改后的Excel文件
	if err := f.Save(); err != nil {
		fmt.Printf("保存Excel文件失败: %v\n", err)
		return
	}

	fmt.Printf("\n✅ 所有文件处理完成！Excel文件已更新。\n")
}

// 生成新文件名的函数（根据你的需求自定义）
func generateNewFileName(oldFileName string) string {
	//原始文件
	//baseName := strings.TrimSuffix(oldFileName, filepath.Ext(oldFileName))

	// 获取后缀
	ext := filepath.Ext(oldFileName)

	// 生成随机文件名：纳秒级时间戳 + 随机数
	newName := generateNanoTimestampRandomName() + ext
	return newName
}

func createFolderIfNotExists(folderPath string) error {
	// 检查文件夹是否存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 文件夹不存在，创建它
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			return fmt.Errorf("创建文件夹失败: %v", err)
		}
		fmt.Printf("文件夹创建成功: %s\n", folderPath)
	} else {
		fmt.Printf("文件夹已存在: %s\n", folderPath)
	}
	return nil
}
