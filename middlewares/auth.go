package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lacion/iothub/config"
)

// Strips 'Bearer ' prefix from bearer token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

// Auth is a basic middleware for authenticating requests
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := stripBearerPrefixFromTokenString(c.Request.Header.Get("Authorization"))

		if err != nil {
			c.AbortWithError(401, err)
		}

		cfg := config.Config()
		if token != cfg.GetString("secret") {
			c.AbortWithStatus(401)
		}
	}
}
