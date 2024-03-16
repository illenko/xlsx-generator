package service

import (
	"context"
	"fmt"
	"github.com/illenko/xlsx-generator/internal/logger"
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/illenko/xlsx-generator/internal/style"
	"github.com/tealeg/xlsx/v3"
	"log/slog"
)

type XlsxService interface {
	Generate(ctx context.Context, request model.XlsxRequest) *xlsx.File
}

type xlsxService struct {
	log *slog.Logger
}

func New(log *slog.Logger) XlsxService {
	return xlsxService{log: log}
}

func (s xlsxService) Generate(ctx context.Context, request model.XlsxRequest) *xlsx.File {
	wb := xlsx.NewFile()

	s.createSheets(ctx, request, wb)

	return wb
}

func (s xlsxService) createSheets(ctx context.Context, request model.XlsxRequest, wb *xlsx.File) {
	for _, sheet := range request.Sheets {
		s.createSheet(ctx, wb, &sheet)
	}
}

func (s xlsxService) createSheet(ctx context.Context, wb *xlsx.File, sheet *model.Sheet) {
	wbSheet, err := wb.AddSheet(sheet.Name)
	ctx = logger.AppendCtx(ctx, slog.String("sheet_name", sheet.Name))

	if err != nil {
		s.log.ErrorContext(ctx, "Error while creating sheet")
		return
	}

	s.log.InfoContext(ctx, "Created sheet")
	currentRowIndex := 0
	if sheet.AdditionalInfo != nil {
		currentRowIndex = s.setAdditionalData(ctx, wbSheet, currentRowIndex, sheet.AdditionalInfo.Top)
	}
	currentRowIndex = s.setColumns(ctx, wbSheet, currentRowIndex, sheet.Columns)
	currentRowIndex = s.setTable(ctx, wbSheet, currentRowIndex, sheet.Columns, sheet.Data)
	if sheet.AdditionalInfo != nil {
		s.setAdditionalData(ctx, wbSheet, currentRowIndex, sheet.AdditionalInfo.Bottom)
	}
	s.adjustColWidth(ctx, wbSheet, sheet)
}

func (s xlsxService) setAdditionalData(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int, data *[]model.AdditionalData) (cIndex int) {
	if data != nil {
		for _, a := range *data {
			currentRowIndex = s.createRow(ctx, sheet, currentRowIndex)
			s.createCell(ctx, sheet, currentRowIndex, 0, a.Title, style.DefaultAdditionalInfoTitle)
			s.createCell(ctx, sheet, currentRowIndex, 1, a.Value, style.Default)
			s.log.InfoContext(ctx, fmt.Sprintf("Created additional data row, title: %v, value: %v, index: %v", a.Title, a.Value, currentRowIndex))
			currentRowIndex++
		}
		currentRowIndex = s.createRow(ctx, sheet, currentRowIndex)
	}
	return currentRowIndex
}

func (s xlsxService) setColumns(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int, columns *[]model.Column) (cIndex int) {
	if columns != nil {
		currentRowIndex = s.createRow(ctx, sheet, currentRowIndex)
		for i, c := range *columns {
			s.createCell(ctx, sheet, currentRowIndex, i, c.Title, style.Resolve(c.Color, style.DefaultColumn))
			s.log.InfoContext(ctx, fmt.Sprintf("Created column, title: %v, index: %v", c.Title, i))
		}
		currentRowIndex++
	}
	return currentRowIndex
}

func (s xlsxService) setTable(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int, columns *[]model.Column, data *[]model.Data) (cIndex int) {
	if columns != nil && data != nil {
		for _, d := range *data {
			currentRowIndex = s.createRow(ctx, sheet, currentRowIndex)
			for j, c := range *columns {
				s.createCellWithType(ctx, sheet, currentRowIndex, j, d[c.ID], style.Default, c.Type)
				s.log.InfoContext(ctx, fmt.Sprintf("Created table cell, value: '%v', row index: %v, col index: %v", d[c.ID], currentRowIndex, j))
			}
			currentRowIndex++
		}
		currentRowIndex = s.createEmptyRow(ctx, sheet, currentRowIndex)
	}
	return currentRowIndex
}

func (s xlsxService) createEmptyRow(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int) (cIndex int) {
	currentRowIndex = s.createRow(ctx, sheet, currentRowIndex)
	currentRowIndex++
	return currentRowIndex
}

func (s xlsxService) createRow(ctx context.Context, sheet *xlsx.Sheet, currentRowIndex int) (cIndex int) {
	_, err := sheet.AddRowAtIndex(currentRowIndex)
	if err != nil {
		s.log.ErrorContext(ctx, fmt.Sprintf("Error while creating row: %v", currentRowIndex))
		return
	}
	return currentRowIndex
}

func (s xlsxService) createCell(ctx context.Context, sheet *xlsx.Sheet, row int, col int, val string, style *xlsx.Style) (cell *xlsx.Cell) {
	return s.createCellWithType(ctx, sheet, row, col, val, style, nil)
}

func (s xlsxService) createCellWithType(ctx context.Context, sheet *xlsx.Sheet, row int, col int, val string, style *xlsx.Style, cellType *model.CellType) (cell *xlsx.Cell) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		s.log.ErrorContext(ctx, fmt.Sprintf("Error while creating cell, row: %v, col: %v, val: %v", row, col, val))
		return
	}
	s.setCellValue(cellType, cell, val)
	cell.SetStyle(style)
	return
}

func (s xlsxService) setCellValue(columnType *model.CellType, cell *xlsx.Cell, val string) {
	if columnType == nil || *columnType == model.StringCell {
		cell.Value = val
	} else if *columnType == model.NumberCell {
		cell.SetNumeric(val)
	}
}

const DefaultColumnsSize = 2

func (s xlsxService) adjustColWidth(ctx context.Context, wbSheet *xlsx.Sheet, sheet *model.Sheet) {
	size := DefaultColumnsSize
	if sheet.Columns != nil && len(*sheet.Columns) > size {
		size = len(*sheet.Columns)
	}
	s.setColAutoWidth(ctx, wbSheet, size)
}

func (s xlsxService) setColAutoWidth(ctx context.Context, wbSheet *xlsx.Sheet, length int) {
	for i := 1; i <= length; i++ {
		err := wbSheet.SetColAutoWidth(i, xlsx.DefaultAutoWidth)
		if err != nil {
			s.log.ErrorContext(ctx, fmt.Sprintf("Error while setting auto width for column: %v", i))
			return
		}
	}
}
