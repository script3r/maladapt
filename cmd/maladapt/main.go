package main

import (
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/worlvlhole/maladapt/internal/requests"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
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

	//Environment Variables
	uploadMaxFileSize, err := strconv.ParseInt(os.Getenv("MALADAPT_MAX_UPLOAD_SIZE"), 10, 64)
	if err != nil {
		log.Fatalf("could not parse `ADMIN_LDAP_PORT`: %v", err)
	}

	addr := os.Getenv("MALADAPT_BIND_ADDR")
	if addr == "" {
		log.Fatal("could not fined `MALADAPT_BIND_ADDR`: %v", err)
	}
	quarantined_zone := os.Getenv("MALADAPT_QUARANTINE_PATH")

	r := chi.NewRouter()
	r.Use(requests.NewStructuredLogger(requestLogger))

	//Create MaladaptService
	service := requests.NewMaladaptService(quarantined_zone, uploadMaxFileSize)

	r.Route("/file", func(r chi.Router) {
		r.Post("/scan", service.UploadFile) // POST /file/scan
	})

	log.WithFields(log.Fields{"bindaddress": addr}).Info("Server Started")
	log.Fatal(http.ListenAndServe(addr, r))
}
