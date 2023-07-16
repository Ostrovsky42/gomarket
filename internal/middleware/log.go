package middleware

import (
	"net/http"
	"time"

	"gomarket/internal/logger"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := loggingResponseWriter{
			ResponseWriter: w,
			respData:       &responseData{},
		}
		h.ServeHTTP(&lw, r)

		logger.Log.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Int("status", lw.respData.status).
			Str("duration", time.Since(start).String()).
			Int("size", lw.respData.size).
			Send()
	}
	return http.HandlerFunc(logFn)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		respData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.respData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.respData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
