package file

import (
	"bytes"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const fileNameSuffix = "_inert"

type zipQuarantiner struct {
	QuarantineZone string
}

func (z *zipQuarantiner) RenderInert(contents []byte, filename string) error {
	logger := log.WithFields(log.Fields{"func": "RenderInert"})

	//Create a File
	inertFilename := filename + fileNameSuffix
	file, err := os.Create(filepath.Join(z.QuarantineZone, inertFilename))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	//Create gzip writer
	var outputb bytes.Buffer
	gz := gzip.NewWriter(&outputb)
	gz.Name = inertFilename

	//Write contents to gzip writer
	w, err := gz.Write(contents)
	if err != nil {
		logger.WithFields(log.Fields{"filename": inertFilename, "bytes_written": w}).Error(err.Error())
		return err
	}

	err = gz.Close()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	//write gzip writer contents to file
	file.Write(outputb.Bytes())
	file.Close()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil

}

func (z *zipQuarantiner) RenderAlive(reader []byte, fileName string) error {
	return nil
}

func NewZipQuarantiner(quarantineZone string) *zipQuarantiner {
	return &zipQuarantiner{quarantineZone}
}
