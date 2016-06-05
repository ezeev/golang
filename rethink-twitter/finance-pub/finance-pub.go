package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/pkg/flagutil"
)

type QuoteResponse struct {
	List struct {
		Meta struct {
			Type  string `json:"type"`
			Start int    `json:"start"`
			Count int    `json:"count"`
		} `json:"meta"`
		Resources []struct {
			Resource struct {
				Classname string            `json:"classname"`
				Fields    map[string]string `json:"fields"`
			} `json:"resource"`
		} `json:"resources"`
	} `json:"list"`
}

func saveQuotes(symbols string) {
	url := "http://finance.yahoo.com/webservice/v1/symbols/" + symbols + "/quote?format=json&view=detail"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var respData QuoteResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&respData)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(respData.List.Meta)
	fmt.Println(respData.List.Resources[0].Resource.Fields["symbol"])
}

func main() {

	flags := flag.NewFlagSet("finance", flag.ExitOnError)
	symbols := flags.String("symbols", "", "Comma separated list of symbols")
	interval := flags.Int("interval", 10, "Interval to poll finance API")
	natsUrls := flags.String("nats-urls", "", "Comma separated list of nats message queue servers")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *symbols == "" || *natsUrls == "" {
		log.Fatal("Symbols required")
	}

	//Finance API URL - http://finance.yahoo.com/webservice/v1/symbols/YHOO,AAPL/quote?format=json&view=detail

	//lastPoll := time.Now()
	tick := time.Tick(1 * time.Second)
	tickCount := 0
	for {
		select {
		case <-tick:
			tickCount++
			if tickCount == *interval {
				go saveQuotes(*symbols)
				tickCount = 0
			}
		}
	}

	//Connect to NATS services
	/*
		nc, _ := nats.Connect(*natsUrls)
		c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		subj := "quotes"
	*/
	//c.Publish(subj, tweet)
	//nc.Flush()
	//fmt.Println(tweet.Text)

}
