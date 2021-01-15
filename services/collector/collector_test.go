package collector

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/config_client"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func InitHelper(t *testing.T) *Collector {
	t.Helper()
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := store.NewStore(cfg, log, db)
	conf := config_client.New(cfg, log)

	cli := clients.NewClients(cfg, log, conf)
	h := New(cfg, log, sto, cli)
	return h
}

func TestCollector_CollectResults(t *testing.T) {
	c := InitHelper(t)
	err := c.CollectResultsRound()
	assert.NoError(t, err)
}
