package traderDB

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// TwitterDBI Interface for database access
type DB interface {
	Connect(string) error
	Close()

	//Tweets
	CreateTweetsTable() error
	CreateTwitterDatabase() error
	SaveTweet(twitter.Tweet) error
	GetTweets() []TweetItem
	StreamTweets(chan TweetItem)

	//Quotes
}

// QuoteItem model representing a stock quote
type QuoteItem struct {
	Id        string
	Symbol    string
	Price     float64
	Timestamp time.Time
}

// TweetItem model representing a tweet
type TweetItem struct {
	Id      string `gorethink:"id,omitempty"`
	Text    string
	Created time.Time
	UserId  string
}

// NewDBI Returns a new instance of a DBI
func NewDBI(dbms string) DB {
	if dbms == "rethinkdb" {
		return &RethinkDB{}
	}
	return nil
}
