package handler

import (
	"github.com/aibotsoft/decimal"
	pb "github.com/aibotsoft/gen/fortedpb"
)

func Profit(sb *pb.Surebet) float64 {
	var prob float64
	for i := range sb.Members {
		prob += probability(sb.Members[i].Check.Price)
	}
	profit := 1/prob*100 - 100
	return RoundDown(profit, 0.001)
}

func probability(price float64) float64 {
	return 1 / price
}

func RoundDown(value float64, roundValue float64) float64 {
	rv := decimal.NewFromFloat(roundValue)
	res, _ := decimal.NewFromFloat(value).Div(rv).Floor().Mul(rv).Float64()
	return res
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
}
func CalcMaxStake(m *pb.SurebetSide) *SurebetError {
	if m.Check.Price == 0 {
		return &SurebetError{Msg: "Check.Price is 0", ServiceName: m.ServiceName}
	}
	//m.CheckCalc.MaxStake = Min(float64(m.Check.Balance), m.Check.MaxBet, float64(m.BetConfig.MaxStake), float64(m.BetConfig.MaxWin)/m.Check.Price)
	m.CheckCalc.MaxStake, _ = decimal.Min(
		decimal.New(m.Check.Balance, 0),
		decimal.NewFromFloat(m.Check.MaxBet),
		decimal.New(m.BetConfig.MaxStake, 0),
		decimal.New(m.BetConfig.MaxWin, 0).DivRound(decimal.NewFromFloat(m.Check.Price), 3)).Float64()
	return nil
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
