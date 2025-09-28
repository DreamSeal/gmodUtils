package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type RowData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func readExcelFile(filename string) ([]RowData, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 获取所有工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel文件中没有工作表")
	}

	// 使用第一个工作表
	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	var result []RowData

	for i, row := range rows {
		// 跳过空行
		if len(row) == 0 {
			continue
		}

		// 处理列数不足的情况
		firstCol := ""
		secondCol := ""

		if len(row) > 0 {
			firstCol = row[0]
		}
		if len(row) > 1 {
			secondCol = row[1]
		}

		record := RowData{
			Key:   firstCol,
			Value: secondCol,
		}
		result = append(result, record)

		fmt.Printf("行%d: Key='%s', Value='%s'\n", i+1, record.Key, record.Value)
	}

	return result, nil
}

func ReadTest() {
	data, err := readExcelFile("image_names.xlsx")
	if err != nil {
		return
	}

	fmt.Printf("\n成功读取 %d 行数据:\n", len(data))
	for i, item := range data {
		fmt.Printf("%d. %s -> %s\n", i+1, item.Key, item.Value)
	}
}
