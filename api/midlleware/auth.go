package midlleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	product_service "github.com/Andrewalifb/alpha-pos-system-product-service/api/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type JWTPayloadWithClaims struct {
	*product_service.JWTPayload
}

func (p *JWTPayloadWithClaims) Valid() error {
	if p.StandardClaims.ExpiresAt < time.Now().Unix() {
		return errors.New("token is expired")
	}
	expectedIssuer := os.Getenv("ISSUER")
	if p.StandardClaims.Issuer != expectedIssuer {
		return errors.New("invalid issuer")
	}
	// Add more validation if needed
	return nil
}

// JWTAuthMiddleware function
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load(".env")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error occurred"})
			c.Abort()
			return
		}
		var secretKey = os.Getenv("SECRET_KEY")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken[1], &JWTPayloadWithClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTPayloadWithClaims); ok && token.Valid {
			c.Set("user", claims.JWTPayload)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
