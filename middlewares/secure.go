package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	stsHeader           = "Strict-Transport-Security"
	stsSubdomainString  = "; includeSubdomains"
)

func process(w http.ResponseWriter, r *http.Request) error {

	// SSL checks
	if strings.EqualFold(r.URL.Scheme, "https") || r.TLS != nil {
		url := r.URL
		url.Scheme = "https"
		url.Host = r.Host

		status := http.StatusMovedPermanently

		http.Redirect(w, r, url.String(), status)
		return fmt.Errorf("Redirecting to HTTPS")
	}


	// 315360000s are 10 years
	w.Header().Add(stsHeader, fmt.Sprintf("max-age=%d%s", 315360000, stsSubdomainString))

	return nil

}

// Secure adds some extra security when using TLS
func Secure() gin.HandlerFunc {

	return func(c *gin.Context) {
		err := process(c.Writer, c.Request)
		if err != nil {
			if c.Writer.Written() {
				c.AbortWithStatus(c.Writer.Status())
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}

}
