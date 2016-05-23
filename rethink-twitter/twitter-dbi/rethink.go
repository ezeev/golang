package twitterDBI

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	r "gopkg.in/dancannon/gorethink.v2"
)

// RethinkDB Implements TwitterDBI interface
type RethinkDB struct {
	session *r.Session
}

// Connect connects to a rethink db server or cluster of servers
func (rt *RethinkDB) Connect(connString string) error {
	var err error
	rt.session, err = r.Connect(r.ConnectOpts{
		Address:       connString,
		DiscoverHosts: true,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	log.Printf("Session connected?: %t", rt.session.IsConnected())
	return nil
}

// CreateTwitterDatabase creates twitter database in rethinkdb
func (rt *RethinkDB) CreateTwitterDatabase() error {
	_, err := r.DBCreate("twitter").RunWrite(rt.session)
	if err != nil {
		log.Print("Unable to create database. Database twitter probably already exists.")
	} else {
		log.Print("Created database twitter.")
	}
	return err
}

// CreateTweetsTable creates tweets table in rethinkdb
func (rt *RethinkDB) CreateTweetsTable() error {
	_, err := r.DB("twitter").TableCreate("tweets").RunWrite(rt.session)
	if err != nil {
		log.Print("Unable to create table. Table twitter.tweets probably already exists.")
	} else {
		log.Print("Created table twitter.tweets")
	}
	return err
}

// SaveTweet Saves a tweet to RethinkDB
func (rt *RethinkDB) SaveTweet(tweet twitter.Tweet) error {
	err := r.DB("twitter").Table("tweets").Insert(tweet).Exec(rt.session)
	return err
}

// Close closes the rethinkdb session
func (rt *RethinkDB) Close() {
	rt.session.Close()
}