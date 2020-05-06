package store

import (
	"context"
	"database/sql"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/cache"
	"github.com/aibotsoft/micro/config"
	"github.com/dgraph-io/ristretto"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Store struct {
	cfg   *config.Config
	log   *zap.SugaredLogger
	db    *sqlx.DB
	cache *ristretto.Cache
}

func NewStore(cfg *config.Config, log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{log: log, db: db, cache: cache.NewCache(cfg)}
}
func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		s.log.Error(err)
	}
	s.cache.Close()
}

func MemberByName(sb *pb.Surebet, serviceName string) *pb.SurebetSide {
	for i := 0; i < len(sb.Members); i++ {
		if sb.Members[i].ServiceName == serviceName {
			return sb.Members[i]
		}
	}
	return nil
}
func MemberNameList(sb *pb.Surebet) []string {
	var nameList []string
	for i := 0; i < len(sb.Members); i++ {
		nameList = append(nameList, sb.Members[i].ServiceName)
	}
	return nameList
}

func Filter(arr []pb.BetConfig, name string) *pb.BetConfig {
	for i := range arr {
		if arr[i].ServiceName == name {
			return &arr[i]
		}
	}
	return nil
}

func (s *Store) LoadConfig(ctx context.Context, sb *pb.Surebet) error {
	var betConfigs []pb.BetConfig
	query, args, err := sqlx.In("SELECT * FROM dbo.BetConfig WHERE ServiceName IN (?)", MemberNameList(sb))
	err = s.db.SelectContext(ctx, &betConfigs, s.db.Rebind(query), args...)
	if err != nil {
		return errors.Wrapf(err, "select BetConfig error")
	}
	for i := range sb.Members {
		sb.Members[i].BetConfig = Filter(betConfigs, sb.Members[i].ServiceName)
	}
	return nil
}

func (s *Store) SaveFortedSurebet(sb *pb.Surebet) error {
	_, err := s.db.Exec("uspSaveFortedSurebet",
		sql.Named("CreatedAt", sb.CreatedAt),
		sql.Named("Starts", sb.Starts),
		sql.Named("FortedHome", sb.FortedHome),
		sql.Named("FortedAway", sb.FortedAway),
		sql.Named("FortedProfit", sb.FortedProfit),
		sql.Named("FortedSport", sb.FortedSport),
		sql.Named("FortedLeague", sb.FortedLeague),
		sql.Named("FilterName", sb.FilterName),
		sql.Named("FortedSurebetId", sb.FortedSurebetId))
	if err != nil {
		return errors.Wrap(err, "uspSaveFortedSurebet error")
	}
	return nil
}

func (s *Store) SaveCalc(sb *pb.Surebet) error {
	_, err := s.db.Exec("uspSaveCalc",
		sql.Named("Profit", sb.Calc.Profit),
		sql.Named("FirstName", sb.Calc.FirstName),
		sql.Named("SecondName", sb.Calc.SecondName),
		sql.Named("LowerWinIndex", sb.Calc.LowerWinIndex),
		sql.Named("HigherWinIndex", sb.Calc.HigherWinIndex),
		sql.Named("FirstIndex", sb.Calc.FirstIndex),
		sql.Named("SecondIndex", sb.Calc.SecondIndex),
		sql.Named("WinDiff", sb.Calc.WinDiff),
		sql.Named("WinDiffRel", sb.Calc.WinDiffRel),
		sql.Named("FortedSurebetId", sb.FortedSurebetId),
		sql.Named("SurebetId", sb.SurebetId),
	)
	if err != nil {
		return errors.Wrap(err, "uspSaveCalc error")
	}
	return nil
}

func (s *Store) SaveSide(sb *pb.Surebet) error {
	for i, side := range sb.Members {
		side.GetCheck()
		_, err := s.db.Exec("uspSaveSide",
			sql.Named("SurebetId", sb.SurebetId),
			sql.Named("SideIndex", i),
			sql.Named("ServiceName", side.ServiceName),
			sql.Named("SportName", side.SportName),
			sql.Named("LeagueName", side.LeagueName),
			sql.Named("Home", side.Home),
			sql.Named("Away", side.Away),
			sql.Named("MarketName", side.MarketName),
			sql.Named("Price", side.Price),
			sql.Named("Initiator", side.Initiator),
			sql.Named("Starts", sb.Starts),
			sql.Named("EventId", side.EventId),

			sql.Named("CheckId", side.Check.Id),
			sql.Named("AccountId", side.Check.AccountId),
			sql.Named("AccountLogin", side.Check.AccountLogin),
			sql.Named("CheckStatus", side.Check.Status),
			sql.Named("StatusInfo", side.Check.StatusInfo),
			sql.Named("CountLine", side.Check.CountLine),
			sql.Named("CountEvent", side.Check.CountEvent),
			sql.Named("AmountEvent", side.Check.AmountEvent),
			sql.Named("MinBet", side.Check.MinBet),
			sql.Named("MaxBet", side.Check.MaxBet),
			sql.Named("Balance", side.Check.Balance),
			sql.Named("CheckPrice", side.Check.Price),
			sql.Named("Currency", side.Check.Currency),
			sql.Named("CheckDone", side.Check.Done),

			sql.Named("CalcStatus", side.GetCheckCalc().GetStatus()),
			sql.Named("MaxStake", side.GetCheckCalc().GetMaxStake()),
			sql.Named("MinStake", side.GetCheckCalc().GetMinStake()),
			sql.Named("MaxWin", side.GetCheckCalc().GetMaxWin()),
			sql.Named("Stake", side.GetCheckCalc().GetStake()),
			sql.Named("Win", side.GetCheckCalc().GetWin()),
			sql.Named("IsFirst", side.GetCheckCalc().GetIsFirst()),

			sql.Named("ToBetId", side.GetToBet().GetId()),
			sql.Named("TryCount", side.GetToBet().GetTryCount()),

			sql.Named("BetStatus", side.GetBet().GetStatus()),
			sql.Named("BetStatusInfo", side.GetBet().GetStatusInfo()),
			sql.Named("Start", side.GetBet().GetStart()),
			sql.Named("Done", side.GetBet().GetDone()),
			sql.Named("BetPrice", side.GetBet().GetPrice()),
			sql.Named("BetStake", side.GetBet().GetStake()),
			sql.Named("ApiBetId", side.GetBet().GetApiBetId()),
		)
		if err != nil {
			return errors.Wrap(err, "uspSaveCalc error")
		}
	}
	return nil
}