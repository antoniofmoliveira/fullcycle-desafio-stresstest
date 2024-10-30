package pool

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

func TestStressEndpoint(t *testing.T) {
	type args struct {
		method  string
		url     string
		payload string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				method:  "POST",
				url:     "http://localhost:8080/hello",
				payload: `{"message":"World"}`,
			},
			wantErr: false,
		},
		{
			name: "Fail",
			args: args{
				method:  "POST",
				url:     "http://localhost:8080/helloo",
				payload: `{"message":"World"}`,
			},
			wantErr: true,
		},
	}
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StressEndpoint(tt.args.method, tt.args.url, tt.args.payload); (err != nil) != tt.wantErr {

				t.Errorf("StressEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
