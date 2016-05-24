package twitterDBI

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// TwitterDBI Interface for database access
type TwitterDBI interface {
	Connect(string) error
	Close()
	CreateTweetsTable() error
	CreateTwitterDatabase() error
	SaveTweet(twitter.Tweet) error
	GetTweets() []TweetItem
	StreamTweets(chan TweetItem)
}

// Model representing a tweet
type TweetItem struct {
	Id      string `gorethink:"id,omitempty"`
	Text    string
	Created time.Time
	UserId  string
}

// NewDBI Returns a new instance of a DBI
func NewDBI(dbms string) TwitterDBI {
	if dbms == "rethinkdb" {
		return &RethinkDB{}
	}
	return nil
}
