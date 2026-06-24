package core

type Article struct {
	Code        string  `json:"Code"`
	Description string  `json:"ItemName"`
	Barcode     string  `json:"BarCode"`
	Price       float64 `json:"UnitPrice"`
}

type PrintJob struct {
	Article  Article
	Quantity int
	RawZPL   string
}

type Settings struct {
	CompanyDB          string `json:"company_db"`
	SAPServiceLayerURL string `json:"sap_service_layer_url"`
	USBPort            string `json:"usb_port"`
	DefaultTemplate    string `json:"default_template"`
}

type Session struct {
	Cookie  string
	BaseURL string
}
