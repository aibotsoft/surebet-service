package store

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/surebet-service/pkg/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func InitHelper(t *testing.T) *Store {
	t.Helper()
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := NewStore(cfg, log, db)
	return sto
}
func TestStore_LoadConfig(t *testing.T) {
	s := InitHelper(t)
	sur := tests.SurebetHelper(t)
	err := s.LoadConfig(context.Background(), sur)
	assert.NoError(t, err)
}

func TestStore_SaveFortedSurebet(t *testing.T) {
	s := InitHelper(t)
	sur := tests.SurebetHelper(t)
	err := s.SaveFortedSurebet(sur)
	assert.NoError(t, err)
}

func TestStore_SaveCalc(t *testing.T) {
	s := InitHelper(t)
	sur := tests.SurebetHelper(t)
	err := s.SaveCalc(sur)
	assert.NoError(t, err)
}

func TestStore_SaveSide(t *testing.T) {
	s := InitHelper(t)
	sur := tests.SurebetHelper(t)
	err := s.SaveSide(sur)
	assert.NoError(t, err)
}
