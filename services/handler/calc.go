package handler

import (
	"context"
	"fmt"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"math"
	"strings"
	"time"
)

var FastServiceList = []string{"bf", "mbook", "bdaq", "pin", "pin88", "isn", "sing2", "penta88"}
var SMList = []string{"bf", "isn", "sing2"}
var SlowServiceList = []string{"daf"}

// oly daf template
//const (
//	FastServiceSkewPercent     = 100.0
//	RepeatMultiplier           = 0.15
//	FastMaxBetDiff             = 20.0
//	ProfitFillFactorMultiplier = 1
//)

// sasha daf template
const (
	FastServiceSkewPercent     = 100.0
	RepeatMultiplier           = 0.05
	FastMaxBetDiff             = 25.0
	ProfitFillFactorMultiplier = 1
)

const (
	maxWinDiffPercent       = 8.0
	minGross                = 0.22
	startsTimeSpreadMinutes = 15.0
	FastVsPinSkew           = 0.0
	RoiFillFactorMultiplier = 1000.0
	ProfitTooLow            = "profit_too_low"
	ProfitTooHigh           = "profit_too_high"
	CountLineLimit          = "count_line_reached_max"
	CountEventLimit         = "count_event_reached_max"
	AmountEventLimit        = "amount_event_reached_max"
	AmountLineLimit         = "amount_line_reached_max"

	MaxStakeTooLow  = "max_stake_lower_min_stake"
	MinStakeTooHigh = "min_stake_higher_max_stake"

	WinDiffTooHigh = "win_diff_rel_too_high"
	GrossTooLow    = "gross_too_Low"
	ROITooLow      = "roi_too_Low"
)

func (h *Handler) CalcSecond(sb *pb.Surebet) (isDone bool) {
	var isSkew bool

	a := sb.Members[sb.Calc.FirstIndex]
	b := sb.Members[sb.Calc.SecondIndex]

	CalcMaxStake(b)
	CalcMinStake(b)
	CalcMaxWin(b)
	b.CheckCalc.Stake = a.Bet.Stake * a.Bet.Price / b.Check.Price

	if util.StringInList(a.Check.SubService, FastServiceList) && util.StringInList(b.Check.SubService, SlowServiceList) {
		h.log.Infow("fast_sub_service", "a_sub", a.Check.SubService, "b_stake", b.CheckCalc.Stake, "skew_sum", b.CheckCalc.Stake*FastServiceSkewPercent/100.0)
		b.CheckCalc.Stake = b.CheckCalc.Stake + b.CheckCalc.Stake*FastServiceSkewPercent/100.0
		isSkew = true
	} else if util.StringInList(b.Check.SubService, FastServiceList) && util.StringInList(a.Check.SubService, SlowServiceList) {
		h.log.Infow("fast_sub_service", "b_sub", b.Check.SubService, "b_stake", b.CheckCalc.Stake, "skew_sum", b.CheckCalc.Stake*FastServiceSkewPercent/100.0)
		b.CheckCalc.Stake = b.CheckCalc.Stake - b.CheckCalc.Stake*FastServiceSkewPercent/100.0
		isSkew = true
	} else if util.StringInList(a.Check.SubService, SMList) && b.Check.SubService == "pin" {
		skew := b.CheckCalc.Stake * FastVsPinSkew / 100.0
		stakeSkew := b.CheckCalc.Stake + skew
		h.log.Infow("fast_vs_pin",
			"first_sub", a.Check.SubService,
			"second_sub", b.Check.SubService,
			"second_stake", b.CheckCalc.Stake,
			"skew_sum", skew,
			"stake_with_skew", stakeSkew,
		)
		b.CheckCalc.Stake = stakeSkew

	} else if util.StringInList(b.Check.SubService, SMList) && a.Check.SubService == "pin" {
		skew := b.CheckCalc.Stake * FastVsPinSkew / 100.0
		stakeSkew := b.CheckCalc.Stake - skew
		h.log.Infow("fast_vs_pin", "first_sub", a.Check.SubService,
			"second_stake", b.CheckCalc.Stake,
			"skew_sum", skew,
			"stake_with_skew", stakeSkew,
		)
		b.CheckCalc.Stake = stakeSkew
	}

	if isSkew && b.CheckCalc.Stake < b.CheckCalc.MinStake {
		h.log.Info("bet_no_bet")
		b.ToBet = &pb.ToBet{Id: util.UnixMsNow()}
		b.Bet = &pb.Bet{
			Status:     status.StatusOk,
			StatusInfo: "bet_no_bet",
			Start:      util.UnixMsNow(),
			Done:       util.UnixMsNow(),
			Price:      b.Check.Price,
			Stake:      0,
			ApiBetId:   "0",
		}
		return true
	}

	if b.CheckCalc.Stake < b.CheckCalc.MinStake {
		b.CheckCalc.Stake = b.CheckCalc.MinStake
	} else if b.CheckCalc.Stake > b.CheckCalc.MaxStake {
		b.CheckCalc.Stake = b.CheckCalc.MaxStake
	}
	CalcWin(b)
	b.CheckCalc.Status = status.StatusOk
	return false
}

