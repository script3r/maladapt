package model

import (
	"crypto/md5"
	"crypto/sha256"
	"github.com/globalsign/mgo/bson"
)

type Result struct {
	Infected bool `json:"infected" bson:"infected"`
}

type ScanResult struct {
	ID        bson.ObjectId     `bson:"_id"`
	Filename  string            `bson:"filename"`
	SHA256    [sha256.Size]byte `bson:"sha256"`
	MD5       [md5.Size]byte    `bson:"md5"`
	Path      string            `bson:"path"`
	Permalink string            `bson:"permalink"`
	Results   map[string]Result `bson:"results"`
}
