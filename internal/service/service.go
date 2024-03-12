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

	for _, sheet := range request.Sheets {
		s.createSheet(wb, &sheet)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err = wb.Write(writer)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (s XlsxServiceImpl) createSheet(wb *xlsx.File, sheet *model.Sheet) {
	wbSheet, err := wb.AddSheet(sheet.Name)
	if err != nil {
		return
	}
	currentRowIndex := 0
	currentRowIndex = s.createAdditionalData(wbSheet, currentRowIndex, sheet.AdditionalInfo.Top)
	currentRowIndex = s.createColumns(wbSheet, currentRowIndex, sheet.Columns)
	currentRowIndex = s.createTable(wbSheet, currentRowIndex, sheet.Columns, sheet.Data)
	currentRowIndex = s.createAdditionalData(wbSheet, currentRowIndex, sheet.AdditionalInfo.Bottom)
}

func (s XlsxServiceImpl) createAdditionalData(sheet *xlsx.Sheet, currentRowIndex int, data *[]model.AdditionalData) (cIndex int) {
	if data != nil {
		for _, a := range *data {
			_, err := sheet.AddRowAtIndex(currentRowIndex)
			if err != nil {
				return
			}
			s.createCell(sheet, currentRowIndex, 0, a.Title)
			s.createCell(sheet, currentRowIndex, 1, a.Value)
			currentRowIndex++
		}
	}
	return currentRowIndex
}

func (s XlsxServiceImpl) createColumns(sheet *xlsx.Sheet, currentRowIndex int, columns *[]model.Column) (cIndex int) {
	if columns != nil {
		_, err := sheet.AddRowAtIndex(currentRowIndex)
		if err != nil {
			return
		}
		for i, c := range *columns {
			s.createCell(sheet, currentRowIndex, i, c.Title)
		}
		currentRowIndex++
	}
	return currentRowIndex
}

func (s XlsxServiceImpl) createTable(sheet *xlsx.Sheet, currentRowIndex int, columns *[]model.Column, data *[]model.Data) (cIndex int) {
	if columns != nil && data != nil {
		for _, d := range *data {
			_, err := sheet.AddRowAtIndex(currentRowIndex)
			if err != nil {
				return
			}
			for j, c := range *columns {
				s.createCell(sheet, currentRowIndex, j, d[c.ID])

			}
			currentRowIndex++
		}
	}
	return currentRowIndex
}

func (s XlsxServiceImpl) createCell(sheet *xlsx.Sheet, row int, col int, val string) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return
	}
	cell.Value = val
}
