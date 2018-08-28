package requests

import (
	"encoding/json"
	"github.com/worlvlhole/maladapt/internal/file"
	"net/http"
)

type maladaptService struct {
	quarantine        *file.QuarantineAdmin
	uploadMaxFileSize int64
}

func (m *maladaptService) UploadFile(w http.ResponseWriter, r *http.Request) {
	files, present := r.MultipartForm.File["file"]
	if !present {
		WriteError(w, http.StatusBadRequest, InvalidKeySupplied)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			json.NewEncoder(w).Encode("Error opening file")
			return
		}

		// when do i close file?? TODO

		//Quarantine file
		//GenerateHash
		m.quarantine.Quaratiner.RenderInert(file, fileHeader.Filename, fileHeader.Size)

	}

	return
}

func NewMaladaptService(quarantineZone string) *maladaptService {
	return &maladaptService{
		quarantine: file.NewQuarantineAdmin(
			quarantineZone,
			file.NewZipQuarantiner(quarantineZone),
		),
	}
}
