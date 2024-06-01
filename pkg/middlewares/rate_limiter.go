package middlewares

import (
	"net"
	"net/http"

	"github.com/alexsuriano/rate-limiter/internal/limiter"
)

func RateLimiter(limit limiter.LimiterInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			token := r.Header.Get("API_KEY")

			var overLimit bool

			if token != "" {
				overLimit, err = limit.OverLimit(token, limiter.Token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				overLimit, err = limit.OverLimit(ip, limiter.IP)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			if overLimit == true {
				errorMessage := "you have reached the maximum number of requests or actions allowed within a certain time frame"
				http.Error(w, errorMessage, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
