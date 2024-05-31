package main

import (
	"log"

	"github.com/alexsuriano/rate-limiter/config"
	"github.com/alexsuriano/rate-limiter/internal/infra/repository"
	"github.com/alexsuriano/rate-limiter/internal/infra/web"
	"github.com/alexsuriano/rate-limiter/internal/infra/web/webserver"
	"github.com/alexsuriano/rate-limiter/pkg/middlewares"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.NewConfig()

	webserver := webserver.NewWebServer(cfg.WebServerPort)

	cache := repository.NewRedisRepository(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB)

	middlewares := middlewares.NewMiddlewares(cache, middlewares.Options{
		LimitRequestIP:    cfg.LimitRequestIP,
		LimitRequestToken: cfg.LimitRequestToken,
		IpBlockingTime:    cfg.IpBlockingTime,
		TokenBlockingTime: cfg.TokenBlockingTime,
	})

	webserver.AddMiddleware(middlewares.RateLimiter)
	webserver.AddMiddleware(middleware.Logger)

	webserver.AddHandler("/", web.LoremIpsum)

	log.Printf("Starting %s on port %s", cfg.WebServerName, cfg.WebServerPort)
	webserver.Start()
}
