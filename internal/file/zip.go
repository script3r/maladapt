package file

import (
	"bytes"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type zipQuarantiner struct {
	QuarantineZone string
}

func (z *zipQuarantiner) RenderInert(reader io.Reader, fileName string, size int64) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		log.WithFields(log.Fields{"filename": fileName, "error": err}).Warn()
		return err
	}

	var outputb bytes.Buffer
	gz := gzip.NewWriter(&outputb)
	gz.Name = fileName

	w, err := gz.Write(b)
	if err != nil {
		log.WithFields(log.Fields{"filename": fileName, "bytes_written": w, "error": err}).Warn()
	}
	gz.Close()

	return nil

}

func (z *zipQuarantiner) RenderAlive(reader io.Reader, fileName string, size int64) error {
	return nil
}

func NewZipQuarantiner(quarantineZone string) *zipQuarantiner {
	return &zipQuarantiner{quarantineZone}
}
