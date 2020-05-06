package handler

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/surebet-service/pkg/clients"
	"github.com/aibotsoft/surebet-service/pkg/store"
	"github.com/aibotsoft/surebet-service/pkg/tests"
	"testing"
	"time"
)

func InitHelper(t *testing.T) *Handler {
	t.Helper()
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := store.NewStore(cfg, log, db)
	cli := clients.NewClients(log)
	h := NewHandler(cfg, log, sto, cli)
	return h
}

func SurebetToProcess(t *testing.T) *pb.Surebet {
	t.Helper()
	return &pb.Surebet{
		CreatedAt: time.Time{}.Format(time.RFC3339Nano),
		//Starts:       time.Time{}.Format(time.RFC3339Nano),
		Starts:       "2020-05-02T10:28:01.568",
		FortedHome:   "FortedHome",
		FortedAway:   "FortedAway",
		FortedProfit: 6.66,
		Members: []*pb.SurebetSide{{
			Num:         1,
			ServiceName: "Sbobet",
			SportName:   "E Sports",
			LeagueName:  "FIFA 20 - Russia Liga Pro (12 Mins)",
			Home:        "EZ1D 11 (EZ1)",
			Away:        "Gambit Esports (GMB)",
			MarketName:  "ТМ(3,25)",
			Price:       2.26,
			Url:         "https://www.sbobet.com/ru-ru/euro/e-sports/fifa-20---russia-liga-pro-(12-mins)/2985196/ez1d-11-(ez1)-vs-gambit-esports-(gmb)",
			Initiator:   true,
		}, {Num: 2,
			ServiceName: "Pinnacle",
			SportName:   "Soccer",
			LeagueName:  "eSoccer - Liga Pro (12 mins)",
			Home:        "Ez1d (EZ1) Esports",
			Away:        "Gambit (GMB) Esports",
			MarketName:  "ТБ(3,25)",
			Price:       1.862,
			Url:         "https://members.pinnacle.com/Sportsbook/Mobile/ru-RU/Enhanced/Regular/SportsBookAll/35/Curacao/Odds/Soccer-29/Market/2/208851/1122224338",
			Initiator:   false,
		}}}
}

func TestHandler_SurebetLoop(t *testing.T) {
	h := InitHelper(t)
	sur := SurebetToProcess(t)
	h.SurebetLoop(sur)
}

func TestHandler_CheckLine(t *testing.T) {
	h := InitHelper(t)
	sur := tests.SurebetHelper(t)
	h.CheckLine(context.Background(), sur, 0, nil)
}

func TestSurebetWithOneMember(t *testing.T) {
	sur := SurebetToProcess(t)
	got := SurebetWithOneMember(sur, 1)
	t.Log(got)
}
