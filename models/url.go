//go:generate ffjson $GOFILE
package models

import (
	"shlink-server/pkg/genid"
	"time"

	"github.com/globalsign/mgo/bson"
)

// URL is a document structure for
// the collection url. It is used to
// store url(s).
type URL struct {
	Success   bool       `bson:",omitempty" json:"success,omitempty"`
	Err       string     `bson:",omitempty" json:"err,omitempty"`
	Hash      string     `bson:"hash" json:"-"`
	ID        string     `bson:"id" json:"id,omitempty"`
	TargetURL string     `bson:",omitempty" json:"targeturl,omitempty"`
	LongURL   string     `bson:"longurl" json:"longurl,omitempty"`
	Timestamp *time.Time `bson:"ts" json:"ts,omitempty"`
	Stats     []Stats    `bson:"stats,omitempty" json:"stats,omitempty"`
}

// Stats is a document structure for
// the collection statistics. It is used
// for making statistics.
type Stats struct {
	Clicks int `bson:"clicks" json:"clicks,omitempty"`
}

// FindURL returns document is a structure.
func FindURL(doc bson.M) (resp URL, err error) {
	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(dbName)

	err = db.C(collections["url"]).Find(&doc).One(&resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// InfoURL returns document in an array of structures.
func InfoURL(id string) (resp []URL, err error) {
	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(dbName)

	pipe := db.C(collections["url"]).Pipe([]bson.M{{"$match": bson.M{"id": id}}, {"$lookup": bson.M{"from": collections["statistics"], "localField": "id", "foreignField": "id", "as": "stats"}}})
	err = pipe.All(&resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// InsertURL inserts a document.
func InsertURL(doc interface{}) (err error) {
	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(dbName)

	err = db.C(collections["url"]).Insert(&doc)
	if err != nil {
		return
	}

	return nil
}

// ReadyToInsert fills document.
func (u *URL) ReadyToInsert(hash string, seq *Counter) {
	u.Hash = hash
	u.ID = genid.IntToBase62(seq.Sequence - 1)
	timeStamp := time.Now()
	u.Timestamp = &timeStamp
}

// UpdateStats updates statistics.
func UpdateStats(id string) (err error) {
	newSession := session.Copy()
	defer newSession.Close()

	db := newSession.DB(dbName)

	if _, err = db.C(collections["statistics"]).Upsert(bson.M{"id": id}, bson.M{"$set": bson.M{"id": id}, "$inc": bson.M{"clicks": 1}}); err != nil {
		return
	}

	return nil
}
