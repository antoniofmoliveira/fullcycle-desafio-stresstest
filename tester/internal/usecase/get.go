package usecase

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/entity"
)

type httpGet struct {
	Client              *http.Client
	Target              string
	ReturnChannel       chan *dto.Red
	NumRequests         int
	IntervalNanoseconds int
}

func NewHttpGet(client *http.Client, target string, numRequests int, interval int, rec chan *dto.Red) *httpGet {
	return &httpGet{
		Client:              client,
		Target:              target,
		ReturnChannel:       rec,
		NumRequests:         numRequests,
		IntervalNanoseconds: interval,
	}
}

func (h *httpGet) ExecuteGet(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < h.NumRequests; i++ {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			if h.IntervalNanoseconds > 0 {
				time.Sleep(time.Duration(h.IntervalNanoseconds) * time.Nanosecond)
			}
			go func(client *http.Client, target string, rec chan *dto.Red, wg *sync.WaitGroup) {
				r := &entity.Red{
					Target: target,
				}
				r.Get(client)
				dto := &dto.Red{Target: r.Target, SentAt: r.SentAt, ReceivedAt: r.ReceivedAt, StatusCode: r.StatusCode, Duration: r.ReceivedAt.Sub(r.SentAt)}
				rec <- dto
				wg.Done()
			}(h.Client, h.Target, h.ReturnChannel, wg)
		}
	}
}
