package httpclients

import (
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

type ipHttpClient struct {
	url  string
	http *http.Client
}

type IpInterface interface {
	GetIp() (string, error)
}

type IP struct {
	IP string `json:"ip"`
}

func NewIPHttpClient() IpInterface {
	url := "https://api.ipify.org?format=json"
	return &ipHttpClient{
		url:  url,
		http: &http.Client{},
	}
}

func (hc *ipHttpClient) GetIp() (string, error) {

	resp, err := hc.http.Get(hc.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ip IP

	err = json.Unmarshal(body, &ip)
	if err != nil {
		return "", err
	}
	return ip.IP, nil
}
