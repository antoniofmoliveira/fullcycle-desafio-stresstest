package pool

import (
	"net/http"
	"strings"
	"time"
)

func GetHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 100,
	}
	return &http.Client{Transport: tr}
}

func TestEndpoint(method string, url string, payload string) error {
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := GetHttpClient().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
