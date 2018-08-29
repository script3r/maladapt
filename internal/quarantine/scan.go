package quarantine

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/model"
)

type Scan struct {
	uploadChan chan model.ScanMessage
}

func NewScan(uploadChan chan model.ScanMessage) *Scan {
	return &Scan{uploadChan: uploadChan}
}

func (s *Scan) Listen() {
	go func() {
		for {
			msg := <-s.uploadChan
			s.preFlightCheck(msg)
		}
	}()
}

func (s *Scan) preFlightCheck(msg model.ScanMessage) error {
	logger := log.WithFields(log.Fields{"func": "PreFlighCheck"})
	logger.Info()
	return nil
}

func (s *Scan) Send(msg model.ScanMessage) {
	logger := log.WithFields(log.Fields{"func": "Send"})
	logger.WithFields(log.Fields{
		"filename": msg.Filename,
		"sha256":   hex.EncodeToString(msg.SHA256[:]),
		"md5":      hex.EncodeToString(msg.MD5[:])}).Info()
	s.uploadChan <- msg
}
