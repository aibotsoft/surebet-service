package receiver

import (
	pb "github.com/aibotsoft/gen/fortedpb"
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
	nc, err := nats.Connect("nats://192.168.1.10:30873")
	if err != nil {
		log.Panic(err)
	}
	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Panic(err)
	}
	return &Receiver{cfg: cfg, log: log, handler: handler, nats: c}
}
func (r *Receiver) Close() {
	r.handler.Close()
}

func (r *Receiver) Job(sb *pb.Surebet) {
	//r.log.Info("receiver_got_sb: ", sb.FortedSurebetId)
	go r.handler.SurebetLoop(sb)
}

func (r *Receiver) Subscribe() {
	_, err := r.nats.Subscribe("surebet", r.Job)
	if err != nil {
		r.log.Error(err)
	}
}

func (r *Receiver) Send(sb *pb.Surebet) {
	r.nats.Publish("surebet", sb)
}
