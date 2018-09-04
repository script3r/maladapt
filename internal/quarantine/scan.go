package quarantine

import (
	"crypto/sha256"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/repository"
	"github.com/worlvlhole/maladapt/pkg/message/rabbit"
)

type Scan struct {
	scanRepository repository.ScanRepository
	producer       *rabbit.Producer
}

func NewScan(scanRepository repository.ScanRepository, producer *rabbit.Producer) *Scan {
	return &Scan{
		scanRepository: scanRepository,
		producer:       producer,
	}
}

func (s *Scan) HandleMessage(msg *rabbit.ScanMessage) {
	logger := log.WithFields(log.Fields{"func": "HandleMessage"})
	logger.WithFields(log.Fields{
		"filename": msg.Filename,
		"sha256":   msg.SHA256,
		"md5":      msg.MD5}).Info()

	seenBefore, err := s.preFlightCheck(msg)
	if err != nil {
		logger.WithFields(log.Fields{"filename": msg.Filename}).Error(err)
		return
	}

	if seenBefore {
		logger.Info("File has been checked before")
		return
	}

	//Send to rabbitmq
	s.producer.Publish(msg)

}

//Check hashes from some db to make sure we don't perform unnecessary work.
func (s *Scan) preFlightCheck(msg *rabbit.ScanMessage) (bool, error) {
	logger := log.WithFields(log.Fields{"func": "PreFlightCheck"})

	d, err := base64.RawURLEncoding.DecodeString(msg.SHA256)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	var arr [sha256.Size]byte
	//TODO
	copy(arr[:], d)

	if res, err := s.scanRepository.FindBySHA256(arr); err != nil || res == nil {
		return false, nil
	}

	return true, nil
}
