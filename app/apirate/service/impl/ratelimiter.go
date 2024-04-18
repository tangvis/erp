package impl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"math"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/tangvis/erp/app/apirate/repository"
	"github.com/tangvis/erp/app/apirate/service"
)

const (
	publicKey = "pub_key" // 如果没有配置限流，应用一个全局的
)

type Limiters struct {
	repo repository.Repo

	pool map[string]*Limiter // key为userID:apiPath

	lock sync.Mutex
}

func (l *Limiters) InitPublic(publicLimitSetting map[string]int) {
	// todo public for user rather than api path
	l.lock.Lock()
	defer l.lock.Unlock()
	for path, limit := range publicLimitSetting {
		l.pool[limiterID(publicKey, path)] = NewRateLimiter(publicKey, path, limit, math.MaxInt)
	}
}

func (l *Limiters) RateLimitWrapper(c *gin.Context) {
	success, allow := l.Allow("", c.Request.URL.Path)
	if !allow {
		c.String(http.StatusTooManyRequests, "too many requests")
		c.Abort()
		return
	}
	c.Next()
	// 只有成功返回才会扣减
	if c.Err() == nil {
		success()
	}
}

func (l *Limiters) GetPublicLimiter(path string) *Limiter {
	return l.pool[limiterID(publicKey, path)]
}

func (l *Limiters) Allow(userID, path string) (func(), bool) {
	var limiter *Limiter
	successAction := func() {
		if limiter != nil {
			limiter.Incr()
		}
	}
	// todo 读要不要加锁
	if limiter, ok := l.pool[limiterID(userID, path)]; ok {
		return successAction, limiter.Allow()
	}
	limiter = l.GetPublicLimiter(path)
	if limiter == nil {
		return successAction, false
	}
	return successAction, limiter.Allow()
}

func limiterID(userID, path string) string {
	return fmt.Sprintf("%s:%s", userID, path)
}

func NewLimiters(
	repo repository.Repo,
) service.APP {
	settings, err := repo.GetRateLimitSettings(context.Background())
	if err != nil {
		panic(err)
	}
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

	return &Limiters{
		pool: pool,
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
	TotalUsed    atomic.Uint64

	mu sync.Mutex
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(userID, apiPath string, qps, totalAllowed int) *Limiter {
	limiter := rate.NewLimiter(rate.Limit(qps), qps)
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
	if r.TotalUsed.Load() >= uint64(r.TotalAllowed) {
		return false
	}

	return r.limiter.Allow()
}

func (r *Limiter) Incr() {
	r.TotalUsed.Add(1)
}
