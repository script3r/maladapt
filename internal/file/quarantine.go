package file

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type Quarantiner interface {
	RenderInert(input []byte, filename string) error
	RenderAlive(input []byte, filename string) error
}

type QuarantineAdmin struct {
	Quaratiner Quarantiner
}

func NewQuarantineAdmin(quarantinelocation string, quarantiner Quarantiner) *QuarantineAdmin {
	return &QuarantineAdmin{quarantiner}
}

func (q *QuarantineAdmin) Handle(input io.Reader, filename string, size int64) error {
	logger := log.WithFields(log.Fields{"func": "Handle"})

	contents, err := ioutil.ReadAll(input)
	if err != nil {
		logger.Error(err)
		return err
	}

	//Compute Hashes
	sha256 := q.computeSHA256(contents)
	md5 := q.computeMD5(contents)
	logger.WithFields(log.Fields{
		"filename": filename,
		"sha256":   hex.EncodeToString(sha256[:]),
		"md5":      hex.EncodeToString(md5[:])}).Info()

	//TODO
	//Use Hash to check if we have seen thisfile before.

	//QuarantineFile
	q.Quaratiner.RenderInert(contents, filename)
	return nil

}

func (q *QuarantineAdmin) computeSHA256(input []byte) [sha256.Size]byte {
	return sha256.Sum256(input)
}

func (q *QuarantineAdmin) computeMD5(input []byte) [md5.Size]byte {
	return md5.Sum(input)
}
