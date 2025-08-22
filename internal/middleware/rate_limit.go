package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		now := time.Now()
		windowStart := now.Add(-rl.window)

		// Clean old requests
		if times, exists := rl.requests[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if t.After(windowStart) {
					validTimes = append(validTimes, t)
				}
			}
			rl.requests[clientIP] = validTimes
		}

		// Check if limit exceeded
		if len(rl.requests[clientIP]) >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		// Add current request
		rl.requests[clientIP] = append(rl.requests[clientIP], now)

		c.Next()
	}
}
