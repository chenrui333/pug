package tui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgo(t *testing.T) {
	now := time.Now()

	assert.Equal(t, "47s ago", ago(now, now.Add(-47*time.Second)))
}