package quarantine

import (
	"crypto/md5"
	"crypto/sha256"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/model"
	"io"
	"io/ioutil"
)

type MaladaptFileMeta struct {
}

type Quarantiner interface {
	RenderInert(input []byte, filename string) error
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

func (q *Manager) HandleScanRequest(input io.Reader, filename string, size int64) error {
	logger := log.WithFields(log.Fields{"func": "HandleScanRequest"})

	contents, err := ioutil.ReadAll(input)
	if err != nil {
		logger.Error(err)
		return err
	}

	//Compute Hashes
	sha256 := q.computeSHA256(contents)
	md5 := q.computeMD5(contents)

	//TODO
	// Send to go channel
	q.Scan.Send(model.ScanMessage{Filename: filename,
		SHA256: sha256,
		MD5:    md5,
		Path:   q.Quaratiner.GetLocation(),
	})

	//QuarantineFile
	q.Quaratiner.RenderInert(contents, filename)
	return nil

}

func (q *Manager) computeSHA256(input []byte) [sha256.Size]byte {
	return sha256.Sum256(input)
}

func (q *Manager) computeMD5(input []byte) [md5.Size]byte {
	return md5.Sum(input)
}
