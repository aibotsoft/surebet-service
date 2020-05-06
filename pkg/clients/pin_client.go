package clients

import (
	"context"
	pb "github.com/aibotsoft/gen/surebetpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

type PinClient struct {
	log  *zap.SugaredLogger
	conn *grpc.ClientConn
	pb.SurebetClient
}

func NewPinClient(cfg *config.Config, log *zap.SugaredLogger) *PinClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, net.JoinHostPort("", cfg.PinService.GrpcPort), grpc.WithInsecure(), grpc.WithBlock())
	logger.Panic(err, log, "NewPinClient DialContext error")
	cli := &PinClient{log: log, conn: conn, SurebetClient: pb.NewSurebetClient(conn)}
	log.Infow("begin ping to pin server", "addr", cfg.PinService.GrpcPort)
	_, err = cli.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Panicw("pin server do not response", "addr", cfg.PinService.GrpcPort)
	}
	log.Infow("pin server responded", "addr", cfg.PinService.GrpcPort)
	return cli
}
func (p *PinClient) Close() {
	err := p.conn.Close()
	if err != nil {
		p.log.Error(err)
	}
}
