package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/codepnw/microservice-ecommerce/token"
	"github.com/gin-gonic/gin"
)

const claimsKey string = "claims"

func GetAuthMiddlewareFunc(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// read the authorization header
		// verify the token
		claims, err := verifyClaimsFromAuthHeader(c, tokenMaker)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// pass the payload/claims down the context
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func GetAdminMiddlewareFunc(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// read the authorization header
		// verify the token
		claims, err := verifyClaimsFromAuthHeader(c, tokenMaker)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !claims.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is not admin"})
			c.Abort()
			return
		}

		// pass the payload/claims down the context
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func verifyClaimsFromAuthHeader(c *gin.Context, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := fields[1]
	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
