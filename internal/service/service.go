package service

import (
	"bufio"
	"bytes"
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"
)

type XlsxService interface {
	Generate(request model.XlsxRequest) (file []byte, err error)
}

type XlsxServiceImpl struct {
	log *zap.Logger
}

func New(log *zap.Logger) XlsxService {
	return XlsxServiceImpl{
		log: log,
	}
}

func (s XlsxServiceImpl) Generate(request model.XlsxRequest) (file []byte, err error) {
	wb := xlsx.NewFile()

	for _, s := range request.Sheets {
		createSheet(wb, &s)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err = wb.Write(writer)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func createSheet(wb *xlsx.File, sheet *model.Sheet) {
	s, err := wb.AddSheet(sheet.Name)
	if err != nil {
		return
	}
	currentRowIndex := 0
	currentRowIndex = createData(s, currentRowIndex, "Test_value")
}

func createData(sheet *xlsx.Sheet, currentRowIndex int, val string) (cIndex int) {
	_, err := sheet.AddRowAtIndex(currentRowIndex)
	if err != nil {
		return
	}

	createCell(sheet, currentRowIndex, 0, val)
	return currentRowIndex
}

func createCell(sheet *xlsx.Sheet, row int, col int, val string) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return
	}
	cell.Value = val
}
