package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Red struct {
	Target     string
	SentAt     time.Time
	ReceivedAt time.Time
	StatusCode int
}

var client *http.Client

// init sets up the http client to use a Transport with a large number of
// idle connections and a long idle connection timeout. This is so that the
// client can handle a large number of concurrent requests without having to
// do a lot of extraneous work.
func init() {
	tr := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 100,
	}
	client = &http.Client{Transport: tr}
}

// Get performs a GET request to the URL in r.Target, retrying up to three
// times if there is an error. It records the time at which the request was
// sent and the time at which the response was received. It also records the
// status code of the response. If any of the requests fail due to an error,
// the final status code is -1.
func (r *Red) Get() *Red {
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
	defer res.Body.Close()
	r.ReceivedAt = time.Now()
	r.StatusCode = res.StatusCode
	return r
}

// RedRoutine runs a single request to the target URL and sends the result
// down the channel. It retries up to three times if there is an error.
func RedRoutine(target string, rec chan *Red) {
	r := &Red{
		Target: target,
	}
	rec <- r.Get()
}

type ResultRED struct {
	NumRequestPerSecond          int
	NumRequestWithErrorPerSecond int
	NumNetworkErrorPerSecond     int
	AverageDuration              time.Duration
	MaxDuration                  time.Duration
	MinDuration                  time.Duration
}

type ResultError struct {
	ErrorType                    int
	NumRequestWithErrorPerSecond int
}

// calculateRed takes a slice of Red structs and returns a map where the keys
// are strings representing the minute and second the request was sent, and the
// values are Results which contain the total number of requests sent in that
// second, the number of requests that returned an error, and the average
// duration of the requests in that second.
func calculateRed(recs []*Red) (map[string]*ResultRED, map[int]*ResultError) {

	mapRecs := make(map[string]*ResultRED)
	mapErrs := make(map[int]*ResultError)

	for _, rec := range recs {
		s := fmt.Sprintf("%d", rec.SentAt.Minute()*100+rec.SentAt.Second())
		if mapRecs[s] == nil {
			mapRecs[s] = &ResultRED{}
		}
		mapRecs[s].NumRequestPerSecond++
		if rec.StatusCode == -1 {
			mapRecs[s].NumNetworkErrorPerSecond++
		}
		if rec.StatusCode != 200 && rec.StatusCode != -1 {
			mapRecs[s].NumRequestWithErrorPerSecond++

			if mapErrs[rec.StatusCode] == nil {
				mapErrs[rec.StatusCode] = &ResultError{ErrorType: rec.StatusCode, NumRequestWithErrorPerSecond: 0}
			}
			mapErrs[rec.StatusCode].NumRequestWithErrorPerSecond++
		}
		mapRecs[s].AverageDuration = mapRecs[s].AverageDuration + rec.ReceivedAt.Sub(rec.SentAt)

		if rec.ReceivedAt.Sub(rec.SentAt) > mapRecs[s].MaxDuration {
			mapRecs[s].MaxDuration = rec.ReceivedAt.Sub(rec.SentAt)
		}
		if rec.ReceivedAt.Sub(rec.SentAt) < mapRecs[s].MinDuration && mapRecs[s].MinDuration != 0 {
			mapRecs[s].MinDuration = rec.ReceivedAt.Sub(rec.SentAt)
		} else if mapRecs[s].MinDuration == 0 {
            mapRecs[s].MinDuration = rec.ReceivedAt.Sub(rec.SentAt)
        }
	}
	return mapRecs, mapErrs
}

// main is the entry point for the CLI application. It parses command-line flags
// to determine the number of tests, the endpoint to test, and the interval
// between requests. It initiates a series of requests to the specified endpoint
// using goroutines, collects the results, and processes them to calculate RED
// metrics. Finally, it prints a summary of these metrics, showing the number of
// requests per second, errors, average request duration, and network errors.
func main() {

	n := flag.Int("numtests", 300000, "Number of tests")
	endpoint := flag.String("endpoint", "http://localhost:8080/hello", "Endpoint")
	interval := flag.Int("interval", 1, "Interval in microseconds")
	flag.Parse()

	numtests := *n
	inter := *interval
	recs := []*Red{}

	log.Println("Tests starting...")
	fmt.Println(" Running ", numtests, " tests with interval ", inter, " microseconds and endpoint ", *endpoint)

	rec := make(chan *Red)
	for i := 0; i < numtests; i++ {
		if inter > 0 {
			time.Sleep(time.Duration(inter) * time.Microsecond)
		}
		go RedRoutine(*endpoint, rec)
	}

	for i := 0; i < numtests; i++ {
		recs = append(recs, <-rec)
	}

	result, errors := calculateRed(recs)

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	p := message.NewPrinter(language.English)
	fmt.Printf("%-7s\t%10s\t%10s\t%10s\t%10s\t%10s\t%10s\n", "Min/Seg", "Rate", "Error", "Avg Time", "Min Time", "Max Time", "Net Error")
	for _, v := range keys {
		fmt.Printf("%-7s\t%10s\t%10s\t%10v\t%10v\t%10v\t%10s\n", v, p.Sprintf("%d", result[v].NumRequestPerSecond), p.Sprintf("%d", result[v].NumRequestWithErrorPerSecond), result[v].AverageDuration/time.Duration(result[v].NumRequestPerSecond-result[v].NumNetworkErrorPerSecond+1), result[v].MinDuration, result[v].MaxDuration, p.Sprintf("%d", result[v].NumNetworkErrorPerSecond))
	}

	keys2 := make([]int, 0, len(errors))
	for k := range errors {
		keys2 = append(keys2, k)
	}
	sort.Ints(keys2)
	fmt.Printf("\n%-7s\t%10s\n", "Error", "# Errors")
	for _, v := range keys2 {
		fmt.Printf("%-7s\t%10s\n", p.Sprintf("%d", v), p.Sprintf("%d", errors[v].NumRequestWithErrorPerSecond))
	}

	log.Println("Tests finished.")

}
