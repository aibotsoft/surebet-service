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

type SboClient struct {
	log  *zap.SugaredLogger
	conn *grpc.ClientConn
	pb.SurebetClient
}

func NewSboClient(cfg *config.Config, log *zap.SugaredLogger) *SboClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, net.JoinHostPort("", cfg.SboService.GrpcPort), grpc.WithInsecure(), grpc.WithBlock())
	logger.Panic(err, log, "NewSboClient DialContext error")
	cli := &SboClient{log: log, conn: conn, SurebetClient: pb.NewSurebetClient(conn)}
	log.Infow("begin ping to sbo server", "addr", cfg.SboService.GrpcPort)
	_, err = cli.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Panicw("sbo server do not response", "addr", cfg.SboService.GrpcPort)
	}
	return cli
}

func (p *SboClient) Close() {
	p.log.Info("closing sbobet connection")
	err := p.conn.Close()
	if err != nil {
		p.log.Error(err)
	}
}
