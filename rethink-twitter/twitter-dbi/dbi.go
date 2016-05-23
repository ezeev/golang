package twitterDBI

import "github.com/dghubble/go-twitter/twitter"

// TwitterDBI Interface for database access
type TwitterDBI interface {
	Connect(string) error
	Close()
	CreateTweetsTable() error
	CreateTwitterDatabase() error
	SaveTweet(twitter.Tweet) error
}

// NewDBI Returns a new instance of a DBI
func NewDBI(dbms string) TwitterDBI {
	if dbms == "rethinkdb" {
		return &RethinkDB{}
	}
	return nil
}
