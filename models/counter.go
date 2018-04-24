package models

import (
	"shlink-server/utils"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Counter is a document structure for
// the collection counter. It is used
// for id sequence.
type Counter struct {
	ID       string `bson:"_id"`
	Sequence int    `bson:"sequence"`
}

func FindAndModify(updateQuery bson.M) (doc Counter, err error) {
	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(utils.Conf.Database.DB)

	changes := mgo.Change{
		Update:    updateQuery,
		Upsert:    true,
		ReturnNew: true,
	}

	_, err = db.C(collections["counter"]).Find(bson.M{}).Apply(changes, &doc)

	return
}
