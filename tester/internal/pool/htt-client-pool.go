package pool

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetHttpClient() *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   200,
		MaxIdleConns:          200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{Transport: tr}

}

func TestEndpoint(method string, url string, payload string) error {
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		slog.Error("TestEndpoint", "msg", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := GetHttpClient().Do(req)
	if err != nil {
		slog.Error("TestEndpoint Do", "msg", err.Error())
		return err
	}
	defer res.Body.Close()
	return nil
}
