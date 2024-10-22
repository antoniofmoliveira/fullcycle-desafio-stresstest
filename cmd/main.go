package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/db"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/pool"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/report"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/stats"
	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/usecase"
)

func main() {

	endpoint, numtests, requestType, payload := handleFlags()

	rec := make(chan *dto.Red)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := db.NewDB(pool.GetDb(), rec)
	defer database.Close()
	go database.Store(ctx)

	client := pool.GetHttpClient()
	defer client.CloseIdleConnections()

	log.Println("Tests starting...")
	fmt.Println(" Running ", numtests, " tests for endpoint ", endpoint)
	wg := sync.WaitGroup{}
	wg.Add(numtests)

	if requestType == "GET" {
		uGet := usecase.NewHttpGet(client, endpoint, numtests, 1, rec)
		go uGet.ExecuteGet(ctx, &wg)

	} else if requestType == "POST" {
		uPost := usecase.NewHttpPost(client, endpoint, numtests, 2, payload, rec)
		go uPost.ExecutePost(ctx, &wg)
	}

	wg.Wait()

	result := stats.CalculateRed(database.GetAllReds())

	errors := stats.CalculateErrors(database.GetRedWithErrors())

	percentiles := stats.CalculatePercentile(database.GetAllReds())

	report.ReportRed(result)

	report.ReportError(errors)

	report.ReportPercentiles(percentiles)

	log.Println("Tests finished.")

}

func handleFlags() (string, int, string, string) {
	n := flag.Int("numtests", 300000, "Number of tests")
	endpoint := flag.String("endpoint", "http://localhost:8080/hello", "Endpoint")
	requestType := flag.String("requesttype", "GET", "Request type GET or POST")
	payload := flag.String("payload", "{}", "Payload JSON")
	flag.Parse()

	errors := []string{}

	if *n < 1 {
		errors = append(errors, "Number of tests must be greater than 0")
	}

	if *requestType != "GET" && *requestType != "POST" {
		errors = append(errors, "Request type must be GET or POST")
	}
	_, err := json.Marshal(payload)
	if err != nil {
		errors = append(errors, "Payload is not valid JSON")
	}
	err = pool.TestEndpoint(*requestType, *endpoint, *payload)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		for _, v := range errors {
			log.Println(v)
		}
		panic("Invalid flags")
	}

	return *endpoint, *n, *requestType, *payload
}
