package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/rs/zerolog/log"
)

func Trace(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.Trace().Enabled() {
			request, err := httputil.DumpRequest(r, false)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			log.Trace().Msg(string(request))
		}
		handler.ServeHTTP(w, r)
	})
}
