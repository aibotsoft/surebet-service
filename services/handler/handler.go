package handler

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/config_client"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/telegram"
	"github.com/aibotsoft/micro/util"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	repeatMinTimeout    = 3 * time.Second
	secondBetMaxTry     = 57
	surebetLoopMaxCount = 20
)

var loopIds sync.Map

type Handler struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	store   *store.Store
	clients clients.Clients
	tel     *telegram.Telegram
	Conf    *config_client.ConfClient
}

func NewHandler(cfg *config.Config, log *zap.SugaredLogger, store *store.Store, clients clients.Clients, conf *config_client.ConfClient) *Handler {
	tel := telegram.New(cfg, log)
	return &Handler{cfg: cfg, log: log, store: store, clients: clients, tel: tel, Conf: conf}
}

func (h *Handler) Close() {
	h.store.Close()
	h.Conf.Close()
}
func (h *Handler) GetCurrency(ctx context.Context) ([]pb.Currency, error) {
	get, b := h.store.Cache.Get("currency_list")
	if b {
		return get.([]pb.Currency), nil
	}
	currency, err := h.Conf.GetCurrency(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get currency error")
	}
	var cur []pb.Currency
	err = copier.Copy(&cur, &currency)
	if err != nil {
		return nil, errors.Wrap(err, "copy error")
	}
	h.store.Cache.SetWithTTL("currency_list", cur, 1, time.Hour)
	return cur, nil
}

func (h *Handler) CheckLine(ctx context.Context, sb *pb.Surebet, i int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	side := sb.Members[i]
	side.Check = &pb.Check{Status: status.StatusError, Id: util.UnixMsNow()}
	response, err := h.clients[side.ServiceName].CheckLine(ctx, &pb.CheckLineRequest{Surebet: SurebetWithOneMember(sb, i)})
	if err != nil {
		h.log.Errorw("check line error", "err", err, "name", side.ServiceName)
		side.Check.Status = status.ServiceError
		side.Check.StatusInfo = "service error"
		side.Check.Done = util.UnixMsNow()
		return
	}
	if response.GetSide() != nil {
		sb.Members[i] = response.GetSide()
		sb.Members[i].Check.Done = util.UnixMsNow()
	}
	if sb.Members[i].GetCheck().GetStatus() != status.StatusOk {
		h.log.Infow("check_line_resp", "name", side.ServiceName, "check", response.Side.Check, "marketName", side.MarketName, "sport", side.SportName, "league", side.LeagueName,
			"home", side.Home, "away", side.Away)
	}
}

func (h *Handler) PlaceBet(ctx context.Context, sb *pb.Surebet, i int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	side := sb.Members[i]
	if side.ToBet == nil {
		side.ToBet = &pb.ToBet{Id: util.UnixMsNow()}
	} else {
		h.log.Info("add try count")
		side.ToBet.TryCount = side.ToBet.TryCount + 1
	}
	side.Bet = &pb.Bet{Status: status.StatusError, Start: util.UnixMsNow()}
	response, err := h.clients[side.ServiceName].PlaceBet(ctx, &pb.PlaceBetRequest{Surebet: SurebetWithOneMember(sb, i)})
	if err != nil {
		h.log.Errorw("place bet error", "err", err, "name", side.ServiceName)
		side.Bet.StatusInfo = "service error"
		side.Bet.Done = util.UnixMsNow()
		return
	}
	if response.GetSide() != nil {
		sb.Members[i] = response.GetSide()
		sb.Members[i].Bet.Done = util.UnixMsNow()
	}
	if sb.Members[i].GetBet().GetStatus() != status.StatusOk {
		h.log.Infow("place bet resp", "name", side.ServiceName, "bet", sb.Members[i].GetBet())
	}
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
	if sb.Calc == nil {
		return
	}
	if err := AllCheckCalcStatusOk(sb); err != nil {
		return
	}
	if err := h.store.SaveFortedSurebet(sb); err != nil {
		h.log.Error(err)
	}
	//h.log.Infow("save surebet", "profit", sb.Calc.Profit, "time", ElapsedFromSurebetId(sb.SurebetId))
	if err := h.store.SaveCalc(sb); err != nil {
		h.log.Error(err)
	}
	//h.log.Infow("save sides", "0", sb.Members[0], "1", sb.Members[1])
	if err := h.store.SaveSide(sb); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) SurebetLoop(sb *pb.Surebet) {
	_, ok := h.store.Cache.Get(sb.FortedSurebetId)
	//_, ok := loopIds.Load(sb.FortedSurebetId)
	if ok {
		h.log.Infow("loop_already_exists", "id", sb.FortedSurebetId)
		return
	}
	h.store.Cache.Set(sb.FortedSurebetId, true, 1)
	defer h.store.Cache.Del(sb.FortedSurebetId)
	//loopIds.Store(sb.FortedSurebetId, nil)
	//defer loopIds.Delete(sb.FortedSurebetId)
	for i := 0; i < surebetLoopMaxCount; i++ {
		ClearSurebet(sb)
		err := h.ProcessSurebet(sb)
		h.SaveSurebet(sb)
		if err != nil {
			var otherName string
			if sb.Members[0].ServiceName == err.ServiceName {
				otherName = sb.Members[1].ServiceName
			} else {
				otherName = sb.Members[0].ServiceName
			}
			h.log.Infow("result", "err", err, "name", err.ServiceName, "time", ElapsedFromSurebetId(sb.SurebetId), "loop", i, "fid", sb.FortedSurebetId, "other", otherName)
			if err.Permanent {
				//h.log.Info("error permanent, so returning...")
				return
			}
		} else {
			go h.tel.Sendf("surebet v=%v, l=%v, p=%v, f=%v, s=%v, t=%v", betAmount(sb), i, sb.Calc.Profit, sb.Members[0].ServiceName, sb.Members[1].ServiceName, ElapsedFromSurebetId(sb.SurebetId))
			h.log.Infow("placed_surebet", "profit", sb.GetCalc().GetProfit(), "time", ElapsedFromSurebetId(sb.SurebetId), "loop", i, "fid", sb.FortedSurebetId)
			i = 0
		}
		if sb.GetCalc().GetProfit() < -9 {
			h.log.Infow("profit_too_low, so returning...", "profit", sb.GetCalc().GetProfit(), "fid", sb.FortedSurebetId)
		}
		time.Sleep(repeatMinTimeout + time.Millisecond*100*time.Duration(i))
	}
}

