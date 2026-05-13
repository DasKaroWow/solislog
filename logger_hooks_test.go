// hooks_test.go
package solislog_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestBeforeHookCanModifyRecord(t *testing.T) {
	var buf bytes.Buffer
	logger := solislog.NewLogger(nil, solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{message} | {extra[hooked]}\n",
		BeforeHook: func(record *solislog.Record) {
			record.Message = "changed"
			record.Extra["hooked"] = "yes"
		},
	}))

	logger.Info("original")

	t.Log(buf.String())

	want := "changed | yes\n"
	got := buf.String()
	require.Equal(t, want, got, "BeforeHook should modify record before rendering")
}

func TestAfterHookReceivesRenderedMessage(t *testing.T) {
	var buf bytes.Buffer
	var gotMessage string
	var gotRendered []byte

	logger := solislog.NewLogger(nil, solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{level} | {message}\n",
		AfterHook: func(record *solislog.Record, msg []byte, successful bool) {
			gotMessage = record.Message
			gotRendered = msg
		},
	}))

	logger.Info("hello")

	want := "INFO | hello\n"
	got := buf.String()

	t.Log(got)

	require.Equal(t, want, got, "output should match expected template")
	assert.Equal(t, "hello", gotMessage, "hook should receive original message")
	assert.Equal(t, want, string(gotRendered), "hook should receive fully rendered output")
}

func TestBeforeHookIsIsolatedPerHandler(t *testing.T) {
	var first, second bytes.Buffer
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&first, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{message}\n",
			BeforeHook: func(record *solislog.Record) {
				record.Message = "changed"
			},
		}),
		solislog.NewHandler(&second, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{message}\n",
		}),
	)

	logger.Info("original")

	t.Log("first:", first.String())
	t.Log("second:", second.String())

	require.Equal(t, "changed\n", first.String(), "first handler should see modified record")
	require.Equal(t, "original\n", second.String(), "second handler should see unmodified record")
}

func TestAfterHookRunsAfterUnlock(t *testing.T) {
	var buf bytes.Buffer
	var logger *solislog.Logger

	logger = solislog.NewLogger(nil, solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{message}\n",
		AfterHook: func(record *solislog.Record, msg []byte, successful bool) {
			if strings.TrimSpace(string(msg)) == "first" {
				logger.Info("second")
			}
		},
	}))

	logger.Info("first")

	want := "first\nsecond\n"
	got := buf.String()

	t.Log(got)

	require.Equal(t, want, got, "AfterHook should run after handler unlock, allowing nested log calls")
}

func TestErrorHandlerIsCalledOnWriteError(t *testing.T) {
	var gotErr error
	var gotMsg []byte
	var gotRecord *solislog.Record

	logger := solislog.NewLogger(nil, solislog.NewHandler(&failingWriter{}, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{level} | {message}\n",
		ErrorHandler: func(record *solislog.Record, msg []byte, err error) {
			gotErr = err
			gotMsg = msg
			gotRecord = record
		},
	}))

	logger.Info("hello")

	require.NotNil(t, gotErr, "ErrorHandler should be called when writer fails")
	assert.Equal(t, "write failed", gotErr.Error())
	assert.Equal(t, "INFO | hello\n", string(gotMsg))
	assert.Equal(t, "hello", gotRecord.Message)
}

func TestErrorHandlerRunsAfterUnlock(t *testing.T) {
	var callCount int
	var logger *solislog.Logger

	logger = solislog.NewLogger(nil, solislog.NewHandler(&failingWriter{}, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{message}\n",
		ErrorHandler: func(record *solislog.Record, msg []byte, err error) {
			callCount++
			if callCount == 1 {
				logger.Info("nested")
			}
		},
	}))

	logger.Info("first")
	assert.Equal(t, 2, callCount, "ErrorHandler should allow nested log calls without deadlock")
}

func TestAfterHookReceivesCorrectSuccessfulFlag(t *testing.T) {
	// 1. Успешная запись → successful == true
	var successFlag bool
	var successBuf bytes.Buffer
	loggerSuccess := solislog.NewLogger(nil, solislog.NewHandler(&successBuf, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{message}\n",
		AfterHook: func(record *solislog.Record, msg []byte, successful bool) {
			successFlag = successful
		},
	}))
	loggerSuccess.Info("ok")

	t.Log("success output:", successBuf.String())
	require.True(t, successFlag, "AfterHook should receive successful=true on normal write")

	// 2. Ошибка записи → successful == false
	var failFlag bool
	loggerFail := solislog.NewLogger(nil, solislog.NewHandler(&failingWriter{}, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{message}\n",
		AfterHook: func(record *solislog.Record, msg []byte, successful bool) {
			failFlag = successful
		},
	}))
	loggerFail.Info("fail")

	require.False(t, failFlag, "AfterHook should receive successful=false on write error")
}
