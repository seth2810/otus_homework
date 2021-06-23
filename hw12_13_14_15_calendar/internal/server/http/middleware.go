package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/app"
)

type statusCodeCatcher struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeCatcher) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler, logger app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		scc := &statusCodeCatcher{w, http.StatusOK}

		next.ServeHTTP(scc, r)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		logger.Info(fmt.Sprintf("%s [%s] %s %s %s %d %d %q",
			ip, time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method, r.URL.Path, r.Proto, scc.statusCode, time.Since(start), r.UserAgent(),
		))
	})
}
