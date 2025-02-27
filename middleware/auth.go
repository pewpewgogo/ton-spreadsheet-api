package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret string

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Get the JWT secret from the environment variable
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable not set")
	}
}

var (
	memoryStore = struct {
		sync.RWMutex
		data map[string]*RequestLimiter
	}{data: make(map[string]*RequestLimiter)}
)

type RequestLimiter struct {
	count       int
	lastReset   time.Time
	limit       int
	resetPeriod time.Duration
}

func NewRequestLimiter(limit int, period time.Duration) *RequestLimiter {
	return &RequestLimiter{
		count:       0,
		lastReset:   time.Now(),
		limit:       limit,
		resetPeriod: period,
	}
}

func (rl *RequestLimiter) isAllowed() bool {
	if time.Since(rl.lastReset) > rl.resetPeriod {
		rl.count = 0
		rl.lastReset = time.Now()
	}
	if rl.count < rl.limit {
		rl.count++
		return true
	}
	return false
}

func rateLimiter(key string, limit int, period time.Duration) bool {
	memoryStore.Lock()
	defer memoryStore.Unlock()

	limiter, exists := memoryStore.data[key]
	if !exists {
		limiter = NewRequestLimiter(limit, period)
		memoryStore.data[key] = limiter
	}

	return limiter.isAllowed()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		var tokenLimit int
		var tokenKey string

		if tokenString == "" {
			tokenLimit = 10
			tokenKey = "no_token_" + c.ClientIP()
		} else {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			tokenLimit = 100
			tokenKey = "token_" + tokenString // Key for users with token
		}

		if !rateLimiter(tokenKey, tokenLimit, time.Second) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
