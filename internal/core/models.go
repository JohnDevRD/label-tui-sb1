package core

type Article struct {
	ItemCode    string      `json:"ItemCode"`
	Description string      `json:"ItemName"`
	Barcode     string      `json:"BarCode"`
	Price       float64     `json:"-"`
	ItemPrices  []ItemPrice `json:"ItemPrices"`
}

type ItemPrice struct {
	PriceList int     `json:"PriceList"`
	Price     float64 `json:"Price"`
	Currency  string  `json:"Currency"`
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
	PriceList          int    `json:"price_list"`
}

type Session struct {
	Cookie  string
	BaseURL string
}
