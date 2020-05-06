package handler

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	firstBetMaxTry  = 1
	secondBetMaxTry = 1
)

var currencyList = []pb.Currency{{Code: "USD", Value: 1}, {Code: "EUR", Value: 0.93}}

type Handler struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	store   *store.Store
	clients clients.Clients
}

func NewHandler(cfg *config.Config, log *zap.SugaredLogger, store *store.Store, clients clients.Clients) *Handler {
	return &Handler{cfg: cfg, log: log, store: store, clients: clients}
}

func SurebetWithOneMember(sb *pb.Surebet, i int) *pb.Surebet {
	copySb := *sb
	copySb.Members = sb.Members[i : i+1]
	return &copySb
}

func (h *Handler) CheckLine(ctx context.Context, sb *pb.Surebet, i int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	name := sb.Members[i].ServiceName
	response, err := h.clients[name].CheckLine(ctx, &pb.CheckLineRequest{Surebet: SurebetWithOneMember(sb, i)})
	if err != nil {
		h.log.Errorw("check line error", "err", err, "name", name)
	}
	if response.GetSide() != nil {
		sb.Members[i] = response.GetSide()
		sb.Members[i].Check.Done = util.UnixMsNow()
	}
	h.log.Infow("check line resp", "name", name, "check", response.Side.Check)
}

func (h *Handler) PlaceBet(ctx context.Context, sb *pb.Surebet, i int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	side := sb.Members[i]
	if side.ToBet == nil {
		side.ToBet = &pb.ToBet{Id: util.UnixMsNow()}
	} else {
		side.ToBet.TryCount += side.ToBet.TryCount
	}
	side.Bet = &pb.Bet{Status: status.StatusError, Start: util.UnixMsNow()}
	response, err := h.clients[side.ServiceName].PlaceBet(ctx, &pb.PlaceBetRequest{Surebet: SurebetWithOneMember(sb, i)})
	if err != nil {
		h.log.Errorw("place bet error", "err", err, "name", side.ServiceName)
	}
	if response.GetSide() != nil {
		sb.Members[i] = response.GetSide()
		sb.Members[i].Bet.Done = util.UnixMsNow()
	}
	h.log.Infow("place bet resp", "name", side.ServiceName, "bet", sb.Members[i].Bet)
}

func (h *Handler) SendCheckLines(ctx context.Context, sb *pb.Surebet) {
	var wg sync.WaitGroup
	for i := range sb.Members {
		wg.Add(1)
		go h.CheckLine(ctx, sb, i, &wg)
	}
	wg.Wait()
}

func (h *Handler) PlaceSecond(ctx context.Context, sb *pb.Surebet) {
	h.PlaceBet(ctx, sb, int(sb.Calc.SecondIndex), nil)
}

func (h *Handler) SendPlaceFirst(ctx context.Context, sb *pb.Surebet) {
	var wg sync.WaitGroup
	for i := range sb.Members {
		wg.Add(1)
		if int64(i) == sb.Calc.FirstIndex {
			go h.PlaceBet(ctx, sb, i, &wg)
		} else {
			go h.CheckLine(ctx, sb, i, &wg)
		}
	}
	wg.Wait()
}
func (h *Handler) SaveSurebet(sb *pb.Surebet) {
	err := h.store.SaveFortedSurebet(sb)
	if err != nil {
		h.log.Error(err)
	}
	if sb.Calc == nil {
		return
	}
	h.log.Infow("save calc", "calc", sb.Calc)
	err = h.store.SaveCalc(sb)
	if err != nil {
		h.log.Error(err)
	}
	h.log.Infow("save sides", "0", sb.Members[0], "1", sb.Members[1])
	err = h.store.SaveSide(sb)
	if err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) SurebetLoop(sb *pb.Surebet) {
	err := h.ProcessSurebet(sb)
	if err != nil {
		h.log.Info(err)
	}
	h.SaveSurebet(sb)
}

func (h *Handler) ProcessSurebet(sb *pb.Surebet) error {
	ctx := context.Background()
	start := time.Now()

	if err := h.AllServicesActive(sb); err != nil {
		return err
	}
	err := h.store.LoadConfig(ctx, sb)
	if err != nil {
		return SurebetError{err, "store.LoadConfig error", true}
	}
	if err := h.AnyDisabled(sb); err != nil {
		return err
	}

	sb.Currency = currencyList
	sb.SurebetId = util.UnixUsNow()

	for i := 0; i < firstBetMaxTry; i++ {
		start := time.Now()
		h.SendCheckLines(ctx, sb)

		if !AllCheckStatusOk(sb) {
			h.log.Info("one of checks not ok")
			continue
		}

		h.Calc(sb)
		if !AllCheckCalcStatusOk(sb) {
			h.log.Infow("не все статусы ок", "profit", sb.Calc.Profit, "time", time.Since(start))
			continue
		}
		if !h.AllSurebet(sb) {
			return nil
		}

		h.log.Infow("begin bet first", "profit", sb.Calc.Profit, "stake", sb.Members[sb.Calc.FirstIndex].CheckCalc.Stake, "time", time.Since(start))
		h.SendPlaceFirst(ctx, sb)
		if sb.Members[sb.Calc.FirstIndex].Bet.Status == status.StatusOk {
			break
		}
	}
	if sb.Calc == nil {
		return SurebetError{nil, "no calc, returning...", false}
	}
	if sb.Members[sb.Calc.FirstIndex].Bet == nil {
		return SurebetError{nil, "no first bet", false}
	}
	if sb.Members[sb.Calc.FirstIndex].Bet.Status != status.StatusOk {
		h.log.Infow("first bet not ok", "name", sb.Members[sb.Calc.FirstIndex].ServiceName)
		return SurebetError{nil, "first bet not ok", false}
	}

	for i := 0; i < secondBetMaxTry; i++ {
		if sb.Members[sb.Calc.SecondIndex].Check.Status != status.StatusOk {
			h.CheckLine(ctx, sb, int(sb.Calc.SecondIndex), nil)
			continue
		}
		h.CalcSecond(sb)
		h.PlaceSecond(ctx, sb)
		if sb.Members[sb.Calc.SecondIndex].Bet.Status == status.StatusOk {
			break
		}
	}
	if sb.Members[sb.Calc.SecondIndex].Bet.Status != status.StatusOk {
		h.log.Infow("second bet not ok", "name", sb.Members[sb.Calc.SecondIndex].ServiceName)
		//	todo: облом, послать уведомление
	} else {
		h.log.Info("surebet done", "time", time.Since(start))
	}
	return nil
}
