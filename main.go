package main

import (
	"context"
	pb "github.com/aibotsoft/gen/surebetpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

func main() {
	cfg := config.New()
	log := logger.New()
	log.Infow("Begin service", "config", cfg)
	db := sqlserver.MustConnect(cfg)
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, net.JoinHostPort("", strconv.Itoa(cfg.ProxyService.GrpcPort)), grpc.WithInsecure(), grpc.WithBlock())
	logger.Panic(err, log, "grpc.DialContext error")
	pinClient := pb.NewSurebetClient(conn)
	line, err := pinClient.CheckLine(context.Background(), &pb.CheckLineRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Info(line)
}
