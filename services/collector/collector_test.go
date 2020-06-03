package collector

import (
	"github.com/aibotsoft/micro/config"
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
	cli := clients.NewClients(nil, log, nil)
	h := New(cfg, log, sto, cli)
	return h
}

func TestCollector_CollectResults(t *testing.T) {
	c := InitHelper(t)
	err := c.CollectResults()
	assert.NoError(t, err)
}
