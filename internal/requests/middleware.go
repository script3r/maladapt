package requests

import (
	"net/http"
)

func MaxBodySize(maxMemory int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Save us from too big of files..
			r.Body = http.MaxBytesReader(w, r.Body, maxMemory)
			next.ServeHTTP(w, r)
			return
		}
		return http.HandlerFunc(fn)
	}
}

func MultipartFormParse(maxMemory int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// This max size is only for RAM. Anything read over is put on disk.
			err := r.ParseMultipartForm(maxMemory)
			if err != nil {
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

}
