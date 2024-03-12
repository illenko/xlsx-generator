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
	Colors struct {
		Font       string `json:"font"`
		Background string `json:"background"`
	} `json:"colors"`
}

type AdditionalInfo struct {
	Top *[]struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"top"`
	Bottom *[]struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"bottom"`
}

type Column struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Width int    `json:"width"`
	Color string `json:"color"`
	Type  string `json:"type"`
}

type Data map[string]interface{}
