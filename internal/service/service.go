package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/illenko/xlsx-generator/internal/logger"
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/illenko/xlsx-generator/internal/style"
	"github.com/tealeg/xlsx/v3"
	"log/slog"
)

type XlsxService interface {
	Generate(ctx context.Context, request model.XlsxRequest) (file []byte, err error)
}

type XlsxServiceImpl struct {
	log *slog.Logger
}

func New(log *slog.Logger) XlsxService {
	return XlsxServiceImpl{log: log}
}

func (s XlsxServiceImpl) Generate(ctx context.Context, request model.XlsxRequest) (file []byte, err error) {
	wb := xlsx.NewFile()

	for _, sheet := range request.Sheets {
		s.createSheet(ctx, wb, &sheet)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err = wb.Write(writer)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (s XlsxServiceImpl) createSheet(ctx context.Context, wb *xlsx.File, sheet *model.Sheet) {
	wbSheet, err := wb.AddSheet(sheet.Name)
	if err != nil {
		return
	}

	ctx = logger.AppendCtx(ctx, slog.String("sheet_name", sheet.Name))
	s.log.InfoContext(ctx, "Created sheet")
	currentRowIndex := 0
	currentRowIndex = s.createAdditionalData(ctx, wbSheet, currentRowIndex, sheet.AdditionalInfo.Top)
	currentRowIndex = s.createColumns(ctx, wbSheet, currentRowIndex, sheet.Columns)
	currentRowIndex = s.createTable(ctx, wbSheet, currentRowIndex, sheet.Columns, sheet.Data)
	currentRowIndex = s.createAdditionalData(ctx, wbSheet, currentRowIndex, sheet.AdditionalInfo.Bottom)
	s.setAutoWidth(wbSheet, sheet)
}

func (s XlsxServiceImpl) createAdditionalData(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int, data *[]model.AdditionalData) (cIndex int) {
	if data != nil {
		for _, a := range *data {
			_, err := sheet.AddRowAtIndex(currentRowIndex)
			if err != nil {
				return
			}
			s.createCell(sheet, currentRowIndex, 0, a.Title, style.DefaultAdditionalInfoTitle)
			s.createCell(sheet, currentRowIndex, 1, a.Value, style.Default)
			s.log.InfoContext(ctx, fmt.Sprintf("Created additional data row, title: %v, value: %v, index: %v", a.Title, a.Value, currentRowIndex))
			currentRowIndex++
		}
		currentRowIndex = s.emptyRow(sheet, currentRowIndex)
	}
	return currentRowIndex
}

func (s XlsxServiceImpl) createColumns(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int, columns *[]model.Column) (cIndex int) {
	if columns != nil {
		_, err := sheet.AddRowAtIndex(currentRowIndex)
		if err != nil {
			return
		}
		for i, c := range *columns {
			s.createCell(sheet, currentRowIndex, i, c.Title, style.Resolve(c.Color, style.DefaultColumn))
			s.log.InfoContext(ctx, fmt.Sprintf("Created column, title: %v, index: %v", c.Title, i))
		}
		currentRowIndex++
	}
	return currentRowIndex
}

func (s XlsxServiceImpl) createTable(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int, columns *[]model.Column, data *[]model.Data) (cIndex int) {
	if columns != nil && data != nil {
		for _, d := range *data {
			_, err := sheet.AddRowAtIndex(currentRowIndex)
			if err != nil {
				return
			}
			for j, c := range *columns {
				s.createCellWithType(sheet, currentRowIndex, j, d[c.ID], style.Default, c.Type)
				s.log.InfoContext(ctx, fmt.Sprintf("Created table cell, value: '%v', row index: %v, col index: %v", d[c.ID], currentRowIndex, j))
			}
			currentRowIndex++
		}
		currentRowIndex = s.emptyRow(sheet, currentRowIndex)
	}
	return currentRowIndex
}

func (s XlsxServiceImpl) emptyRow(sheet *xlsx.Sheet, currentRowIndex int) (cIndex int) {
	_, err := sheet.AddRowAtIndex(currentRowIndex)
	if err != nil {
		return
	}
	currentRowIndex++
	return currentRowIndex
}

func (s XlsxServiceImpl) createCell(sheet *xlsx.Sheet, row int, col int, val string, style *xlsx.Style) (cell *xlsx.Cell) {
	return s.createCellWithType(sheet, row, col, val, style, nil)
}

func (s XlsxServiceImpl) createCellWithType(sheet *xlsx.Sheet, row int, col int, val string, style *xlsx.Style, cellType *model.CellType) (cell *xlsx.Cell) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return
	}
	s.setCellValue(cellType, cell, val)
	cell.SetStyle(style)
	return
}

func (s XlsxServiceImpl) setCellValue(columnType *model.CellType, cell *xlsx.Cell, val string) {
	if columnType == nil || *columnType == model.StringCell {
		cell.Value = val
	} else if *columnType == model.NumberCell {
		cell.SetNumeric(val)
	}
}

func (s XlsxServiceImpl) setAutoWidth(wbSheet *xlsx.Sheet, sheet *model.Sheet) {
	length := 2
	if sheet.Columns != nil && len(*sheet.Columns) > length {
		length = len(*sheet.Columns)
	}
	setWidth(wbSheet, length)
}

func setWidth(wbSheet *xlsx.Sheet, length int) {
	for i := 1; i <= length; i++ {
		err := wbSheet.SetColAutoWidth(i, xlsx.DefaultAutoWidth)
		if err != nil {
			return
		}
	}
}
