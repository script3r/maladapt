package quarantine

import (
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
		"md5":      msg.MD5}).Info("Sending message to rabbit")

	//Send to rabbitmq
	s.producer.Publish(msg)

}
