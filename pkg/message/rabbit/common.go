package rabbit

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"github.com/streadway/amqp"
)

type ScanMessage struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
	MD5      string `json:"md5"`
	Path     string `json:"path"`
}

func NewScanMessage(filename string, hSha256 [sha256.Size]byte, hMd5 [md5.Size]byte, path string) *ScanMessage {
	return &ScanMessage{
		Filename: filename,
		SHA256:   base64.RawURLEncoding.EncodeToString(hSha256[:]),
		MD5:      base64.RawURLEncoding.EncodeToString(hMd5[:]),
		Path:     path,
	}
}

type Rabbit struct {
	config Configuration
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewRabbit(config Configuration) *Rabbit {
	return &Rabbit{
		config: config,
	}
}
