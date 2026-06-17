package handler

import (
	"myproject/internal/logger"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// у каждого IP свой limiter + время последнего запроса
type ipLimiter struct {
	entries map[string]*entry
	mu      sync.Mutex
	r       rate.Limit
	b       int
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	l := &ipLimiter{
		entries: make(map[string]*entry),
		r:       r,
		b:       b,
	}
	// чистим мёртвые IP каждые 5 минут
	go l.cleanup(5 * time.Minute)
	return l
}

func (l *ipLimiter) cleanup(interval time.Duration) {
	for {
		time.Sleep(interval)
		l.mu.Lock()
		for ip, e := range l.entries {
			if time.Since(e.lastSeen) > interval {
				delete(l.entries, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, exists := l.entries[ip]
	if !exists {
		e = &entry{limiter: rate.NewLimiter(l.r, 