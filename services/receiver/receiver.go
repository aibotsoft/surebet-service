package receiver

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/surebet-service/services/handler"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Receiver struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	handler *handler.Handler
	nats    *nats.EncodedConn
}

func New(cfg *config.Config, log *zap.SugaredLogger, handler *handler.Handler) *Receiver {
	return &Receiver{cfg: cfg, log: log, handler: handler}
}
