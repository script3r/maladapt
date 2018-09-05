package quarantine

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/repository"
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

type Quarantiner interface {
	RenderInert(input []byte, filename string) (string, error)
	RenderAlive(input []byte, filename string) error
	GetLocation() string
}

type Manager struct {
	quarantiner    Quarantiner
	producer       *rabbit.Producer
	scanRepository repository.ScanRepository
}

func NewManager(quarantiner Quarantiner, scan repository.ScanRepository, producer *rabbit.Producer) *Manager {
	return &Manager{
		quarantiner:    quarantiner,
		producer:       producer,
		scanRepository: scan,
	}
}

func (q *Manager) HandleRequest(input io.Reader, uploadedFilename string, size int64) (ScanResponse, error) {
	logger := log.WithFields(log.Fields{"func": "HandleRequest"})

	contents, err := ioutil.ReadAll(input)
	if err != nil {
		logger.Error(err)
		return ScanResponse{}, err
	}

	//Quarantine File
	inertFilename, err := q.quarantiner.RenderInert(contents, uploadedFilename)
	if err != nil {
		logger.Error(err)
		return ScanResponse{}, err
	}

	//Compute Hashes
	sha256 := q.computeSHA256(contents)
	md5 := q.computeMD5(contents)

	//Send rabbitmq
	go q.SendToScanners(rabbit.NewScanMessage(inertFilename,
		sha256,
		md5,
		q.quarantiner.GetLocation(),
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

func (q *Manager) SendToScanners(msg *rabbit.ScanMessage) {
	logger := log.WithFields(log.Fields{"func": "SendToScanners"})
	logger.WithFields(log.Fields{
		"filename": msg.Filename,
		"sha256":   msg.SHA256,
		"md5":      msg.MD5}).Info("Sending message to rabbit")

	q.producer.Publish(msg)

}
