package handler

import (
	"context"
	"fmt"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/config_client"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/telegram"
	"github.com/aibotsoft/micro/util"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/loop_count"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"github.com/aibotsoft/surebet-service/pkg/time_id"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"strconv"
	"sync"
	"time"
)

const (
	repeatMinTimeout    = 2 * time.Second
	checkTimeOut        = 5 * time.Second
	secondBetMaxTry     = 95
	surebetLoopMaxCount = 2
)

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

func (h *Handler) CheckLine(ctx context.Context, sb *pb.Surebet, i int64, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	side := sb.Members[i]
	side.Check = &pb.Check{Status: status.StatusError, Id: util.UnixMsNow()}
	response, err := h.clients[side.ServiceName].CheckLine(ctx, &pb.CheckLineRequest{Surebet: SurebetWithOneMember(sb, i)})
	if err != nil {
		side.Check.Done = util.UnixMsNow()
		st := gstatus.Convert(err)
		if st.Code() == codes.DeadlineExceeded {
			side.Check.Status = status.DeadlineExceeded
		} else {
			h.log.Errorw("check_line_error", "err", err, "name", side.ServiceName)
			side.Check.Status = status.ServiceError
			side.Check.StatusInfo = fmt.Sprintf("code:%q, msg:%q", st.Code(), st.Message())
		}
	} else if response.GetSide() != nil {
		sb.Members[i] = response.GetSide()
		sb.Members[i].Check.Done = util.UnixMsNow()
		switch sb.Members[i].GetCheck().GetStatus() {
		case status.StatusOk:
		case status.Suspended:
		case status.ServiceBusy:
		case status.MarketClosed:
		case status.HandicapChanged:
		case status.StatusNotFound:
		case status.PitchersRequired:
		default:
			//h.log.Infow("check_line_resp", "name", side.ServiceName, "check", response.Side.Check, "marketName", side.MarketName, "sport", side.SportName, "league", side.LeagueName,
			//	"home", side.Home, "away", side.Away)
			h.log.Infow("check_not_ok", "status", response.Side.Check.Status, "info", response.Side.Check.StatusInfo, "name", side.ServiceName, "m", side.MarketName, "s", side.SportName, "l", side.LeagueName,
				"h", side.Home, "a", side.Away, "fid", sb.FortedSurebetId)
		}
	}
}

func (h *Handler) PlaceBet(ctx context.Context, sb *pb.Surebet, i int64, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	side := sb.Members[i]
	if side.ToBet == nil {
		side.ToBet = &pb.ToBet{Id: util.UnixMsNow()}
	} else {
		h.log.Infow("add_try_count", "fid", sb.FortedSurebetId)
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
		h.log.Infow("place_bet_not_ok", "name", side.ServiceName, "bet", sb.Members[i].GetBet(), "fid", sb.FortedSurebetId)
	}
}

func (h *Handler) SendCheckLines(ctx context.Context, sb *pb.Surebet) {
	var wg sync.WaitGroup
	wg.Add(2)
	go h.CheckLine(ctx, sb, 0, &wg)
	go h.CheckLine(ctx, sb, 1, &wg)
	wg.Wait()
}

func (h *Handler) PlaceSecond(ctx context.Context, sb *pb.Surebet) {
	h.PlaceBet(ctx, sb, sb.Calc.SecondIndex, nil)
}

func (h *Handler) SendPlaceFirst(ctx context.Context, sb *pb.Surebet) {
	var wg sync.WaitGroup
	wg.Add(2)
	go h.PlaceBet(ctx, sb, sb.Calc.FirstIndex, &wg)
	go h.CheckLine(ctx, sb, sb.Calc.SecondIndex, &wg)
	wg.Wait()
}
func (h *Handler) SaveSurebet(sb *pb.Surebet) {
	if !h.HasAnyBet(sb) {
		return
	}

	if err := h.store.SaveFortedSurebet(sb); err != nil {
		h.log.Error(err)
	}
	if err := h.store.SaveCalc(sb); err != nil {
		h.log.Error(err)
	}
	if err := h.store.SaveSide(sb); err != nil {
		h.log.Error(err)
	}
}

