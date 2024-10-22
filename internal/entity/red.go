package entity

import (
	"io"
	"log"
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
		panic(err)
	}
	r.SentAt = time.Now()
	var res *http.Response
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Microsecond)
		res, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		r.ReceivedAt = time.Now()
		r.StatusCode = -1
		return r
	}
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		log.Println(err)
	}
	res.Body.Close()

	r.ReceivedAt = time.Now()
	r.StatusCode = res.StatusCode
	return r
}

func (r *Red) Post(client *http.Client) *Red {
	req, err := http.NewRequest("POST", r.Target, strings.NewReader(r.Payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r.SentAt = time.Now()
	var res *http.Response
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Microsecond)
		res, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		r.ReceivedAt = time.Now()
		r.StatusCode = -1
		return r
	}

	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		log.Println(err)
	}
	res.Body.Close()
	
	r.ReceivedAt = time.Now()
	r.StatusCode = res.StatusCode
	return r
}
