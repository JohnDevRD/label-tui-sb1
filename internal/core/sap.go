package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type loginRequest struct {
	CompanyDB string `json:"CompanyDB"`
	UserName  string `json:"UserName"`
	Password  string `json:"Password"`
}

type SapClient struct {
	BaseURL   string
	SessionID string
	PriceList int
	client    *http.Client
}

func NewSapClient(baseURL string, priceList int) *SapClient {
	return &SapClient{
		BaseURL:   strings.TrimRight(baseURL, "/"),
		PriceList: priceList,
		client:    &http.Client{},
	}
}

func (s *SapClient) Login(companyDB, user, password string) error {
	body := loginRequest{
		CompanyDB: companyDB,
		UserName:  user,
		Password:  password,
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", s.BaseURL+"/Login", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed (status %d): %s", resp.StatusCode, string(raw))
	}

	for _, c := range resp.Cookies() {
		if c.Name == "B1SESSION" {
			s.SessionID = c.Value
			return nil
		}
	}

	return fmt.Errorf("no B1SESSION cookie in response")
}

type SapItemsResponse struct {
	Value []Article `json:"value"`
}

func (s *SapClient) QueryArticles(filter string) ([]Article, error) {
	u := s.BaseURL + "/Items?$select=ItemCode,ItemName,BarCode"
	if filter != "" {
		u += "&$filter=" + url.QueryEscape(filter)
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", "B1SESSION="+s.SessionID)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("items query failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("items query failed (status %d): %s", resp.StatusCode, string(raw))
	}

	var result SapItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode items response: %w", err)
	}

	for i := range result.Value {
		result.Value[i].Price = result.Value[i].priceForList(s.PriceList)
	}

	return result.Value, nil
}

func (a Article) priceForList(listNum int) float64 {
	for _, p := range a.ItemPrices {
		if p.PriceList == listNum {
			return p.Price
		}
	}
	return 0
}
