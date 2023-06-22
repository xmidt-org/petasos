package main

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/xmidt-org/candlelight"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

// LoggerFunc is a strategy for adding key/value pairs (possibly) based on an HTTP request.
// Functions of this type must append key/value pairs to the supplied slice and then return
// the new slice.
type LoggerFunc func([]interface{}, *http.Request) []interface{}

func sanitizeHeaders(headers http.Header) (filtered http.Header) {
	filtered = headers.Clone()
	if authHeader := filtered.Get("Authorization"); authHeader != "" {
		filtered.Del("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 {
			filtered.Set("Authorization-Type", parts[0])
		}
	}
	return
}

func setLogger(logger *zap.Logger, lf ...LoggerFunc) func(delegate http.Handler) http.Handler {

	if logger == nil {
		panic("The base Logger cannot be nil")
	}

	return func(delegate http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				kvs := []interface{}{"requestHeaders", sanitizeHeaders(r.Header), "requestURL", r.URL.EscapedPath(), "method", r.Method}
				for _, f := range lf {
					if f != nil {
						kvs = f(kvs, r)
					}
				}
				kvs, _ = candlelight.AppendTraceInfo(r.Context(), kvs)
				ctx := r.Context()
				ctx = addFieldsToLog(ctx, logger, kvs)
				delegate.ServeHTTP(w, r.WithContext(ctx))
			})
	}
}

func addFieldsToLog(ctx context.Context, logger *zap.Logger, kvs []interface{}) context.Context {

	for i := 0; i <= len(kvs)-2; i += 2 {
		logger = logger.With(zap.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
	}

	return sallust.With(ctx, logger)

}

func header(headerName, keyName string) LoggerFunc {
	headerName = textproto.CanonicalMIMEHeaderKey(headerName)

	return func(kv []interface{}, request *http.Request) []interface{} {
		values := request.Header[headerName]
		switch len(values) {
		case 0:
			return append(kv, keyName, "")
		case 1:
			return append(kv, keyName, values[0])
		default:
			return append(kv, keyName, values)
		}
	}
}
