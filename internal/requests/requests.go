package requests

import (
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/quarantine"
	"net/http"
)

type maladaptService struct {
	manager *quarantine.Manager
}

func (m *maladaptService) UploadFile(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{"func": "UploadFile"})

	files, present := r.MultipartForm.File["file"]
	if !present {
		WriteError(w, http.StatusBadRequest, InvalidKeySupplied)
		return
	}

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Insufficient resources")
			logger.Error(err.Error())
			return
		}

		res, err := m.manager.HandleRequest(file, fileHeader.Filename, fileHeader.Size)
		if err != nil {

		}

		WriteSuccess(w, res)

		if err := file.Close(); err != nil {
			logger.Error(err.Error())
			return
		}
	}
}

func (m *maladaptService) DownloadFile(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{"func": "DownloadFile"})
	logger.Info(chi.URLParam(r, "hash"))

}

func (m *maladaptService) GetResults(w http.ResponseWriter, r *http.Request) {

}
func NewMaladaptService(manager *quarantine.Manager) *maladaptService {
	return &maladaptService{manager}
}
