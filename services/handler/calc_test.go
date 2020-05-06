package handler

import (
	"github.com/aibotsoft/surebet-service/pkg/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinFloat(t *testing.T) {
	got := Min(0, 1)
	assert.Equal(t, float64(0), got)
	got = Min(0, 1, -1)
	assert.Equal(t, float64(-1), got)
	got = Min()
	assert.Equal(t, float64(0), got)

}

func TestHandler_Profit(t *testing.T) {
	got := Profit(tests.SurebetHelper(t))
	t.Log(got)
}
