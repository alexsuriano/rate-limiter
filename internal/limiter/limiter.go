package limiter

import "time"

type RegisterType int

const (
	Token = 1
	IP    = 2
)

type Limiter struct {
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

func NewLimiter(cache CacheInterface, options Options) *Limiter {
	return &Limiter{
		Cache:   cache,
		Options: options,
	}
}

func (l *Limiter) OverLimit(register string, registerType RegisterType) (bool, error) {

	limitRequest, blockingTime := l.getConfig(registerType)

	currentCount, err := l.Cache.Get(register)
	if err != nil {
		return false, err
	}

	if currentCount >= limitRequest {
		return true, nil
	}

	_, err = l.Cache.Increment(register, time.Duration(blockingTime*int(time.Second)))
	if err != nil {
		return false, err
	}

	return false, nil
}

func (l *Limiter) getConfig(registerType RegisterType) (limitRequest, blockingTime int) {
	switch registerType {
	case Token:
		return l.Options.LimitRequestToken, l.Options.TokenBlockingTime
	case IP:
		return l.Options.LimitRequestIP, l.Options.IpBlockingTime
	default:
		return 0, 0
	}
}
