package main

import (
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/worlvlhole/maladapt/internal/config"
	"github.com/worlvlhole/maladapt/internal/quarantine"
	"github.com/worlvlhole/maladapt/internal/repository"
	"github.com/worlvlhole/maladapt/internal/requests"
	"github.com/worlvlhole/maladapt/pkg/message/rabbit"
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

	//Producer Rabbit
	producer := rabbit.NewProducer(rabbit.NewRabbit(config.RabbitConfig))
	if err := producer.Connect(); err != nil {
		log.Fatal("Rabbitmq could not connect: ", err)
	}

	service := requests.NewMaladaptService(
		quarantine.NewManager(
			quarantine.NewZipQuarantiner(config.QuarantinePath),
			repository.NewScanMongoRepository(config.DBConfig),
			producer),
	)

	//Create Router
	r := chi.NewRouter()
	r.Use(requests.NewStructuredLogger(requestLogger))

	r.Route("/file", func(r chi.Router) {
		upload := r.Group(nil)
		upload.Use(requests.MaxBodySize(config.MaxUploadSize))
		upload.Use(requests.MultipartFormParse(config.MaxUploadSize))
		upload.Post("/scan", service.UploadFile) // POST /quarantine/scan
		upload.Get("/scan", service.GetResults)  // POST /quarantine/scan

		download := r.Group(nil)
		download.Get("/download/{hash}", service.DownloadFile) // POST /quarantine/scan
	})

	log.WithFields(log.Fields{"bind_address": config.BindAddress}).Info("Server Started")

	log.Fatal(http.ListenAndServe(config.BindAddress, r))
}