var lc = loop_count.NewLoopCount()

func (h *Handler) SurebetLoop(sb *pb.Surebet) {
	lc.Add()
	defer lc.Remove()
	_, ok := h.store.Cache.Get(sb.FortedSurebetId)
	if ok {
		return
	}
	_, ok = h.store.Cache.Get("permanent:" + strconv.FormatInt(sb.FortedSurebetId, 10))
	if ok {
		return
	}
	h.store.Cache.Set(sb.FortedSurebetId, true, 1)
	defer h.store.Cache.Del(sb.FortedSurebetId)
	for i := 0; i < surebetLoopMaxCount; i++ {
		start := util.UnixMsNow()
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
			h.log.Infow("", "e", err, "n", err.ServiceName, "fid", sb.FortedSurebetId, "o", otherName, "t", ElapsedFromId(start),
				"g", fmt.Sprintf("%v-%v-%v:%v", sb.FortedSport, sb.FortedHome, sb.FortedAway, sb.Members[0].MarketName))
			if err.Permanent {
				h.store.Cache.SetWithTTL("permanent:"+strconv.FormatInt(sb.FortedSurebetId, 10), true, 1, time.Minute*7)
				return
			}
			time.Sleep(repeatMinTimeout + time.Millisecond*100*time.Duration(i))
			if lc.Get() > 1 {
				//h.log.Info("loop_count: ", lc.Get())
				return
			}
		} else {
			rp := ActualProfit(sb)
			go h.tel.Sendf("surebet v=%v, p=%v, f=%v, s=%v, t=%v", betAmount(sb), rp, sb.Members[0].ServiceName, sb.Members[1].ServiceName, ElapsedFromId(start))
			h.log.Infow("placed_surebet", "p", rp, "t", ElapsedFromId(start), "amount", betAmount(sb), "fid", sb.FortedSurebetId,
				"fs", sb.Members[0].ServiceName, "ss", sb.Members[1].ServiceName, "s", sb.FortedSport, "h", sb.FortedHome, "a", sb.FortedAway)
			i = 0
			if sb.Members[0].BetConfig.MaxCountLine == 1 || sb.Members[1].BetConfig.MaxCountLine == 1 {
				h.log.Infow("max_count_line_is_one_so_exit", "fid", sb.FortedSurebetId)
				h.store.Cache.SetWithTTL("permanent:"+strconv.FormatInt(sb.FortedSurebetId, 10), true, 1, time.Minute*5)
				time.Sleep(time.Millisecond * 100)
				return
			}
			time.Sleep(time.Millisecond * 200)
		}
	}
	//LoopCount[sb.FortedSurebetId] += 1
	////h.log.Debugw("loop_count", "fid", sb.FortedSurebetId, "count", LoopCount[sb.FortedSurebetId])
	//if LoopCount[sb.FortedSurebetId] > 3 && sb.Calc.Profit < 0 {
	//	//h.log.Debug("loop_count>2")
	//	h.store.Cache.SetWithTTL("permanent:"+strconv.FormatInt(sb.FortedSurebetId, 10), true, 1, time.Minute)
	//}
}

func (h *Handler) ReleaseChecks(sb *pb.Surebet) {
	var wg sync.WaitGroup
	wg.Add(2)
	go h.ReleaseCheck(sb, 0, &wg)
	go h.ReleaseCheck(sb, 1, &wg)
	wg.Wait()
}
func (h *Handler) ReleaseCheck(sb *pb.Surebet, num int64, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	_, err := h.clients[sb.Members[num].ServiceName].ReleaseCheck(context.Background(), &pb.ReleaseCheckRequest{Surebet: SurebetWithOneMember(sb, num)})
	if err != nil {
		h.log.Error(err)
	}
}

var timeIdProducer = time_id.NewTimeId(1)

