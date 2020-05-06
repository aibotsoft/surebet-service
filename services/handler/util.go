package handler

import (
	"github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
)

func (h *Handler) AllServicesActive(sb *fortedpb.Surebet) error {
	for i := 0; i < len(sb.Members); i++ {
		if h.clients[sb.Members[i].ServiceName] == nil {
			return SurebetError{nil, "service " + sb.Members[i].ServiceName + " not active", true}
		}
	}
	return nil
}
func (h *Handler) AllSurebet(sb *fortedpb.Surebet) bool {
	for i := 0; i < len(sb.Members); i++ {
		if sb.Members[i].BetConfig.Regime != status.RegimeSurebet {
			h.log.Infow("regime not Surebet", "regime", sb.Members[i].BetConfig.Regime, "name", sb.Members[i].ServiceName)
			return false
		}
	}
	return true
}

func AllCheckCalcStatusOk(sb *fortedpb.Surebet) bool {
	for i := 0; i < len(sb.Members); i++ {
		if sb.Members[i].CheckCalc.Status != status.StatusOk {
			return false
		}
	}
	return true
}
func AllCheckStatusOk(sb *fortedpb.Surebet) bool {
	for i := 0; i < len(sb.Members); i++ {
		if sb.Members[i].Check.Status != status.StatusOk {
			return false
		}
	}
	return true
}
func (h *Handler) AnyDisabled(sb *fortedpb.Surebet) error {
	for i := range sb.Members {
		if sb.Members[i].BetConfig.Regime == status.StatusDisabled {
			return SurebetError{nil, "service " + sb.Members[i].ServiceName + " disabled", true}
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
