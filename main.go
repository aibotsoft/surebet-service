package main

import (
	"fmt"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"github.com/aibotsoft/surebet-service/services/handler"
	"github.com/aibotsoft/surebet-service/services/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()
	log := logger.New()
	log.Infow("Begin service", "config", cfg)
	db := sqlserver.MustConnectX(cfg)
	sto := store.NewStore(cfg, log, db)
	cli := clients.NewClients(log)
	h := handler.NewHandler(cfg, log, sto, cli)
	s := server.NewServer(cfg, log, h)
	// Инициализируем Close
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errc <- s.Serve()
	}()
	defer func() {
		log.Debug("begin closing services")
		s.Close()
		_ = db.Close()
	}()

	log.Info("exit: ", <-errc)

}
