package confighelper

import (
	"context"
	"github.com/leyle/go-api-starter/logmiddleware"
	"testing"
)

type Config struct {
	Debug   bool              `yaml:"debug"`
	Server  *ConnectionOption `yaml:"server"`
	Couchdb *ConnectionOption `yaml:"couchdb"`
}

func TestLoadConfig(t *testing.T) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)
	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	cfgPath := "../test/example_config.yaml"

	var cfg *Config
	err := LoadConfig(ctx, cfgPath, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(cfg.Server.Host)
	t.Log(cfg.Server.Port)
	t.Log(cfg.Server.User)
}
