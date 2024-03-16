package service

import (
	"context"
	"encoding/json"
	"github.com/illenko/xlsx-generator/internal/logger"
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/xlsx/v3"
	"os"
	"testing"
)

func TestXlsxServiceImpl_Generate(t *testing.T) {
	cases := []struct {
		name string
		file string
	}{
		{
			name: "One sheet, full data",
			file: "one_sheet_full_data",
		},
		{
			name: "One sheet, no table",
			file: "one_sheet_no_table",
		},
		{
			name: "One sheet, only table",
			file: "one_sheet_only_table",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := readRequest(tc)

			sut := New(logger.New())

			wb := sut.Generate(context.Background(), request)

			assert.Equal(t, len(request.Sheets), len(wb.Sheets))

			for i, sheet := range request.Sheets {
				wbSheet := wb.Sheets[i]

				assert.Equal(t, sheet.Name, wbSheet.Name)

				currentRowIndex := 0

				if sheet.AdditionalInfo != nil && sheet.AdditionalInfo.Top != nil {
					currentRowIndex = assertAdditionalData(t, sheet.AdditionalInfo.Top, wbSheet, currentRowIndex)
				}

				if sheet.Columns != nil {
					currentRowIndex = assertColumns(t, sheet, wbSheet, currentRowIndex)
				}

				if sheet.Data != nil {
					currentRowIndex = assertData(t, sheet, wbSheet, currentRowIndex)
					currentRowIndex++
				}

				if sheet.AdditionalInfo != nil && sheet.AdditionalInfo.Bottom != nil {
					assertAdditionalData(t, sheet.AdditionalInfo.Bottom, wbSheet, currentRowIndex)
				}
			}
		})
	}
}

func readRequest(tc struct {
	name string
	file string
}) model.XlsxRequest {
	file, _ := os.ReadFile("../../test_data/" + tc.file + ".json")

	var request model.XlsxRequest
	_ = json.Unmarshal(file, &request)
	return request
}

func assertAdditionalData(t *testing.T, additionalData *[]model.AdditionalData, wbSheet *xlsx.Sheet, currentRowIndex int) int {
	for _, ad := range *additionalData {
		titleCell, err := wbSheet.Cell(currentRowIndex, 0)
		if err != nil {
			assert.Error(t, err, "Cell does not exist")
		}
		valueCell, err := wbSheet.Cell(currentRowIndex, 1)
		if err != nil {
			assert.Error(t, err, "Cell does not exist")
		}
		assert.Equal(t, ad.Title, titleCell.Value)
		assert.Equal(t, ad.Value, valueCell.Value)
		currentRowIndex++ // empty row
	}
	return currentRowIndex
}

func assertColumns(t *testing.T, sheet model.Sheet, wbSheet *xlsx.Sheet, currentRowIndex int) int {
	for i, c := range *sheet.Columns {
		columnCell, err := wbSheet.Cell(currentRowIndex, i)
		if err != nil {
			assert.Error(t, err, "Cell does not exist")
		}
		assert.Equal(t, c.Title, columnCell.Value)
	}
	currentRowIndex++
	return currentRowIndex
}

func assertData(t *testing.T, sheet model.Sheet, wbSheet *xlsx.Sheet, currentRowIndex int) int {
	for _, d := range *sheet.Data {
		for j, c := range *sheet.Columns {
			dataCell, err := wbSheet.Cell(currentRowIndex, j)
			if err != nil {
				assert.Error(t, err, "Cell does not exist")
			}
			assert.Equal(t, d[c.ID], dataCell.Value)
		}
		currentRowIndex++
	}
	return currentRowIndex

}
