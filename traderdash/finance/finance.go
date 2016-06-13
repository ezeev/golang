package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coreos/pkg/flagutil"
	"github.com/ezeev/golang/traderdash/dbi"
)

var db traderDB.DB

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
	/*
		type QuoteItem struct {
			Id        string
			Symbol    string
			Price     float64
			Volume	  int64
			Timestamp int64
		}
	*/
	for v := range respData.List.Resources {
		q := respData.List.Resources[v].Resource.Fields
		symbol := q["symbol"]
		price, _ := strconv.ParseFloat(q["price"], 64)
		volume, _ := strconv.ParseInt(q["volume"], 20, 64)
		timestamp, _ := strconv.ParseInt(q["ts"], 10, 64)
		//id := symbol + string(timestamp)

		item := traderDB.QuoteItem{
			Symbol:    symbol,
			Price:     price,
			Volume:    volume,
			Timestamp: timestamp,
		}

		err := db.SaveQuote(item)
		if err != nil {
			log.Fatalf("Error saving quote, error: %s", err)
		}
	}
}

func main() {

	flags := flag.NewFlagSet("finance", flag.ExitOnError)
	symbols := flags.String("symbols", "", "Comma separated list of symbols")
	interval := flags.Int("interval", 10, "Interval to poll finance API")
	dbType := flags.String("db-type", "", "Database type to use")
	connStr := flags.String("conn-str", "", "Connection string to database")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "FINANCE")

	if *symbols == "" || *dbType == "" {
		log.Fatal("Symbols required")
	}

	db = traderDB.NewDBI(*dbType)
	db.Connect(*connStr)
	defer db.Close()

	db.CreateQuotesTable()

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
