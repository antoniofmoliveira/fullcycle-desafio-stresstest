package entity

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Red struct {
	Target     string
	SentAt     time.Time
	ReceivedAt time.Time
	StatusCode int
	Payload    string
}

func (r *Red) Get(client *http.Client) *Red {
	req, err := http.NewRequest("GET", r.Target, nil)
	if err != nil {
		slog.Error("(*Red).Get", "msg", err.Error())
		panic(err)
	}
	r.SentAt = time.Now()

	res, err := client.Do(req)
	if err != nil {
		r.ReceivedAt = time.Now()
		r.StatusCode = -1
		return r
	}
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		slog.Error("(*Red).Get io.Copy", "msg", err.Error())
	}
	res.Body.Close()

	r.ReceivedAt = time.Now()
	r.StatusCode = res.StatusCode
	return r
}

func (r *Red) Post(client *http.Client) *Red {
	req, err := http.NewRequest("POST", r.Target, strings.NewReader(r.Payload))
	if err != nil {
		slog.Error("(*Red).Post", "msg", err.Error())
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r.SentAt = time.Now()

	res, err := client.Do(req)
	if err != nil {
		r.ReceivedAt = time.Now()
		r.StatusCode = -1
		return r
	}

	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		slog.Error("(*Red).Post io.Copy", "msg", err.Error())
	}
	res.Body.Close()

	r.ReceivedAt = time.Now()
	r.StatusCode = res.StatusCode
	return r
}
