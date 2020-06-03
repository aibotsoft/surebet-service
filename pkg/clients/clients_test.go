package clients

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/config_client"
	"github.com/aibotsoft/micro/logger"
	"testing"
)

func TestNewClients(t *testing.T) {
	cfg := config.New()
	log := logger.New()
	conf := config_client.New(cfg, log)
	c := NewClients(cfg, log, conf)
	log.Info(c)
}
