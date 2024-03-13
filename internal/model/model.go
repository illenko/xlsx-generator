package model

type XlsxRequest struct {
	Sheets []Sheet `json:"sheets"`
}

type Sheet struct {
	Name           string          `json:"name"`
	AdditionalInfo *AdditionalInfo `json:"additionalInfo"`
	Columns        *[]Column       `json:"columns"`
	Data           *[]Data         `json:"data"`
}

type AdditionalInfo struct {
	Top    *[]AdditionalData `json:"top"`
	Bottom *[]AdditionalData `json:"bottom"`
}

type AdditionalData struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Column struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Type  *CellType `json:"type"`
	Color *Color    `json:"color"`
}

type Data map[string]string

type Color struct {
	Font       *string `json:"font"`
	Background *string `json:"background"`
}

const (
	NumberCell CellType = "number"
	StringCell CellType = "string"
)

type CellType string
