package report

import (
	"fmt"
	"sort"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ReportRed(result map[string]*dto.ResultRed) {
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
}

func ReportError(errors map[int]*dto.ResultError) {
	keys := make([]int, 0, len(errors))
	for k := range errors {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	p := message.NewPrinter(language.English)
	fmt.Printf("\n%-7s\t%10s\n", "Status", "#Errors")
	for _, v := range keys {
		fmt.Printf("%-7s\t%10s\n", p.Sprintf("%d", v), p.Sprintf("%d", errors[v].NumRequestWithErrorPerSecond))
	}
}

func ReportPercentiles(perc dto.Percentiles) {
	fmt.Printf("\n%-10s\t%10s\n", "Percentile", "Duration")
	fmt.Printf("%-10s\t%10v\n", "P10", perc.P10)
	fmt.Printf("%-10s\t%10v\n", "P25", perc.P25)
	fmt.Printf("%-10s\t%10v\n", "P50", perc.P50)
	fmt.Printf("%-10s\t%10v\n", "P75", perc.P75)
	fmt.Printf("%-10s\t%10v\n", "P90", perc.P90)
	fmt.Printf("%-10s\t%10v\n", "P99", perc.P99)
}
