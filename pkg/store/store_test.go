package store

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var s *Store

func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	s = NewStore(cfg, log, db)
	m.Run()
	s.Close()
}

func TestStore_GetConfigByName(t *testing.T) {
	got, err := s.GetConfigByName(context.Background(), "Pinnacle")
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got)
	}
	time.Sleep(time.Millisecond)
	got, err = s.GetConfigByName(context.Background(), "Pinnacle")

}
