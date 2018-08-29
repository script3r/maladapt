package model

import (
	"crypto/md5"
	"crypto/sha256"
	"github.com/globalsign/mgo/bson"
)

type Result struct {
	Infected bool `json:"infected" bson:"infected"`
}

type ScanResponse struct {
	Filename  string `json:"filename" bson:"filename"`
	SHA256    string `json:"sha256" bson:"sha256"`
	MD5       string `json:"md5" bson:"md5"`
	Permalink string `json:"permalink" bson:"permalink"`
}

type ScanResult struct {
	ID        bson.ObjectId     `json:"id" bson:"_id"`
	Filename  string            `json:"filename" bson:"filename"`
	SHA256    string            `json:"sha256" bson:"sha256"`
	MD5       string            `json:"md5" bson:"md5"`
	Path      string            `json:"path" bson:"path"`
	Permalink string            `json:"permalink" bson:"permalink"`
	Results   map[string]Result `json:"results" bson:"results"`
}

type ScanMessage struct {
	Filename string
	SHA256   [sha256.Size]byte
	MD5      [md5.Size]byte
	Path     string
}
