package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// loggerMiddleware is a middleware handler that does request logging
type loggerMiddleware struct {
	handler http.Handler
	logger  *logrus.Logger
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (lm *loggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lm.handler.ServeHTTP(w, r)
	lm.logger.Infof("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

// newLoggerMiddleware constructs a new Logger middleware handler
func newLoggerMiddleware(h http.Handler, l *logrus.Logger) *loggerMiddleware {
	return &loggerMiddleware{
		handler: h,
		logger:  l,
	}
}

// timerResponseMiddleware is a middleware response write to add 'X-Response-Time' to server response headers
type timerResponseMiddleware struct {
	http.ResponseWriter

	start time.Time
	ok    bool
}

// WriteHeader add 'X-Response-Time' header with value in microseconds
func (tm *timerResponseMiddleware) WriteHeader(statusCode int) {
	tm.Header().Set("X-Response-Time", strconv.FormatInt(time.Since(tm.start).Microseconds(), 10))
	tm.ResponseWriter.WriteHeader(statusCode)
	tm.ok = true
}

// Write write response with additional headers
func (tm *timerResponseMiddleware) Write(b []byte) (int, error) {
	if !tm.ok {
		tm.WriteHeader(http.StatusOK)
	}
	return tm.ResponseWriter.Write(b)
}

// headerMiddleware is a middleware handler that adds X-Server-Name and X-Response-Time
type headerMiddleware struct {
	handler http.Handler
}

// ServeHTTP handles the request and pass it to real handler
// adding response headers
func (hm *headerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("X-Server-Name", r.Host)
	hm.handler.ServeHTTP(&timerResponseMiddleware{
		ResponseWriter: w,
		ok:             false,
		start:          time.Now(),
	}, r)
}

// NewServerHeader constructs a new headerMiddleware middleware handler
func newHeaderMiddleware(h http.Handler) *headerMiddleware {
	return &headerMiddleware{handler: h}
}

// panicRecoveryMiddleware is a middleware handler that recover from panic and return InternalServer error
type panicRecoveryMiddleware struct {
	handler http.Handler
	logger  *logrus.Logger
}

// ServeHTTP recover from panic and return InternalServer error
// adding response headers
func (pcm *panicRecoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			pcm.logger.Errorf("panic recovery:%v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(&ErrorResponse{
				Error: fmt.Sprintf("panic recovered. Details:%v", err),
			}); err != nil {
				pcm.logger.Error(err.Error())
			}
		}
	}()

	pcm.handler.ServeHTTP(w, r)
}

// newPanicRecoveryMiddleware constructs a new panicRecoveryMiddleware middleware handler
func newPanicRecoveryMiddleware(h http.Handler, l *logrus.Logger) *panicRecoveryMiddleware {
	return &panicRecoveryMiddleware{
		handler: h,
		logger:  l,
	}
}
