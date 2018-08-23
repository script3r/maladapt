package main

import (
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
	"github.com/go-chi/chi"
	"net/http"
	"os"
	"github.com/worlvlhole/maladapt/internal/requests"
)

func main() {
	addr := os.Getenv("MALADAPT_BIND_ADDR")

	log.SetFormatter(&log.JSONFormatter{})

	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "maladapt")
	if err != nil {
		log.Fatal("could not setup syslog logger")
	}
	log.AddHook(hook)

	requestLogger := log.New()
	requestLogger.Formatter = &log.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}

	r := chi.NewRouter()
	r.Use(requests.NewStructuredLogger(requestLogger))

	service := requests.NewMaladaptService()

	r.Route("/file", func(r chi.Router) {
		r.Post("/scan", service.UploadFile) // POST /file/scan
	})

	http.ListenAndServe(addr, r)
}
