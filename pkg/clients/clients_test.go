package clients

import (
	"github.com/aibotsoft/micro/logger"
	"testing"
)

func TestNewClients(t *testing.T) {
	log := logger.New()
	got := NewClients(log)
	t.Log(got)
}