func (h *Handler) Calc(ctx context.Context, sb *pb.Surebet) *SurebetError {
	ctx, span := h.tracer.Start(ctx, "Calc")
	defer span.End()

	sb.Calc.Profit = Profit(sb)
	fStarts, err2 := h.ConvertStartTime(sb.Starts)
	if err2 != nil {
		return &SurebetError{Msg: fmt.Sprintf("parse_forted_starts_error:%v", sb.Starts)}
	}
	sb.Calc.HoursBeforeEvent = util.TruncateFloat(fStarts.Sub(time.Now()).Hours(), 2)

	sb.Calc.Roi = int64(sb.Calc.Profit * 24 / (sb.Calc.HoursBeforeEvent + 2) * 365)

	if sb.Members[0].Check.SubService == "daf" && util.StringInList(sb.Members[1].Check.SubService, SMList) {
		//lower:=sb.Members[0].Check.Price*sb.Members[0].Check.MaxBet/sb.Members[1].Check.Price*FastServiceSkewPercent/100
		//h.log.Infow("lower_max_bet", "max", sb.Members[1].Check.MaxBet, "lower", float64(10*sb.Members[0].Check.CountLine), "second_bet", lower)
		sb.Members[1].Check.MaxBet = sb.Members[1].Check.MaxBet - float64(FastMaxBetDiff*sb.Members[0].Check.CountLine)
	} else if sb.Members[1].Check.SubService == "daf" && util.StringInList(sb.Members[0].Check.SubService, SMList) {
		//lower:=sb.Members[1].Check.Price*sb.Members[1].Check.MaxBet/sb.Members[0].Check.Price*FastServiceSkewPercent/100
		//h.log.Infow("lower_max_bet", "max", sb.Members[0].Check.MaxBet, "lower", float64(10*sb.Members[1].Check.CountLine), "second_bet", lower)
		sb.Members[0].Check.MaxBet = sb.Members[0].Check.MaxBet - float64(FastMaxBetDiff*sb.Members[1].Check.CountLine)
	} else if sb.Members[0].Check.SubService == "pin" && util.StringInList(sb.Members[1].Check.SubService, SMList) {
		sb.Members[0].Check.MaxBet = sb.Members[0].Check.MaxBet - sb.Members[0].Check.MaxBet*FastVsPinSkew/100
	} else if sb.Members[1].Check.SubService == "pin" && util.StringInList(sb.Members[0].Check.SubService, SMList) {
		sb.Members[1].Check.MaxBet = sb.Members[1].Check.MaxBet - sb.Members[1].Check.MaxBet*FastVsPinSkew/100
	}

	var err *SurebetError
	for i := range sb.Members {
		m := sb.Members[i]
		m.CheckCalc = &pb.CheckCalc{}
		if m.Check.Price == 0 {
			return &SurebetError{Msg: "check_price_is_zero", ServiceName: m.ServiceName}
		}
		CalcMaxStake(m)
		CalcMinStake(m)
		CalcMaxWin(m)
		m.CheckCalc.Stake = m.CheckCalc.MaxStake
		CalcWin(m)
		m.CheckCalc.Status = status.StatusOk
		minProfit := util.TruncateFloat(m.BetConfig.MinPercent+m.Check.FillFactor*ProfitFillFactorMultiplier, 2)
		minRoi := m.BetConfig.MinRoi + int64(m.Check.FillFactor*RoiFillFactorMultiplier)

		if m.Check.SubService == "sbo" {
			m.BetConfig.Priority += 100
			m.BetConfig.RoundValue = 1
		} else if m.Check.SubService == "isn" {
			m.BetConfig.Priority += 5
			m.BetConfig.RoundValue = 1
		} else if m.Check.SubService == "pin" {
			m.BetConfig.Priority = 1
		} else if m.Check.SubService == "pin88" {
			m.BetConfig.Priority = 1
		} else if m.Check.SubService == "sing2" {
			m.BetConfig.RoundValue = 1
			m.BetConfig.Priority += 5
		} else if m.Check.SubService == "daf" {
			dafMinProfitCoefficient := float64(m.Check.CountLine) * RepeatMultiplier * FastServiceSkewPercent / 100
			minProfit = util.TruncateFloat(minProfit+dafMinProfitCoefficient, 2)
		}
		//h.log.Info("fill_factor: ", m.Check.FillFactor)

		if m.Check.CountLine >= m.BetConfig.MaxCountLine {
			m.CheckCalc.Status = CountLineLimit
			err = &SurebetError{Msg: fmt.Sprintf("%s:%v, max:%v, sub:%v", m.CheckCalc.Status, m.Check.CountLine, m.BetConfig.MaxCountLine, m.Check.SubService), Permanent: true, ServiceName: m.ServiceName}

		} else if m.Check.CountEvent >= m.BetConfig.MaxCountEvent {
			m.CheckCalc.Status = CountEventLimit
			err = &SurebetError{Msg: fmt.Sprintf("%s:%v, max:%v, sub:%v", m.CheckCalc.Status, m.Check.CountEvent, m.BetConfig.MaxCountEvent, m.Check.SubService), Permanent: true, ServiceName: m.ServiceName}

		} else if m.Check.AmountEvent >= m.BetConfig.MaxAmountEvent {
			m.CheckCalc.Status = AmountEventLimit
			err = &SurebetError{Msg: fmt.Sprintf("%s:%v, max:%v, sub:%v", m.CheckCalc.Status, m.Check.AmountEvent, m.BetConfig.MaxAmountEvent, m.Check.SubService), Permanent: true, ServiceName: m.ServiceName}
		} else if m.Check.AmountLine >= m.BetConfig.MaxAmountLine {
			m.CheckCalc.Status = AmountLineLimit
			err = &SurebetError{Msg: fmt.Sprintf("%s:%v, max:%v, sub:%v", m.CheckCalc.Status, m.Check.AmountLine, m.BetConfig.MaxAmountLine, m.Check.SubService), Permanent: true, ServiceName: m.ServiceName}

		} else if sb.Calc.Profit < minProfit {
			m.CheckCalc.Status = ProfitTooLow
			err = &SurebetError{Msg: fmt.Sprintf("%s:%.2f, min:%v, fill:%v", m.CheckCalc.Status, sb.Calc.Profit, minProfit, m.Check.FillFactor), Permanent: false, ServiceName: m.ServiceName}

		} else if sb.Calc.Profit > float64(m.BetConfig.MaxPercent) {
			m.CheckCalc.Status = ProfitTooHigh
			err = &SurebetError{Msg: fmt.Sprintf("%s:%.2f, max:%v", m.CheckCalc.Status, sb.Calc.Profit, m.BetConfig.MaxPercent), Permanent: false, ServiceName: m.ServiceName}

		} else if m.CheckCalc.MaxStake < m.CheckCalc.MinStake {
			m.CheckCalc.Status = MaxStakeTooLow
			err = &SurebetError{Msg: fmt.Sprintf("%s:%v, min:%v, s:%v", m.CheckCalc.Status, m.CheckCalc.MaxStake, m.CheckCalc.MinStake, m.SportName), Permanent: false, ServiceName: m.ServiceName}

		} else if m.CheckCalc.MinStake > m.CheckCalc.MaxStake {
			m.CheckCalc.Status = MinStakeTooHigh
			err = &SurebetError{Msg: fmt.Sprintf("%s:%v, max:%v", m.CheckCalc.Status, m.CheckCalc.MinStake, m.CheckCalc.MaxStake), Permanent: false, ServiceName: m.ServiceName}

		} else if sb.Calc.Roi < minRoi {
			m.CheckCalc.Status = ROITooLow
			err = &SurebetError{Msg: fmt.Sprintf("roi_too_low:%v, min:%v, ff:%v, hours:%v, p:%v", sb.Calc.Roi, minRoi, m.Check.FillFactor, sb.Calc.HoursBeforeEvent, sb.Calc.Profit), Permanent: false, ServiceName: m.ServiceName}
		}

		startTime, _ := h.ConvertStartTime(m.Starts)
		if math.Abs(fStarts.Sub(startTime).Minutes()) > startsTimeSpreadMinutes {
			m.CheckCalc.Status = "starts_diff_too_match"
			h.log.Infow("starts_diff", "service", m.ServiceName, "sStarts", m.Starts, "fStarts", fStarts, "diff", fStarts.Sub(startTime))
			err = &SurebetError{Msg: fmt.Sprintf("starts_diff_too_match:%v!=%v", fStarts, startTime), Permanent: true, ServiceName: m.ServiceName}
		}
		if err != nil && err.Permanent {
			return err
		}
	}
	if err != nil {
		return err
	}
	FirstSecond(sb)
	a := sb.Members[sb.Calc.LowerWinIndex]
	b := sb.Members[sb.Calc.HigherWinIndex]
	if util.StringInList(a.ServiceName, clients.SportMarketClones) && util.StringInList(b.ServiceName, clients.SportMarketClones) {
		h.log.Infow("oba_sportmarket", "f", sb.Calc.FirstName, "s", sb.Calc.SecondName, "a_event", a.EventId, "b_event", b.EventId)
		if a.EventId != b.EventId {
			return &SurebetError{Msg: fmt.Sprintf("diff_events, a:%v, b:%v, a_event:%v, b_event:%v", a.ServiceName, b.ServiceName, a.EventId, b.EventId), Permanent: true, ServiceName: b.ServiceName}
		}
	}
	if strings.Replace(a.Check.SubService, "pin88", "pin", -1) == strings.Replace(b.Check.SubService, "pin88", "pin", -1) {
		return &SurebetError{Msg: fmt.Sprintf("sub_services_equal, p:%v, a_sub:%v, b_sub:%v", sb.Calc.Profit, a.Check.SubService, b.Check.SubService), Permanent: true, ServiceName: b.ServiceName}
	}
	b.CheckCalc.Stake = CalcStake(a.CheckCalc.Win, b.Check.Price)

	if b.CheckCalc.Stake < b.CheckCalc.MinStake {
		b.CheckCalc.Stake = b.CheckCalc.MinStake
	}
	CalcWin(b)

	sb.Calc.WinDiff = CalcWinDiff(a.CheckCalc.Win, b.CheckCalc.Win)
	sb.Calc.WinDiffRel = CalcWinDiffRel(a.CheckCalc.Win, b.CheckCalc.Win)

	if sb.Calc.WinDiffRel > maxWinDiffPercent {
		b.CheckCalc.Status = WinDiffTooHigh
		//h.log.Infow(b.CheckCalc.Status, "aWin", a.CheckCalc.Win, "bWin", b.CheckCalc.Win, "aService", a.ServiceName, "bService", b.ServiceName)
		return &SurebetError{Msg: fmt.Sprintf("%s:%v, max:%v, %v:%v, %v:%v", b.CheckCalc.Status, sb.Calc.WinDiffRel, maxWinDiffPercent, a.ServiceName, a.CheckCalc.Win, b.ServiceName, b.CheckCalc.Win), Permanent: false, ServiceName: b.ServiceName}
	}
	sb.Calc.Gross = util.TruncateFloat(a.CheckCalc.Win*sb.Calc.Profit/100, 2)
	if sb.Calc.Gross < minGross {
		a.CheckCalc.Status = GrossTooLow
		return &SurebetError{Msg: fmt.Sprintf("gross_too_low:%v, min:%v", sb.Calc.Gross, minGross), Permanent: false, ServiceName: a.ServiceName}
	}
	h.log.Infow("check_calc_done", "a_check_calc", a.CheckCalc, "b_check_calc", b.CheckCalc, "calc", sb.Calc)
	return nil
}
