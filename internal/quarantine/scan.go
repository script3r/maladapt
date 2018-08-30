package quarantine

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/model"
	"github.com/worlvlhole/maladapt/internal/repository"
)

type Scan struct {
	uploadChan     chan model.ScanMessage
	scanRepository repository.ScanRepository
}

func NewScan(uploadChan chan model.ScanMessage, scanRepository repository.ScanRepository) *Scan {
	return &Scan{uploadChan: uploadChan, scanRepository: scanRepository}
}

func (s *Scan) Listen() {
	go func() {
		for {
			s.handleMessage(<-s.uploadChan)
		}
	}()
}

func (s *Scan) handleMessage(msg model.ScanMessage) {
	logger := log.WithFields(log.Fields{"func": "handleMessage"})

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

}

//Check hashes from some db to make sure we don't perform unnecessary work.
func (s *Scan) preFlightCheck(msg model.ScanMessage) (bool, error) {
	logger := log.WithFields(log.Fields{"func": "PreFlighCheck"})

	if res, err := s.scanRepository.FindBySHA256(msg.SHA256); err != nil || res == nil {
		return false, nil
	}

	logger.Info()
	return true, nil
}

func (s *Scan) Send(msg model.ScanMessage) {
	logger := log.WithFields(log.Fields{"func": "Send"})
	logger.WithFields(log.Fields{
		"filename": msg.Filename,
		"sha256":   hex.EncodeToString(msg.SHA256[:]),
		"md5":      hex.EncodeToString(msg.MD5[:])}).Info()
	s.uploadChan <- msg
}
