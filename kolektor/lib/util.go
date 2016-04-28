package kolektor

import (
	"sort"
)

func timerStats(arr []float64) map[string]float64 {

	sort.Float64s(arr)

	stats := make(map[string]float64)

	//average
	var total float64
	for _, value := range arr {
		total += value
	}
	avg := total / float64(len(arr))
	stats["avg"] = avg

	//median
	n := len(arr)
	if len(arr)%2 == 0 { //even
		stats["med"] = (arr[n/2] + arr[(n/2)-1]) / 2
	} else { //odd
		stats["med"] = arr[(n-1)/2]
	}
	//min
	stats["min"] = arr[0]
	//max
	stats["max"] = arr[n-1]
	//count
	stats["count"] = float64(n)

	return stats
}
