package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"social-app/internal/config"
)

func Auth(cfg config.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		token, err := parseToken(tokenString, []byte(cfg.JWTSecret))
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		if err := injectClaimsIntoContext(c, claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := c.Query("token")
		if tokenString != "" {
			return tokenString
		}
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func parseToken(tokenString string, secret []byte) (*jwt.Token, error) {
	j, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return j, nil
}

func injectClaimsIntoContext(c *gin.Context, claims jwt.MapClaims) error {
	sub, ok := claims["sub"].(float64)
	if !ok {
		return fmt.Errorf("Missing 'sub' in token")
	}
	c.Set("user_id", uint64(sub))
	c.Set("role", getStringClaim(claims, "role", "anonymous"))
	c.Set("email", getStringClaim(claims, "email", ""))
	c.Set("username", getStringClaim(claims, "username", ""))
	c.Set("verified", getBoolClaim(claims, "verified", false))
	c.Set("need_2fa", getBoolClaim(claims, "need_2fa", false))
	return nil
}

func getStringClaim(claims jwt.MapClaims, key, fallback string) string {
	if val, exists := claims[key].(string); exists {
		return val
	}
	return fallback
}

func getBoolClaim(claims jwt.MapClaims, key string, fallback bool) bool {
	if val, exists := claims[key].(bool); exists {
		return val
	}
	return fallback
}
