package apirate

import (
	"fmt"
	"golang.org/x/time/rate"
	"sync"

	"github.com/tangvis/erp/app/apirate/repository"
)

const (
	publicKey = "pub_key" // 如果没有配置限流，应用一个全局的
)

type Limiters struct {
	Pool map[string]*Limiter // key为userID:apiPath

	lock sync.Mutex
}

func (l *Limiters) GetPublicLimiter(path string) *Limiter {
	return l.Pool[limiterID(publicKey, path)]
}

func (l *Limiters) Allow(userID, path string) bool {
	// todo 读要不要加锁
	if limiter, ok := l.Pool[limiterID(userID, path)]; ok {
		return limiter.Allow()
	}
	pubLimiter := l.GetPublicLimiter(path)
	if pubLimiter == nil {
		return false
	}
	return pubLimiter.Allow()
}

func limiterID(userID, path string) string {
	return fmt.Sprintf("%s:%s", userID, path)
}

func NewLimiters(settings []repository.RateSetting) *Limiters {
	pool := make(map[string]*Limiter)
	for _, setting := range settings {
		if !setting.Valid() {
			continue
		}
		pool[limiterID(setting.UserID, setting.Path)] = NewRateLimiter(
			setting.UserID,
			setting.Path,
			setting.QPSLimit,
			setting.TotalLimit,
		)
	}
	for _, router := range public.AllRouters {
		_ = router
	}

	return &Limiters{
		Pool: pool,
		lock: sync.Mutex{},
	}
}

// Limiter holds the limiter settings and the limiter itself.
type Limiter struct {
	UserID       string
	APIPath      string
	QPS          int // Queries per second
	TotalAllowed int
	limiter      *rate.Limiter
	TotalUsed    int

	mu sync.Mutex
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(userID, apiPath string, qps, totalAllowed int) *Limiter {
	limiter := rate.NewLimiter(rate.Limit(qps), 1) // Bucket size of 1 to smooth the rate limiting
	return &Limiter{
		UserID:       userID,
		APIPath:      apiPath,
		QPS:          qps,
		TotalAllowed: totalAllowed,
		limiter:      limiter,
	}
}

// Allow checks if a request is allowed under the rate limiting rules.
func (r *Limiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.TotalUsed >= r.TotalAllowed {
		return false
	}

	if r.limiter.Allow() {
		r.TotalUsed++
		return true
	}

	return false
}
