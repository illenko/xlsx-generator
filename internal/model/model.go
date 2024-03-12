package model

type XlsxRequest struct {
	Sheets []Sheet `json:"sheets"`
}

type Sheet struct {
	Name           string          `json:"name"`
	Header         *Header         `json:"header"`
	AdditionalInfo *AdditionalInfo `json:"additionalInfo"`
	Columns        *[]Column       `json:"columns"`
	Data           *[]Data         `json:"data"`
}

type Header struct {
	Title  string `json:"title"`
	Colors *Color `json:"color"`
}

type Color struct {
	Font       string `json:"font"`
	Background string `json:"background"`
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
	ID    string `json:"id"`
	Title string `json:"title"`
	Width int    `json:"width"`
	Color string `json:"color"`
	Type  string `json:"type"`
}

type Data map[string]string
