//go:generate ffjson $GOFILE
package models

import (
	"time"
)

// URL is a document structure for
// the collection url. It is used to
// store url(s).
type URL struct {
	Success   bool       `bson:",omitempty" json:"success,omitempty"`
	Err       string     `bson:",omitempty" json:"err,omitempty"`
	Hash      string     `bson:"hash" json:"-"`
	ID        string     `bson:"id" json:"id,omitempty"`
	ShortURL  string     `bson:",omitempty" json:"targeturl,omitempty"`
	LongURL   string     `bson:"longurl" json:"longurl,omitempty"`
	Timestamp *time.Time `bson:"ts" json:"ts,omitempty"`
	Stats     []Stats    `bson:"stats,omitempty" json:"stats,omitempty"`
}

type Stats struct {
	Clicks int `bson:"clicks" json:"clicks,omitempty"`
}
