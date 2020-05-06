package tests

import (
	pb "github.com/aibotsoft/gen/fortedpb"
	"testing"
	"time"
)

func SurebetHelper(t *testing.T) *pb.Surebet {
	t.Helper()
	return &pb.Surebet{
		CreatedAt:       time.Time{}.Format(time.RFC3339Nano),
		Starts:          time.Time{}.Format(time.RFC3339Nano),
		FortedHome:      "FortedHome",
		FortedAway:      "FortedAway",
		FortedProfit:    6.66,
		FortedSport:     "FortedSport",
		FortedLeague:    "FortedLeague",
		FilterName:      "FilterName",
		FortedSurebetId: 0,
		Calc: &pb.Calc{
			Profit:         0,
			FirstName:      "Sbobet",
			SecondName:     "Pinnacle",
			LowerWinIndex:  0,
			HigherWinIndex: 1,
			FirstIndex:     0,
			SecondIndex:    1,
			WinDiff:        0,
			WinDiffRel:     0,
		},
		Members: []*pb.SurebetSide{{
			Num:         1,
			ServiceName: "Sbobet",
			SportName:   "Table Tennis",
			LeagueName:  "Ukraine Win Cup - Men's Singles (Set Handicap)",
			Home:        "Mikhail Varchenko",
			Away:        "Ilya Glivenko",
			MarketName:  "П1",
			Price:       2.05,
			Url:         "https://www.sbobet.com/ru-ru/euro/table-tennis/ukraine-win-cup---men-s-singles-(set-handicap)/2973913/mikhail-varchenko-vs-ilya-glivenko",
			Initiator:   true,
			Check: &pb.Check{
				Price:       1.35,
				CountLine:   0,
				CountEvent:  0,
				AmountEvent: 0,
				MinBet:      3,
				MaxBet:      3,
				Status:      "Ok",
				StatusInfo:  "",
				Balance:     100,
			},
		}, {Num: 2,
			ServiceName: "Pinnacle",
			SportName:   "TestSport_2",
			LeagueName:  "TestLeague_2",
			Home:        "TestHome_2",
			Away:        "TestAway_2",
			MarketName:  "TestMarket_2",
			Price:       2.05,
			Url:         "http://testurl.loc",
			Initiator:   false,
			Check: &pb.Check{
				Price:       4.25,
				CountLine:   0,
				CountEvent:  0,
				AmountEvent: 0,
				MinBet:      1,
				MaxBet:      5,
				Status:      "Ok",
				StatusInfo:  "",
				Balance:     100,
			},
		}}}
}
