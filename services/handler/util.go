package handler

import (
	"context"
	"fmt"
	"github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
	"github.com/jinzhu/copier"
	"time"
)

func (h *Handler) AllServicesActive(sb *fortedpb.Surebet) *SurebetError {
	for i := range sb.Members {
		if h.clients[sb.Members[i].ServiceName] == nil {
			return &SurebetError{Msg: "service_not_active", Permanent: true, ServiceName: sb.Members[i].ServiceName}
		}
	}
	return nil
}
func (h *Handler) AnyDisabled(sb *fortedpb.Surebet) *SurebetError {
	for i := range sb.Members {
		if sb.Members[i].BetConfig == nil {
			return &SurebetError{Msg: "service_has_no_config", Permanent: true, ServiceName: sb.Members[i].ServiceName}
		}
		if sb.Members[i].BetConfig.GetRegime() == status.StatusDisabled {
			return &SurebetError{Msg: "service_disabled", Permanent: true, ServiceName: sb.Members[i].ServiceName}
		}
	}
	return nil
}
func AnyHasName(sb *fortedpb.Surebet, serviceName string) bool {
	for i := range sb.Members {
		if sb.Members[i].ServiceName == serviceName {
			return true
		}
	}
	return false
}
func betAmount(sb *fortedpb.Surebet) int {
	var amount float64
	for i := range sb.Members {
		amount = amount + sb.Members[i].GetBet().GetStake()
	}
	return int(amount)
}

func AllCheckStatusOk(sb *fortedpb.Surebet) *SurebetError {
	for i := 0; i < len(sb.Members); i++ {
		if sb.Members[i].Check.Status != status.StatusOk {
			return &SurebetError{Msg: "check status not Ok", Permanent: false, ServiceName: sb.Members[i].ServiceName}
		}
	}
	return nil
}
func AllCheckStatus(sb *fortedpb.Surebet) *SurebetError {
	var err *SurebetError
	for i := 0; i < len(sb.Members); i++ {
		switch sb.Members[i].Check.Status {
		case status.StatusOk:
			continue
		case status.StatusNotFound:
			err = &SurebetError{Msg: fmt.Sprintf("not_found, i:%v", sb.Members[i].Check.StatusInfo), Permanent: false, ServiceName: sb.Members[i].ServiceName}
		case status.ServiceBusy:
			err = &SurebetError{Msg: fmt.Sprintf("service_busy"), Permanent: false, ServiceName: sb.Members[i].ServiceName}
		case status.StatusError:
			err = &SurebetError{Msg: fmt.Sprintf("error, i:%v", sb.Members[i].Check.StatusInfo), Permanent: false, ServiceName: sb.Members[i].ServiceName}
		case status.BadBettingStatus:
			err = &SurebetError{Msg: status.BadBettingStatus, Permanent: true, ServiceName: sb.Members[i].ServiceName}
		case status.ServiceSportDisabled:
			err = &SurebetError{Msg: status.ServiceSportDisabled, Permanent: true, ServiceName: sb.Members[i].ServiceName}
		case status.PitchersRequired:
			err = &SurebetError{Msg: status.PitchersRequired, Permanent: true, ServiceName: sb.Members[i].ServiceName}
		case status.ServiceError:
			err = &SurebetError{Msg: status.ServiceError, Permanent: true, ServiceName: sb.Members[i].ServiceName}
		case status.MarketClosed:
			err = &SurebetError{Msg: status.MarketClosed, Permanent: true, ServiceName: sb.Members[i].ServiceName}
		case status.PriceTooLow:
			err = &SurebetError{Msg: status.PriceTooLow, Permanent: true, ServiceName: sb.Members[i].ServiceName}
		case status.Suspended:
			err = &SurebetError{Msg: status.Suspended, Permanent: false, ServiceName: sb.Members[i].ServiceName}
		case status.HandicapChanged:
			err = &SurebetError{Msg: fmt.Sprintf("handicap_changed, i:%v", sb.Members[i].Check.StatusInfo), Permanent: false, ServiceName: sb.Members[i].ServiceName}
		default:
			err = &SurebetError{Msg: "check_status_not_ok", Permanent: false, ServiceName: sb.Members[i].ServiceName}
		}
		if err.Permanent {
			return err
		}
	}
	//if sb.Members[0].Check.Status != status.StatusOk && sb.Members[1].Check.Status != status.StatusOk {
	//	return &SurebetError{Msg: "both_check_status_not_ok", Permanent: true, ServiceName: fmt.Sprintf("f:%s, s:%s", sb.Members[0].ServiceName, sb.Members[1].ServiceName)}
	//}
	return err
}

