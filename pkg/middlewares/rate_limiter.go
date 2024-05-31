package middlewares

import (
	"net"
	"net/http"
)

type Middlewares struct {
	Cache   CacheInterface
	Options Options
}

type Options struct {
	LimitRequestIP    int
	LimitRequestToken int
	IpBlockingTime    int
	TokenBlockingTime int
	Cache             CacheInterface
}

func NewMiddlewares(cache CacheInterface, options Options) *Middlewares {
	return &Middlewares{
		Cache:   cache,
		Options: options,
	}
}

func (m *Middlewares) RateLimiter(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err, http.StatusInternalServerError)
			return
		}

		token := r.Header.Get("API_TOKEN")

	})
}
