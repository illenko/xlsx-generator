package style

import (
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/tealeg/xlsx/v3"
)

var Default = func() *xlsx.Style {
	style := xlsx.NewStyle()
	style.Font.Size = 8
	return style
}()

var DefaultColumn = func() *xlsx.Style {
	style := xlsx.NewStyle()
	style.Fill.FgColor = "A1AAAD"
	style.Fill.PatternType = "solid"
	style.Font.Size = 8
	style.Font.Bold = true
	return style
}()

var DefaultAdditionalInfoTitle = func() *xlsx.Style {
	style := xlsx.NewStyle()
	style.Font.Size = 8
	style.Font.Bold = true
	return style
}()

func Resolve(color *model.Color, defaultStyle *xlsx.Style) *xlsx.Style {
	if color == nil {
		return defaultStyle
	}

	style := *defaultStyle

	if color.Background != nil {
		style.Fill.FgColor = *color.Background
	}

	if color.Font != nil {
		style.Font.Color = *color.Font
	}

	return &style
}
