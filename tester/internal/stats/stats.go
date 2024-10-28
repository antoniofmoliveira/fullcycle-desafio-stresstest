package stats

import (
	"fmt"
	"sort"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
)

func CalculateRed(recs []*dto.Red) map[string]*dto.ResultRed {
	mapRecs := make(map[string]*dto.ResultRed)
	for _, rec := range recs {
		s := fmt.Sprintf("%d", rec.SentAt.Minute()*100+rec.SentAt.Second())
		if mapRecs[s] == nil {
			mapRecs[s] = &dto.ResultRed{}
		}
		mapRecs[s].NumRequestPerSecond++
		if rec.StatusCode == -1 {
			mapRecs[s].NumNetworkErrorPerSecond++
		}
		if rec.StatusCode != 200 && rec.StatusCode != -1 {
			mapRecs[s].NumRequestWithErrorPerSecond++
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
	return mapRecs
}

func CalculateErrors(recs []*dto.Red) map[int]*dto.ResultError {

	mapErrs := make(map[int]*dto.ResultError)
	for _, rec := range recs {
		if mapErrs[rec.StatusCode] == nil {
			mapErrs[rec.StatusCode] = &dto.ResultError{ErrorType: rec.StatusCode, NumRequestWithErrorPerSecond: 0}
		}
		mapErrs[rec.StatusCode].NumRequestWithErrorPerSecond++
	}
	return mapErrs
}

func CalculatePercentile(recs []*dto.Red) dto.Percentiles {
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Duration < recs[j].Duration
	})
	l := len(recs)

	p10 := int(float64(l) * 10 / 100)
	p25 := int(float64(l) * 25 / 100)
	p50 := int(float64(l) * 50 / 100)
	p75 := int(float64(l) * 75 / 100)
	p90 := int(float64(l) * 90 / 100)
	p99 := int(float64(l) * 99 / 100)

	per := dto.Percentiles{
		P10: recs[p10].Duration,
		P25: recs[p25].Duration,
		P50: recs[p50].Duration,
		P75: recs[p75].Duration,
		P90: recs[p90].Duration,
		P99: recs[p99].Duration,
	}
	return per

}
