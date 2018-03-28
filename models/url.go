//go:generate ffjson $GOFILE
package models

import "time"

// URL basic mongo structure
type URL struct {
	Success   bool      `json:"success" bson:",omitempty"`
	Err       string    `json:"err,omitempty" bson:",omitempty"`
	Hash      string    `json:"-" bson:"hash"`
	ID        string    `json:"id,omitempty" bson:"id"`
	ShortURL  string    `json:"targeturl,omitempty" bson:",omitempty"`
	LongURL   string    `json:"longurl,omitempty" bson:"longurl"`
	Timestamp time.Time `json:"ts,omitempty" bson:"ts"`
	TTL       time.Time `json:"ttl,omitempty" bson:"ttl,omitempty"`
}
