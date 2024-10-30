package entity

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/pool"
)

func TestRed_Get(t *testing.T) {
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		r    *Red
		args args
		want *Red
	}{
		{
			name: "Success",
			r: &Red{
				Target:     "http://localhost:8080/hello",
				StatusCode: -1,
				SentAt:     time.Now(),
				ReceivedAt: time.Now(),
			},
			args: args{
				client: pool.GetHttpClient(),
			},
			want: &Red{
				Target:     "http://localhost:8080/hello",
				StatusCode: 200,
				SentAt:     time.Now(),
				ReceivedAt: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		server := &http.Server{Addr: ":8080"}
		defer server.Close()
		http.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello, World!"))
		}))
		go server.ListenAndServe()

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Get(tt.args.client); got.StatusCode != tt.want.StatusCode {
				t.Errorf("Red.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRed_Post(t *testing.T) {
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		r    *Red
		args args
		want *Red
	}{
		{
			name: "Success",
			r: &Red{
				Target:     "http://localhost:8080/hello",
				StatusCode: -1,
				SentAt:     time.Now(),
				ReceivedAt: time.Now(),
				Payload:    `{"message":"World"}`,
			},
			args: args{
				client: pool.GetHttpClient(),
			},
			want: &Red{
				Target:     "http://localhost:8080/hello",
				StatusCode: 200,
				SentAt:     time.Now(),
				ReceivedAt: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		server := &http.Server{Addr: ":8080"}
		defer server.Close()
		http.HandleFunc("POST /hello", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Can't read body", "error", err)
				http.Error(w, "Can't read body", http.StatusInternalServerError)
				return
			}
			var message struct {
				Message string `json:"message"`
			}
			err = json.Unmarshal(body, &message)
			if err != nil {
				slog.Error("Can't unmarshal json", "error", err)
				http.Error(w, "Can't unmarshal json", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Hello, %s!", message.Message)
		})
		go server.ListenAndServe()

		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Post(tt.args.client); got.StatusCode != tt.want.StatusCode {
				t.Errorf("Red.Post() = %v, want %v", got, tt.want)
			}
		})
	}
}