func (h *Handler) AllSurebet(sb *fortedpb.Surebet) *SurebetError {
	for i := range sb.Members {
		if sb.Members[i].BetConfig.Regime != status.RegimeSurebet {
			return &SurebetError{Msg: "regime_not_Surebet", Permanent: true, ServiceName: sb.Members[i].ServiceName}
		}
	}
	return nil
}
func AllCheckCalcStatusOk(sb *fortedpb.Surebet) *SurebetError {
	for i := range sb.Members {
		if sb.Members[i].CheckCalc.GetStatus() != status.StatusOk {
			return &SurebetError{Msg: "CheckCalc status not Ok", Permanent: false, ServiceName: sb.Members[i].ServiceName}
		}
	}
	return nil
}

func Min(vn ...float64) float64 {
	var m float64
	for i, e := range vn {
		if i == 0 || e < m {
			m = e
		}
	}
	return m
}
func Max(vn ...float64) float64 {
	var m float64
	for i, e := range vn {
		if i == 0 || e > m {
			m = e
		}
	}
	return m
}
func ElapsedFromSurebetId(surebetId int64) int64 {
	return (util.UnixUsNow() - surebetId) / 1000
}
func ElapsedFromId(value int64) int64 {
	return util.UnixMsNow() - value
}

func ClearSurebet(sb *fortedpb.Surebet) {
	sb.Calc = fortedpb.Calc{}
	for i := range sb.Members {
		sb.Members[i].BetConfig = nil
		sb.Members[i].Check = nil
		sb.Members[i].CheckCalc = nil
		sb.Members[i].ToBet = nil
		sb.Members[i].Bet = nil
	}
}

func SurebetWithOneMember(sb *fortedpb.Surebet, i int64) *fortedpb.Surebet {
	copySb := *sb
	copySb.Members = sb.Members[i : i+1]
	return &copySb
}
func (h *Handler) LoadConfig(ctx context.Context, sb *fortedpb.Surebet) *SurebetError {
	for i := range sb.Members {
		conf, err := h.store.GetConfigByName(ctx, sb.Members[i].ServiceName)
		if err != nil {
			return &SurebetError{Err: err, Msg: "GetConfigByName_error", Permanent: true}
		}
		sb.Members[i].BetConfig = &conf
	}
	return nil
}
func (h *Handler) LoadConfigForSub(ctx context.Context, sb *fortedpb.Surebet) {
	for i := range sb.Members {
		conf, err := h.store.GetConfigBySub(ctx, sb.Members[i].Check.SubService)
		if err == nil {
			sb.Members[i].BetConfig = &conf
		}
	}
}
func (h *Handler) GetCurrency(ctx context.Context, sb *fortedpb.Surebet) *SurebetError {
	get, b := h.store.Cache.Get("currency_list")
	if b {
		sb.Currency = get.([]fortedpb.Currency)
		return nil
	}
	currency, err := h.Conf.GetCurrency(ctx)
	if err != nil {
		return &SurebetError{Err: err, Msg: "get_currency_error", Permanent: true}
	}
	err = copier.Copy(&sb.Currency, &currency)
	if err != nil {
		return &SurebetError{Err: err, Msg: "copy_currency_slice_error", Permanent: true}
	}
	h.store.Cache.SetWithTTL("currency_list", sb.Currency, 1, time.Hour)
	return nil
}

func (h *Handler) HasAnyBet(sb *fortedpb.Surebet) bool {
	for i := range sb.Members {
		bet := sb.Members[i].GetBet()
		if bet != nil {
			if bet.GetStatus() == status.StatusOk {
				return true
			}
		}
	}
	return false
}

func (h *Handler) ConvertStartTime(startTime string) (time.Time, error) {
	starts, err := time.Parse(util.ISOFormat, startTime)
	if err == nil {
		return starts, nil
	}
	starts, err = time.Parse(time.RFC3339, startTime)
	if err == nil {
		return starts, nil
	}
	return time.Time{}, err
}
