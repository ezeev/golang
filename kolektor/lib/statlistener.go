package kolektor

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ezeev/golang/kolektor/output"
	"github.com/golang/glog"
)

//Stat  blah
type Stat struct {
	key       string
	statType  string
	value     float64
	increment bool
}

func NewStat() *Stat {
	return &Stat{}
}

func ParseStat(input string) *Stat {
	S := NewStat()
	str := strings.Split(input, "|")
	strKeyAndVal := strings.Split(str[0], ":")
	S.key = strKeyAndVal[0]
	S.statType = str[1]
	if strings.Contains(strKeyAndVal[1], "+") || strings.Contains(strKeyAndVal[1], "-") {
		// this means we're incrementing a gauge
		S.increment = true
	} else {
		S.increment = false
	}
	S.value, _ = strconv.ParseFloat(strKeyAndVal[1], 64)
	return S

}

//StatCollector struct
type StatListener struct {
	gauges    map[string]float64
	counters  map[string]float64
	timers    map[string][]float64
	lastFlush time.Time
}

func NewStatListener() *StatListener {
	return &StatListener{
		gauges:   make(map[string]float64),
		counters: make(map[string]float64),
		timers:   make(map[string][]float64),
	}
}

func (s *StatListener) FlushStats(flushInterval float64, be output.Backend) {

	//start timer
	s.lastFlush = time.Now()

	for {
		t := time.Now()
		tUnix := t.Unix()
		d := t.Sub(s.lastFlush)
		//fmt.Println(d.Seconds())
		if d.Seconds() < flushInterval {
			//don't check every single millisecond.
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("Flushing stats, last flush was at ", s.lastFlush)

		stats := make(map[string]float64)
		stats2 := make([]output.Metric, 0)

		//timers
		var key string
		for k := range s.timers {
			key = "timers." + k
			timerStats := timerStats(s.timers[k])

			stats2 = append(stats2, output.Metric{Name: key + ".avg", Value: timerStats["avg"], Timestamp: tUnix})
			stats2 = append(stats2, output.Metric{Name: key + ".med", Value: timerStats["med"], Timestamp: tUnix})
			stats2 = append(stats2, output.Metric{Name: key + ".min", Value: timerStats["min"], Timestamp: tUnix})
			stats2 = append(stats2, output.Metric{Name: key + ".max", Value: timerStats["max"], Timestamp: tUnix})
			stats2 = append(stats2, output.Metric{Name: key + ".count", Value: timerStats["count"], Timestamp: tUnix})

		}
		//reset timers for next flush
		s.timers = make(map[string][]float64)
		//gauges
		for k := range s.gauges {
			key = "gauges." + k
			stats2 = append(stats2, output.Metric{Name: key, Value: s.gauges[k], Timestamp: tUnix})
		}
		//counters
		for k := range s.counters {
			key = "counters." + k
			stats[key] = s.counters[k]
			stats2 = append(stats2, output.Metric{Name: key, Value: s.counters[k], Timestamp: tUnix})
			//now reset to 0 before next flush
			s.counters[k] = 0
		}

		//flush on it's own thread
		go be.Flush(stats2)

		s.lastFlush = t
	}
}

func (s *StatListener) AddGauge(stat *Stat) {
	s.gauges[stat.key] = stat.value
}
func (s *StatListener) IncrementGauge(stat *Stat) {
	s.gauges[stat.key] += stat.value
}
func (s *StatListener) AddCounter(stat *Stat) {
	s.counters[stat.key] += stat.value
}
func (s *StatListener) AddTimer(stat *Stat) {
	s.timers[stat.key] = append(s.timers[stat.key], stat.value)
}

//Lock comments
func (s *StatListener) Lock() {

}

//UnLock This is some text
func (s *StatListener) UnLock() {

}

//ListenForStats listens for stats
func (s *StatListener) ListenForStats(port string) {

	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		glog.Fatalf("Error resolving UDP Address: %s", err)
	}
	Conn, err := net.ListenUDP("udp", ServerAddr)
	defer Conn.Close()
	if err != nil {
		glog.Fatalf("Error listening on UDP: %s", err)
	}
	buf := make([]byte, 1024)
	for {
		n, addr, err := Conn.ReadFromUDP(buf)
		stat := ParseStat(string(buf[0 : n-1]))
		if stat.statType == "g" {
			if stat.increment == false {
				s.AddGauge(stat)
			} else {
				s.IncrementGauge(stat)
			}
		} else if stat.statType == "c" {
			s.AddCounter(stat)
		} else if stat.statType == "ms" {
			s.AddTimer(stat)
		}
		if err != nil {
			fmt.Println("Error: ", err, "from ", addr)
		}
	}
}
