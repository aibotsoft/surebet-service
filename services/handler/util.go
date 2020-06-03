package handler

import (
	"fmt"
	"github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
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
			return &SurebetError{Msg: "service has no config", Permanent: true, ServiceName: sb.Members[i].ServiceName}
		}
		if sb.Members[i].BetConfig.GetRegime() == status.StatusDisabled {
			return &SurebetError{Msg: "service disabled", Permanent: true, ServiceName: sb.Members[i].ServiceName}
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
			err = &SurebetError{Msg: "check_status: NotFound", Permanent: false, ServiceName: sb.Members[i].ServiceName}
		case status.StatusError:
			err = &SurebetError{Msg: "check_status: Error", Permanent: false, ServiceName: sb.Members[i].ServiceName}
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
		default:
			err = &SurebetError{Msg: "check_status not Ok", Permanent: false, ServiceName: sb.Members[i].ServiceName}
		}
		if err.Permanent {
			return err
		}
	}
	if sb.Members[0].Check.Status != status.StatusOk && sb.Members[1].Check.Status != status.StatusOk {
		return &SurebetError{Msg: "both_check_status_not_ok", Permanent: true, ServiceName: fmt.Sprintf("f: %q, s: %q", sb.Members[0].ServiceName, sb.Members[1].ServiceName)}
	}
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
		if sb.Members[i].CheckCalc.Status != status.StatusOk {
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

func ClearSurebet(sb *fortedpb.Surebet) {
	sb.Calc = nil
	for i := range sb.Members {
		sb.Members[i].Check = nil
		sb.Members[i].BetConfig = nil
		sb.Members[i].ToBet = nil
		sb.Members[i].Bet = nil
	}
}

func SurebetWithOneMember(sb *fortedpb.Surebet, i int) *fortedpb.Surebet {
	copySb := *sb
	copySb.Members = sb.Members[i : i+1]
	return &copySb
}
