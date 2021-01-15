package handler

import (
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/util"
	"github.com/shopspring/decimal"
)

func Profit(sb *pb.Surebet) (prob float64) {
	for i := range sb.Members {
		prob += 1 / sb.Members[i].Check.Price
	}
	profit := 1/prob*100 - 100
	return util.TruncateFloat(profit, 3)
}
func ActualProfit(sb *pb.Surebet) (prob float64) {
	for i := range sb.Members {
		prob += 1 / sb.Members[i].Bet.Price
	}
	profit := 1/prob*100 - 100
	return util.TruncateFloat(profit, 2)
}

func FirstSecond(sb *pb.Surebet) {
	if sb.Members[0].BetConfig.Priority >= sb.Members[1].BetConfig.Priority {
		sb.Calc.FirstIndex = 0
		sb.Calc.SecondIndex = 1
	} else {
		sb.Calc.FirstIndex = 1
		sb.Calc.SecondIndex = 0
	}
	sb.Calc.FirstName = sb.Members[sb.Calc.FirstIndex].ServiceName
	sb.Calc.SecondName = sb.Members[sb.Calc.SecondIndex].ServiceName
	sb.Members[sb.Calc.FirstIndex].CheckCalc.IsFirst = true
	if sb.Members[0].CheckCalc.MaxWin <= sb.Members[1].CheckCalc.MaxWin {
		sb.Calc.LowerWinIndex = 0
		sb.Calc.HigherWinIndex = 1
	} else {
		sb.Calc.LowerWinIndex = 1
		sb.Calc.HigherWinIndex = 0
	}
}
func CalcMaxStake(m *pb.SurebetSide) {
	m.CheckCalc.MaxStake, _ = decimal.Min(
		decimal.New(m.Check.Balance, 0),
		decimal.NewFromFloat(m.Check.MaxBet),
		decimal.New(m.BetConfig.MaxStake, 0),
		decimal.New(m.BetConfig.MaxWin, 0).DivRound(decimal.NewFromFloat(m.Check.Price), 3)).Float64()
}
func CalcMinStake(m *pb.SurebetSide) {
	m.CheckCalc.MinStake = Max(m.Check.MinBet, float64(m.BetConfig.MinStake))
}
func CalcMaxWin(m *pb.SurebetSide) {
	m.CheckCalc.MaxWin, _ = decimal.NewFromFloat(m.CheckCalc.MaxStake).Mul(decimal.NewFromFloat(m.Check.Price)).Float64()
}
func CalcWin(m *pb.SurebetSide) {
	m.CheckCalc.Win, _ = decimal.NewFromFloat(m.CheckCalc.Stake).Mul(decimal.NewFromFloat(m.Check.Price)).Round(5).Float64()
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
