package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ezeev/golang/rethink-twitter/dbi"
	"github.com/nats-io/nats"
	r "gopkg.in/dancannon/gorethink.v2"
)

var session *r.Session

func printTweets(ch chan traderDB.TweetItem) {
	for v := range ch {
		fmt.Println(v)
	}

}

func main() {

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	natsUrls := flags.String("nats-urls", "", "Comma separated list of nats message queue servers")
	dbType := flags.String("db-type", "", "The DBMS where tweets will be saved")
	connStr := flags.String("conn-str", "", "The DBMS connection string")
	flags.Parse(os.Args[1:])
	if *natsUrls == "" {
		log.Fatal("natsUrls required")
	}
	if *dbType == "" {
		log.Fatal("db-type required")
	}
	if *connStr == "" {
		log.Fatal("conn-str required")
	}

	nc, _ := nats.Connect(*natsUrls)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	db := traderDB.NewDBI(*dbType)
	db.Connect(*connStr)
	db.CreateTwitterDatabase()
	db.CreateTweetsTable()

	//testing
	ch := make(chan traderDB.TweetItem)
	go db.StreamTweets(ch)
	//go db.ReceiveTweets(ch)
	go printTweets(ch)
	//fmt.Println(tweets)

	//nc, _ := nats.Connect(nats.DefaultURL)
	subj := "tweets"
	//nc.Subscribe(subj, func(msg *nats.Msg) {
	c.Subscribe(subj, func(tweet *twitter.Tweet) {
		err := db.SaveTweet(*tweet)
		//bytesGauge.Update(syscall.Getrusage())
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(tweet.Text)
	})

	runtime.Goexit()

}
