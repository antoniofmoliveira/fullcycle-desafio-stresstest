package usecase

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/entity"
)

type HttpPost struct {
	Client        *http.Client
	Target        string
	ReturnChannel chan *dto.Red
	NumRequests   int
	Interval      int
	Payload       string
}

func NewHttpPost(client *http.Client, target string, numRequests int, interval int, payload string, rec chan *dto.Red) *HttpPost {
	return &HttpPost{
		Client:        client,
		Target:        target,
		ReturnChannel: rec,
		NumRequests:   numRequests,
		Interval:      interval,
		Payload:       payload,
	}
}

func (h *HttpPost) ExecutePost(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < h.NumRequests; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			if h.Interval > 0 {
				time.Sleep(time.Duration(h.Interval) * time.Nanosecond)
			}
			go func(client *http.Client, target string, rec chan *dto.Red, wg *sync.WaitGroup) {
				r := &entity.Red{
					Target:  target,
					Payload: h.Payload,
				}
				r.Post(client)
				dto := &dto.Red{Target: r.Target, SentAt: r.SentAt, ReceivedAt: r.ReceivedAt, StatusCode: r.StatusCode, Duration: r.ReceivedAt.Sub(r.SentAt)}
				rec <- dto
				wg.Done()
			}(h.Client, h.Target, h.ReturnChannel, wg)
		}
	}
}
