package main

import (
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/worlvlhole/maladapt/internal/config"
	"github.com/worlvlhole/maladapt/internal/requests"
	"log/syslog"
	"net/http"
)

func main() {

	//Setup Logger
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

	config := config.Initialize()
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	//Create MaladaptService
	service := requests.NewMaladaptService(config.QuarantinePath)

	//Create Router
	r := chi.NewRouter()
	r.Use(requests.NewStructuredLogger(requestLogger))

	r.Route("/file", func(r chi.Router) {
		upload := r.Group(nil)
		upload.Use(requests.MaxBodySize(config.MaxUploadSize))
		upload.Use(requests.MultipartFormParse(config.MaxUploadSize))
		upload.Post("/scan", service.UploadFile) // POST /file/scan

		download := r.Group(nil)
		download.Get("/download/{hash}", service.DownloadFile) // POST /file/scan
	})

	log.WithFields(log.Fields{"bindaddress": config.BindAddress}).Info("Server Started")

	log.Fatal(http.ListenAndServe(config.BindAddress, r))
}
