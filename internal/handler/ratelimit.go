package handler

import (
	"myproject/internal/logger"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

// у каждого IP свой limiter
type ipLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	lim, exists := l.limiters[ip]
	if !exists {
		lim = rate.NewLimiter(l.r, l.b)
		l.limiters[ip] = lim
	}
	return lim
}

// 10 запросов в секунду, burst 20
var limiter = newIPLimiter(10, 20)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.get(ip).Allow() {
			logger.Log.Warn().Str("ip", ip).Msg("rate limit exceeded")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}
