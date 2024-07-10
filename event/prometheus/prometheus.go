package prometheus

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/emersonary/go-authentication/log/handler"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Call the next handler, which can be another middleware or the final handler
		next.ServeHTTP(w, r)

		status := w.(interface{ Status() int }).Status()

		logStr := r.URL.Path + " " + strconv.Itoa(status)

		log.LogHandler.AddLog(logStr, nil)
		// Log the details of the request
		fmt.Println("Log Prometheus " + logStr)

	})
}
