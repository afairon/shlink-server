package models

// Counter is a document structure for
// the collection counter. It is used
// for id sequence.
type Counter struct {
	ID       string `bson:"_id"`
	Sequence int    `bson:"sequence"`
}
