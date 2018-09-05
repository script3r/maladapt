package quarantine

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/pkg/message/rabbit"
	"io"
	"io/ioutil"
)

type ScanResponse struct {
	Filename  string `json:"filename" bson:"filename"`
	SHA256    string `json:"sha256" bson:"sha256"`
	MD5       string `json:"md5" bson:"md5"`
	Permalink string `json:"permalink" bson:"permalink"`
}

type MaladaptFileMeta struct {
}

type Quarantiner interface {
	RenderInert(input []byte, filename string) (string, error)
	RenderAlive(input []byte, filename string) error
	GetLocation() string
}

type Manager struct {
	Quaratiner Quarantiner
	Scan       *Scan
}

func NewManager(quarantiner Quarantiner, scan *Scan) *Manager {
	//listen on channel
	return &Manager{quarantiner, scan}
}

func (q *Manager) HandleScanRequest(input io.Reader, uploadedFilename string, size int64) (ScanResponse, error) {
	logger := log.WithFields(log.Fields{"func": "HandleScanRequest"})

	contents, err := ioutil.ReadAll(input)
	if err != nil {
		logger.Error(err)
		return ScanResponse{}, err
	}

	//Quarantine File
	inertFilename, err := q.Quaratiner.RenderInert(contents, uploadedFilename)
	if err != nil {
		logger.Error(err)
		return ScanResponse{}, err
	}

	//Compute Hashes
	sha256 := q.computeSHA256(contents)
	md5 := q.computeMD5(contents)

	//Send rabbitmq
	go q.Scan.HandleMessage(rabbit.NewScanMessage(inertFilename,
		sha256,
		md5,
		q.Quaratiner.GetLocation(),
	))

	return ScanResponse{Filename: inertFilename,
		SHA256:    hex.EncodeToString(sha256[:]),
		MD5:       hex.EncodeToString(md5[:]),
		Permalink: "file/scan/" + base64.RawURLEncoding.EncodeToString(sha256[:]),
	}, nil

}

func (q *Manager) computeSHA256(input []byte) [sha256.Size]byte {
	return sha256.Sum256(input)
}

func (q *Manager) computeMD5(input []byte) [md5.Size]byte {
	return md5.Sum(input)
}
