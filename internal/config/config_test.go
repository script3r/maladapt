package config

import (
	"testing"
)

func TestConfig(t *testing.T) {

	config := Initialize()
	if config.BindAddress == "" {
		t.Error("empty bind_addr")
	}

	if config.QuarantinePath == "" {
		t.Error("quarantinePath")
	}

	t.Logf("%+v", config)

}
