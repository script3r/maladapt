package requests

import (
	"github.com/worlvlhole/maladapt/internal/file"
	"net/http"
)

type maladaptService struct {
	quarantine        *file.QuarantineAdmin
	uploadMaxFileSize int64
}

func (m *maladaptService) UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(m.uploadMaxFileSize)

	return
}

func NewMaladaptService(quarantineZone string, uploadMaxFileSize int64) *maladaptService {
	return &maladaptService{
		quarantine:        file.NewQuarantineAdmin(quarantineZone),
		uploadMaxFileSize: uploadMaxFileSize,
	}
}
