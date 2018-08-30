package quarantine

import (
	"bytes"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

const fileNameSuffix = "_inert"

type zipQuarantiner struct {
	quarantineZone string
}

func (z *zipQuarantiner) RenderInert(contents []byte, filename string) (inertFilename string, err error) {
	logger := log.WithFields(log.Fields{"func": "RenderInert"})

	//Create a File
	file, err := ioutil.TempFile(z.quarantineZone, filename)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	inertFilename = filepath.Base(file.Name())
	//Create gzip writer
	var outputb bytes.Buffer
	gz := gzip.NewWriter(&outputb)
	gz.Name = inertFilename

	//Write contents to gzip writer
	_, err = gz.Write(contents)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = gz.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	//write gzip writer contents to file
	file.Write(outputb.Bytes())
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return inertFilename, nil

}
func (z *zipQuarantiner) GetLocation() string {
	return z.quarantineZone
}

func (z *zipQuarantiner) RenderAlive(reader []byte, fileName string) error {
	return nil
}

func NewZipQuarantiner(quarantineZone string) *zipQuarantiner {
	return &zipQuarantiner{quarantineZone}
}