var sbobetLock = &StatefulLock{}
var dafLock = &StatefulLock{}

func (h *Handler) ProcessSurebet(sb *pb.Surebet) *SurebetError {
	ctx := context.Background()
	sb.SurebetId = util.UnixUsNow()
	var err error
	if err := h.AllServicesActive(sb); err != nil {
		return err
	}
	if err := h.store.LoadConfig(ctx, sb); err != nil {
		return &SurebetError{Err: err, Msg: "store.LoadConfig error", Permanent: true}
	}
	if err := h.AnyDisabled(sb); err != nil {
		return err
	}

	sb.Currency, err = h.GetCurrency(ctx)
	if err != nil {
		return &SurebetError{Err: err, Msg: "get currency error", Permanent: true}
	}
	if AnyHasName(sb, "Sbobet") {
		sbobetLock.Take(sb.SurebetId)
		defer sbobetLock.Release(sb.SurebetId)
	}
	if AnyHasName(sb, "Dafabet") {
		dafLock.Take(sb.SurebetId)
		defer dafLock.Release(sb.SurebetId)
	}
	//h.log.Infow("marketId", "0", sb.Members[0].MarketId, "1", sb.Members[1].MarketId)

	h.SendCheckLines(ctx, sb)

	if err := AllCheckStatus(sb); err != nil {
		return err
	}
	if err := h.Calc(sb); err != nil {
		return err
	}
	if err := AllCheckCalcStatusOk(sb); err != nil {
		return err
	}
	if err := h.AllSurebet(sb); err != nil {
		return err
	}
	h.log.Infow("begin_bet_first", "profit", sb.Calc.Profit, "stake", sb.Members[sb.Calc.FirstIndex].CheckCalc.Stake, "time", ElapsedFromSurebetId(sb.SurebetId),
		"fid", sb.FortedSurebetId, "price", sb.Members[sb.Calc.FirstIndex].Check.Price)
	h.SendPlaceFirst(ctx, sb)

	if sb.Members[sb.Calc.FirstIndex].Bet.Status != status.StatusOk {
		if sb.Members[sb.Calc.FirstIndex].Bet.Status == status.AboveEventMax {
			return &SurebetError{Msg: "first_bet_not_ok", Permanent: true, ServiceName: sb.Members[sb.Calc.FirstIndex].ServiceName}
		}
		return &SurebetError{Msg: "first_bet_not_ok", Permanent: false, ServiceName: sb.Members[sb.Calc.FirstIndex].ServiceName}

	}
	switch sb.Members[sb.Calc.FirstIndex].ServiceName {
	case "Sbobet":
		h.log.Info("release sbobet lock after first bet")
		sbobetLock.Release(sb.SurebetId)
	case "Dafabet":
		h.log.Info("release dafabet lock after first bet")
		dafLock.Release(sb.SurebetId)
	}

	for i := 0; i < secondBetMaxTry; i++ {
		if sb.Members[sb.Calc.SecondIndex].Check.GetStatus() != status.StatusOk {
			time.Sleep(time.Millisecond * 100 * time.Duration(i))
			h.CheckLine(ctx, sb, int(sb.Calc.SecondIndex), nil)
			continue
		}
		h.CalcSecond(sb)
		h.PlaceSecond(ctx, sb)
		if sb.Members[sb.Calc.SecondIndex].GetBet().GetStatus() == status.StatusOk {
			break
		} else if i < secondBetMaxTry {
			sb.Members[sb.Calc.SecondIndex].Check.Status = status.StatusError
		}
	}
	if sb.Members[sb.Calc.SecondIndex].GetBet().GetStatus() != status.StatusOk {
		h.log.Infow("second bet not ok", "name", sb.Members[sb.Calc.SecondIndex].ServiceName)
		go h.tel.Send("oblom")
		return &SurebetError{Msg: "second bet not ok", Permanent: true, ServiceName: sb.Members[sb.Calc.SecondIndex].ServiceName}
	}
	return nil
}