func (h *Handler) ProcessSurebet(sb *pb.Surebet) *SurebetError {
	ctx := context.Background()
	sb.SurebetId = timeIdProducer.GetId()
	start := util.UnixMsNow()
	if err := h.AllServicesActive(sb); err != nil {
		return err
	}
	if err := h.LoadConfig(ctx, sb); err != nil {
		return err
	}

	//if err := h.AnyDisabled(sb); err != nil {
	//	return err
	//}
	if err := h.GetCurrency(ctx, sb); err != nil {
		return err
	}
	checkCtx, cancel := context.WithTimeout(ctx, checkTimeOut)
	h.SendCheckLines(checkCtx, sb)
	cancel()
	defer h.ReleaseChecks(sb)

	if err := AllCheckStatus(sb); err != nil {
		return err
	}
	//h.LoadConfigForSub(ctx, sb)
	if err := h.Calc(sb); err != nil {
		return err
	}
	if err := h.AllSurebet(sb); err != nil {
		return err
	}
	h.log.Infow("begin_bet_first", "profit", sb.Calc.Profit, "stake", sb.Members[sb.Calc.FirstIndex].CheckCalc.Stake, "time", ElapsedFromId(start),
		"fid", sb.FortedSurebetId, "price", sb.Members[sb.Calc.FirstIndex].Check.Price, "name", sb.Members[sb.Calc.FirstIndex].ServiceName)
	h.SendPlaceFirst(ctx, sb)

	switch sb.Members[sb.Calc.FirstIndex].Bet.Status {
	case status.StatusOk:
		h.log.Infow("first_ok", "p", sb.Calc.Profit,
			"check_stake", sb.Members[sb.Calc.FirstIndex].CheckCalc.Stake,
			"bet_stake", sb.Members[sb.Calc.FirstIndex].Bet.Stake,
			"check_price", sb.Members[sb.Calc.FirstIndex].Check.Price,
			"bet_price", sb.Members[sb.Calc.FirstIndex].Bet.Price,
			"time", ElapsedFromId(start), "n", sb.Members[sb.Calc.FirstIndex].ServiceName, "fid", sb.FortedSurebetId)
		go h.SaveSurebet(sb)
	case status.AboveEventMax:
		return &SurebetError{Msg: status.AboveEventMax, Permanent: true, ServiceName: sb.Members[sb.Calc.FirstIndex].ServiceName}
	case status.MarketClosed:
		return &SurebetError{Msg: "market_closed", Permanent: true, ServiceName: sb.Members[sb.Calc.FirstIndex].ServiceName}
	default:
		return &SurebetError{Msg: "first_bet_not_ok", Permanent: false, ServiceName: sb.Members[sb.Calc.FirstIndex].ServiceName}
	}
	go h.ReleaseCheck(sb, sb.Calc.FirstIndex, nil)

	for i := 0; i < secondBetMaxTry; i++ {
		if sb.Members[sb.Calc.SecondIndex].Check.GetStatus() != status.StatusOk {
			time.Sleep(time.Millisecond * 100 * time.Duration(i))
			h.CheckLine(ctx, sb, sb.Calc.SecondIndex, nil)
			continue
		}
		isDone := h.CalcSecond(sb)
		if isDone {
			return nil
		}
		h.PlaceSecond(ctx, sb)
		if sb.Members[sb.Calc.SecondIndex].GetBet().GetStatus() == status.StatusOk {
			break
		} else if i < secondBetMaxTry {
			sb.Members[sb.Calc.SecondIndex].Check.Status = status.StatusError
		}
	}
	if sb.Members[sb.Calc.SecondIndex].GetBet().GetStatus() != status.StatusOk {
		go h.tel.Sendf("oblom_sur v=%v, p=%v, f=%v, s=%v, t=%v", betAmount(sb), sb.Calc.Profit, sb.Members[0].ServiceName, sb.Members[1].ServiceName, ElapsedFromId(start))
		return &SurebetError{Msg: "second_bet_not_ok", Permanent: true, ServiceName: sb.Members[sb.Calc.SecondIndex].ServiceName}
	}
	return nil
}
