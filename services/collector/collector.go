package collector

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"go.uber.org/zap"
	"time"
)

type Collector struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	store   *store.Store
	clients clients.Clients
}

const CollectJobPeriod = 4 * time.Minute

func (c *Collector) CollectJob() {
	time.Sleep(time.Second)
	for {
		err := c.CollectResultsRound()
		if err != nil {
			c.log.Error(err)
		}
		time.Sleep(CollectJobPeriod)
	}
}
func (c *Collector) CollectResultsRound() error {
	var res []pb.BetResult
	for name, client := range c.clients {
		_, ok := clients.CloneMap[name]
		if ok {
			continue
		}
		results, err := client.GetResults(context.Background(), &pb.GetResultsRequest{})
		if err != nil {
			c.log.Error(err)
			continue
		}
		//c.log.Infow("res_count", "name", name, "count", len(results.GetResults()))
		//if name == "Dafabet" {
		//	c.log.Infow("", "", results.GetResults())
		//}
		res = append(res, results.GetResults()...)
	}
	err := c.store.SaveBetList(res)
	return err
}

func New(cfg *config.Config, log *zap.SugaredLogger, store *store.Store, clients clients.Clients) *Collector {
	return &Collector{cfg: cfg, log: log, store: store, clients: clients}
}

//func (c *Collector) CollectBalance() error {
//	var res []pb.BetResult
//	for _, client := range c.clients {
//		results, err := client.GetResults(context.Background(), &pb.GetResultsRequest{})
//		if err != nil {
//			c.log.Error(err)
//			continue
//		}
//		res = append(res, results.GetResults()...)
//	}
//	err := c.store.SaveBetList(res)
//	return err
//}
