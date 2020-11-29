package messaging

import (
	"context"
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewBrokerContext(t *testing.T) {
	ctx := context.Background()
	brokerCtx := NewBrokerContext(ctx)
	assert.Equal(t, ctx, brokerCtx.Context(), "context should be equals")
}
