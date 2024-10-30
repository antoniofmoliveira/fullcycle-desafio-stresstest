package usecase

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/pool"
)

func TestNewHttpGet(t *testing.T) {
	type args struct {
		client      *http.Client
		target      string
		numRequests int
		interval    int
		rec         chan *dto.Red
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Success",
			args: args{
				client:      pool.GetHttpClient(),
				target:      "http://localhost:8080/hello",
				numRequests: 2,
				interval:    500,
				rec:         make(chan *dto.Red),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := &http.Server{Addr: ":8080"}
			defer server.Close()
			http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello, World!"))
			})
			go server.ListenAndServe()
			got := NewHttpGet(tt.args.client, tt.args.target, tt.args.numRequests, tt.args.interval, tt.args.rec)
			wg := &sync.WaitGroup{}
			wg.Add(tt.args.numRequests)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			go got.ExecuteGet(ctx, wg)

			count := 0
			for range tt.args.numRequests {
				z := <-tt.args.rec
				count++
				slog.Info("Waiting for results...", "count", count, "result", z)
			}
			if got.NumRequests != tt.want {
				t.Errorf("NewHttpGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
