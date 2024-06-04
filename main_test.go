package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alexsuriano/rate-limiter/config"
	"github.com/alexsuriano/rate-limiter/internal/infra/repository"
	"github.com/alexsuriano/rate-limiter/internal/infra/web"
	"github.com/alexsuriano/rate-limiter/internal/infra/web/webserver"
	"github.com/alexsuriano/rate-limiter/internal/limiter"
	"github.com/alexsuriano/rate-limiter/pkg/middlewares"
	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestTokenRequest(t *testing.T) {

	t.Run("When execute one request", func(t *testing.T) {
		cfg := config.NewConfig()
		webserver := webserver.NewWebServer(":8181")

		cache, err := repository.NewRedisRepository(
			"localhost",
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDB)
		assert.NoError(t, err)

		limit := limiter.NewLimiter(cache, limiter.Options{
			LimitRequestIP:    1,
			LimitRequestToken: 1,
			IpBlockingTime:    1,
			TokenBlockingTime: 1,
		})

		webserver.AddMiddleware(middlewares.RateLimiter(limit))
		webserver.AddMiddleware(middleware.Logger)

		webserver.AddHandler("/", web.LoremIpsum)
		go webserver.Start()
		time.Sleep(50 * time.Millisecond)

		url := fmt.Sprintf("http://localhost:8181")
		request, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		request.Header.Set("API_KEY", "abc1")

		response, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("When execute request until block", func(t *testing.T) {
		cfg := config.NewConfig()
		webserver := webserver.NewWebServer(":8282")

		cache, err := repository.NewRedisRepository(
			"localhost",
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDB)
		assert.NoError(t, err)

		limit := limiter.NewLimiter(cache, limiter.Options{
			LimitRequestIP:    1,
			LimitRequestToken: 10,
			IpBlockingTime:    1,
			TokenBlockingTime: 2,
		})

		webserver.AddMiddleware(middlewares.RateLimiter(limit))
		webserver.AddMiddleware(middleware.Logger)

		webserver.AddHandler("/", web.LoremIpsum)
		go webserver.Start()
		time.Sleep(50 * time.Millisecond)

		url := fmt.Sprintf("http://localhost:8282")
		request, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		request.Header.Set("API_KEY", "abc2")

		for i := range 30 {
			response, err := http.DefaultClient.Do(request)
			assert.NoError(t, err)
			defer response.Body.Close()

			if i < 10 {
				assert.Equal(t, http.StatusOK, response.StatusCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
			}
		}
	})
}

func TestIPRequest(t *testing.T) {

	t.Run("When execute request until block", func(t *testing.T) {
		cfg := config.NewConfig()
		webserver := webserver.NewWebServer(":8383")

		cache, err := repository.NewRedisRepository(
			"localhost",
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDB)
		assert.NoError(t, err)

		limit := limiter.NewLimiter(cache, limiter.Options{
			LimitRequestIP:    10,
			LimitRequestToken: 1,
			IpBlockingTime:    2,
			TokenBlockingTime: 1,
		})

		webserver.AddMiddleware(middlewares.RateLimiter(limit))
		webserver.AddMiddleware(middleware.Logger)

		webserver.AddHandler("/", web.LoremIpsum)
		go webserver.Start()
		time.Sleep(50 * time.Millisecond)

		url := fmt.Sprintf("http://localhost:8383")
		request, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		for i := range 30 {
			response, err := http.DefaultClient.Do(request)
			assert.NoError(t, err)
			defer response.Body.Close()

			if i < 10 {
				assert.Equal(t, http.StatusOK, response.StatusCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
			}
		}
	})
}
