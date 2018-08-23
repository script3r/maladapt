package requests

import (
	"net/http"
)

type maladaptService struct {
}

func (m * maladaptService) UploadFile(w http.ResponseWriter, r *http.Request) {
	 return
}

func NewMaladaptService() (*maladaptService) {
	return &maladaptService{}
}



