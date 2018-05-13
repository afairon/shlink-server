package models

import (
	"github.com/afairon/shlink-server/utils"

	"github.com/globalsign/mgo"
)

var collections = map[string]string{
	"url":        "url",
	"statistics": "statistics",
	"counter":    "counter",
}

var session *mgo.Session
var err error

// Connect establishes a new MongoDB session.
func Connect(c *utils.Config) (*mgo.Session, error) {
	if c.Database.Host == "" {
		c.Database.Host = "127.0.0.1"
	}
	if c.Database.Port == "" {
		c.Database.Port = "27017"
	}

	session, err = mgo.Dial(c.Database.Host + ":" + c.Database.Port)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// CreateIndexes creates indexes for
// three collections: url, statistics and counter.
func CreateIndexes() {
	newSession := session.Copy()

	defer newSession.Close()

	db := newSession.DB(utils.Conf.Database.DB)

	if err := db.C(collections["url"]).EnsureIndex(mgo.Index{
		Key:    []string{"hash", "id"},
		Unique: true,
	}); err != nil {
		panic(err)
	}

	if err := db.C(collections["statistics"]).EnsureIndex(mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}); err != nil {
		panic(err)
	}

	if err := db.C(collections["counter"]).EnsureIndex(mgo.Index{
		Key:    []string{"_id", "sequence"},
		Unique: true,
	}); err != nil {
		panic(err)
	}
}
