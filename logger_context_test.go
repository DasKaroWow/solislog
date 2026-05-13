package solislog_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/assert"
)

func TestContextualizeAndFromContext(t *testing.T) {
	base := solislog.NewLogger(nil, solislog.NewHandler(io.Discard, solislog.InfoLevel, nil))
	ctx := base.Contextualize(context.Background(), solislog.Extra{"trace_id": "abc"})

	logger, ok := solislog.FromContext(ctx)
	assert.NotSame(t, base, logger)
	assert.True(t, ok, "FromContext should find the logger")
	assert.NotNil(t, logger, "retrieved logger should not be nil")
}

func TestFromContextReturnsFalseOnEmptyContext(t *testing.T) {
	_, ok := solislog.FromContext(context.Background())
	assert.False(t, ok, "FromContext should return false for empty context")
}

func TestContextualizeDoesNotMutateParentContext(t *testing.T) {
	base := solislog.NewLogger(nil, solislog.NewHandler(io.Discard, solislog.InfoLevel, nil))
	parent := context.Background()
	child := base.Contextualize(parent, solislog.Extra{"k": "v"})

	_, okInParent := solislog.FromContext(parent)
	_, okInChild := solislog.FromContext(child)

	assert.False(t, okInParent, "parent context should not contain logger")
	assert.True(t, okInChild, "child context should contain logger")
}

func TestContextualizeChainsExtraFields(t *testing.T) {
	var buf bytes.Buffer
	base := solislog.NewLogger(
		solislog.Extra{"l1": "a"},
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{extra[l1]} {extra[l2]} {extra[l3]} | {message}\n",
		}),
	)

	ctx1 := base.Contextualize(context.Background(), solislog.Extra{"l2": "b"})
	logger1, _ := solislog.FromContext(ctx1)

	ctx2 := logger1.Contextualize(ctx1, solislog.Extra{"l3": "c"})
	logger2, _ := solislog.FromContext(ctx2)

	logger2.Info("chain test")
	assert.Equal(t, buf.String(), "a b c | chain test\n")
}
