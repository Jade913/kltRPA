package utils

import (
	"os"

	"github.com/xuri/excelize/v2"
)

func ImportTableFromExcel(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	// 获取第一个工作表的名称
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func SaveFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, os.ModePerm)
}
