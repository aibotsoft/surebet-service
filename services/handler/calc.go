package handler

import (
	"github.com/aibotsoft/decimal"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
)

const maxWinDiffPercent = 5

func (h *Handler) CalcSecond(sb *pb.Surebet) {
	a := sb.Members[sb.Calc.FirstIndex]
	b := sb.Members[sb.Calc.SecondIndex]
	CalcMaxStake(b)
	CalcMinStake(b)
	CalcMaxWin(b)
	b.CheckCalc.Stake = a.Bet.Stake * a.Bet.Price / b.Check.Price

	if b.CheckCalc.Stake < b.CheckCalc.MinStake {
		b.CheckCalc.Stake = b.CheckCalc.MinStake
	}
	CalcWin(b)
	b.CheckCalc.Status = status.StatusOk
}
func (h *Handler) Calc(sb *pb.Surebet) {
	sb.Calc = &pb.Calc{Profit: Profit(sb)}
	for i := range sb.Members {
		m := sb.Members[i]
		m.CheckCalc = &pb.CheckCalc{}
		CalcMaxStake(m)
		CalcMinStake(m)
		CalcMaxWin(m)
		m.CheckCalc.Stake = m.CheckCalc.MaxStake
		CalcWin(m)
		m.CheckCalc.Status = status.StatusOk

		if m.Check.Status != status.StatusOk {
			m.CheckCalc.Status = "Check.Status not Ok"
			h.log.Infow(m.CheckCalc.Status, "status", m.Check.Status, "StatusInfo", m.Check.StatusInfo, "name", m.ServiceName)

		} else if sb.Calc.Profit < m.BetConfig.MinPercent {
			m.CheckCalc.Status = "Profit lower then MinPercent"
			h.log.Infow(m.CheckCalc.Status, "profit", sb.Calc.Profit, "MinPercent", m.BetConfig.MinPercent, "name", m.ServiceName)

		} else if sb.Calc.Profit > float64(m.BetConfig.MaxPercent) {
			m.CheckCalc.Status = "Profit higher then MaxPercent"
			h.log.Infow(m.CheckCalc.Status, "profit", sb.Calc.Profit, "MaxPercent", m.BetConfig.MaxPercent, "name", m.ServiceName)

		} else if m.Check.CountLine >= m.BetConfig.MaxCountLine {
			m.CheckCalc.Status = "CountLine has reached MaxCountLine"
			h.log.Infow(m.CheckCalc.Status, "CountLine", m.Check.CountLine, "MaxCountLine", m.BetConfig.MaxCountLine, "name", m.ServiceName)

		} else if m.Check.CountEvent >= m.BetConfig.MaxCountEvent {
			m.CheckCalc.Status = "CountEvent has reached MaxCountEvent"
			h.log.Infow(m.CheckCalc.Status, "CountEvent", m.Check.CountEvent, "MaxCountEvent", m.BetConfig.MaxCountEvent, "name", m.ServiceName)

		} else if m.Check.AmountEvent >= m.BetConfig.MaxAmountEvent {
			m.CheckCalc.Status = "AmountEvent has reached MaxAmountEvent"
			h.log.Infow(m.CheckCalc.Status, "AmountEvent", m.Check.AmountEvent, "MaxAmountEvent", m.BetConfig.MaxAmountEvent, "name", m.ServiceName)

		} else if m.CheckCalc.MaxStake < m.CheckCalc.MinStake {
			m.CheckCalc.Status = "MaxStake lower CheckCalc.MinStake"
			h.log.Infow(m.CheckCalc.Status, "MaxStake", m.CheckCalc.MaxStake, "MinStake", m.CheckCalc.MinStake, "name", m.ServiceName)

		} else if m.CheckCalc.MinStake > m.CheckCalc.MaxStake {
			m.CheckCalc.Status = "MinStake higher CheckCalc.MaxStake"
			h.log.Infow(m.CheckCalc.Status, "MinStake", m.CheckCalc.MinStake, "MaxStake", m.CheckCalc.MaxStake, "name", m.ServiceName)

		}
	}
	if !AllCheckCalcStatusOk(sb) {
		return
	}
	FirstSecond(sb)
	if sb.Members[0].CheckCalc.MaxWin <= sb.Members[1].CheckCalc.MaxWin {
		//левая сторона слабее, ее ставка базовая
		sb.Calc.LowerWinIndex = 0
		sb.Calc.HigherWinIndex = 1
	} else {
		sb.Calc.LowerWinIndex = 1
		sb.Calc.HigherWinIndex = 0
	}
	a := sb.Members[sb.Calc.LowerWinIndex]
	b := sb.Members[sb.Calc.HigherWinIndex]

	b.CheckCalc.Stake = CalcStake(a.CheckCalc.Win, b.Check.Price)

	if b.CheckCalc.Stake < b.CheckCalc.MinStake {
		b.CheckCalc.Stake = b.CheckCalc.MinStake
	}
	CalcWin(b)

	sb.Calc.WinDiff = CalcWinDiff(a.CheckCalc.Win, b.CheckCalc.Win)
	sb.Calc.WinDiffRel = CalcWinDiffRel(a.CheckCalc.Win, b.CheckCalc.Win)

	if sb.Calc.WinDiffRel > maxWinDiffPercent {
		b.CheckCalc.Status = "WinDiffRel too high"
		h.log.Infow(b.CheckCalc.Status, "WinDiffRel", sb.Calc.WinDiffRel, "WinDiff", sb.Calc.WinDiff, "name", b.ServiceName)
	}
	h.log.Infow("check calc done", "a_check_calc", a.CheckCalc, "b_check_calc", b.CheckCalc, "calc", sb.Calc)
}

func CalcStake(aWin float64, bPrice float64) float64 {
	f, _ := decimal.NewFromFloat(aWin).DivRound(decimal.NewFromFloat(bPrice), 5).Float64()
	return f
}
func CalcWinDiff(aWin float64, bWin float64) float64 {
	f, _ := decimal.NewFromFloat(aWin).Sub(decimal.NewFromFloat(bWin)).Abs().Float64()
	return f
}

var d100 = decimal.New(100, 0)

func CalcWinDiffRel(aWin float64, bWin float64) float64 {
	aWinD := decimal.NewFromFloat(aWin)
	bWinD := decimal.NewFromFloat(bWin)
	sumWinD := aWinD.Add(bWinD)
	res, _ := aWinD.Sub(bWinD).Abs().Mul(d100).DivRound(sumWinD, 2).Float64()
	return res
}
