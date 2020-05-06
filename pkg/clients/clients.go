package clients

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

//type Clients struct {
//	log  *zap.SugaredLogger
//	conn *grpc.ClientConn
//	pb.SurebetClient
//}

type Clients map[string]pb.FortedClient

var clientList = map[string]string{"Pinnacle": "50054", "Sbobet": "50053"}

//var clientList = map[string]string{ "Sbobet": "50053"}

func NewClients(log *zap.SugaredLogger) Clients {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clients := make(Clients)
	for name, addr := range clientList {
		conn, err := grpc.DialContext(ctx, net.JoinHostPort("", addr), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Panicw("dial to client error", "name", name, "addr", addr)
		}
		clients[name] = pb.NewFortedClient(conn)
		log.Infow("begin ping to server", "name", name, "addr", addr)
		_, err = clients[name].Ping(ctx, &pb.PingRequest{})
		if err != nil {
			log.Panicw("server do not response to ping", "name", name, "addr", addr)
		}
		log.Infow("server responded", "name", name, "addr", addr)
	}
	return clients
}

//func (p *Clients) Close() {
//	err := p.conn.Close()
//	if err != nil {
//		p.log.Error(err)
//	}
//}
