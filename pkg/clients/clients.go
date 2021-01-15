package clients

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/config_client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"strconv"
)

type Clients map[string]pb.FortedClient

var CloneMap = map[string]string{
	"Betdaq":       "SportMarket",
	"Matchbook":    "SportMarket",
	"Betfair":      "SportMarket",
	"SbobetMarket": "SportMarket",
	"Black":        "SportMarket",
}
var SportMarketClones = []string{"Betdaq", "Matchbook", "Betfair", "SbobetMarket", "SportMarket", "Black"}

func NewClients(cfg *config.Config, log *zap.SugaredLogger, conf *config_client.ConfClient) Clients {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Service.GrpcTimeout)
	services, err := conf.GetServices(ctx)
	cancel()
	if err != nil {
		log.Panic(err)
	}
	log.Infow("service_list", "", services)
	clients := make(Clients)
	for _, s := range services {
		addr := ""
		if os.Getenv("is_debug") != "" {
			addr = net.JoinHostPort("", strconv.FormatInt(s.GrpcPort, 10))
		} else {
			addr = net.JoinHostPort(s.ServiceName, strconv.FormatInt(50051, 10))
		}
		log.Info("begin_dial_service: ", addr, " err: ", err)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Service.GrpcTimeout)
		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		cancel()
		if err != nil {
			log.Infow("dial_to_client_error", "name", s.ServiceName, "addr", addr, " err", err)
			continue
		}
		c := pb.NewFortedClient(conn)
		log.Infow("begin ping to server", "name", s.ServiceName, "addr", addr)
		ctx, cancel = context.WithTimeout(context.Background(), cfg.Service.GrpcTimeout)
		_, err = c.Ping(ctx, &pb.PingRequest{})
		cancel()
		if err != nil {
			log.Errorw("server do not response to ping", "name", s.ServiceName, "addr", addr)
			continue
		}
		log.Infow("server responded", "name", s.ServiceName, "addr", addr)
		clients[s.FortedName] = c
	}
	for cloneName, serviceName := range CloneMap {
		clients[cloneName] = clients[serviceName]
	}
	return clients
}
