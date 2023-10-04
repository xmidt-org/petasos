package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap/zaptest"
)

func TestSanitizeHeaders(t *testing.T) {
	testCases := []struct {
		Description string
		Input       http.Header
		Expected    http.Header
	}{
		{
			Description: "Filtered",
			Input:       http.Header{"Authorization": []string{"Basic xyz"}, "HeaderA": []string{"x"}},
			Expected:    http.Header{"HeaderA": []string{"x"}, "Authorization-Type": []string{"Basic"}},
		},
		{
			Description: "Handled human error",
			Input:       http.Header{"Authorization": []string{"BasicXYZ"}, "HeaderB": []string{"y"}},
			Expected:    http.Header{"HeaderB": []string{"y"}},
		},
		{
			Description: "Not a perfect system",
			Input:       http.Header{"Authorization": []string{"MySecret IWantToLeakIt"}},
			Expected:    http.Header{"Authorization-Type": []string{"MySecret"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			assert := assert.New(t)
			actual := sanitizeHeaders(tc.Input)
			assert.Equal(tc.Expected, actual)
		})

	}
}

func TestAddFieldToLog_NotNil(t *testing.T) {
	testCases := []struct {
		name string
		kvs  []interface{}
	}{
		{
			name: "header",
			kvs:  []interface{}{"headerName", "petasos"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx, cancel := context.WithCancel(context.Background())
			defer cancel()
			testLogger := zaptest.NewLogger(t)

			assert := assert.New(t)
			ctx := addFieldsToLog(testCtx, testLogger, tc.kvs)

			assert.NotNil(sallust.Get(ctx))

		})
	}

}
