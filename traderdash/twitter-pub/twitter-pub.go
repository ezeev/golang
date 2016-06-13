package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/nats-io/nats"
	"github.com/rcrowley/go-metrics"
)

var tweetCount metrics.Counter

func main() {

	tweetCount = metrics.NewCounter()
	metrics.Register("tweetCount", tweetCount)
	go metrics.Log(metrics.DefaultRegistry, 10*time.Second, log.New(os.Stderr, "twitter-pub metrics: ", log.Lmicroseconds))

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	tracks := flags.String("tracks", "", "Comma separated list of tracks to search")
	natsUrls := flags.String("nats-urls", "", "Comma separated list of nats message queue servers")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" || *tracks == "" || *natsUrls == "" {
		log.Fatal("Consumer key/secret and Access token/secret and tracks required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	//Connect to NATS services
	nc, _ := nats.Connect(*natsUrls)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	subj := "tweets"

	demux.Tweet = func(tweet *twitter.Tweet) {
		//nc.Publish(subj, []byte(tweet.Text))
		c.Publish(subj, tweet)
		nc.Flush()
		tweetCount.Inc(1)
		//fmt.Println(tweet.Text)
	}

	log.Println("Starting Stream...")

	// FILTER
	arrTracks := strings.Split(*tracks, ",")

	filterParams := &twitter.StreamFilterParams{
		Track:         arrTracks,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()

}
