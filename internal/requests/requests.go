package requests

import (
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/file"
	"net/http"
)

type maladaptService struct {
	quarantine *file.QuarantineAdmin
}

func (m *maladaptService) UploadFile(w http.ResponseWriter, r *http.Request) {
	log := log.WithFields(log.Fields{"func": "UploadFile"})

	files, present := r.MultipartForm.File["file"]
	if !present {
		WriteError(w, http.StatusBadRequest, InvalidKeySupplied)
		return
	}

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Insufficient resources")
			log.Error(err.Error())
			return
		}

		m.quarantine.Handle(file, fileHeader.Filename, fileHeader.Size)

		if err := file.Close(); err != nil {
			log.Error(err.Error())
			return
		}
	}
}

func (m *maladaptService) DownloadFile(w http.ResponseWriter, r *http.Request) {

}

func NewMaladaptService(quarantineZone string) *maladaptService {
	return &maladaptService{
		quarantine: file.NewQuarantineAdmin(
			quarantineZone,
			file.NewZipQuarantiner(quarantineZone),
		),
	}
}
