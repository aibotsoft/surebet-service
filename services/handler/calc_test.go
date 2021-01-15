package handler

import (
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
