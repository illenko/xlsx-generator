package model

type XlsxRequest struct {
	Sheets []Sheet `json:"sheets"`
}

type Sheet struct {
	Name           string          `json:"name" example:"Payments"`
	AdditionalInfo *AdditionalInfo `json:"additionalInfo"`
	Columns        *[]Column       `json:"columns"`
	Data           *[]Data         `json:"data"`
}

type AdditionalInfo struct {
	Top    *[]AdditionalData `json:"top"`
	Bottom *[]AdditionalData `json:"bottom"`
}

type AdditionalData struct {
	Title string `json:"title" example:"User"`
	Value string `json:"value" example:"user@test.com"`
}

type Column struct {
	ID    string    `json:"id" example:"payment_id"`
	Title string    `json:"title" example:"Payment ID"`
	Type  *CellType `json:"type" example:"number"`
	Color *Color    `json:"color"`
}

type Data map[string]string

type Color struct {
	Font       *string `json:"font" example:"1D1E1A"`
	Background *string `json:"background" example:"C4DC8F"`
}

const (
	NumberCell CellType = "number"
	StringCell CellType = "string"
)

type CellType string
