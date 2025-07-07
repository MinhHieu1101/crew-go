package graphql

import (
	"context"
	"strings"
	"time"

	"user-service/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type claimsKey struct{}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		var accessClaims *auth.AccessClaims
		if parts := strings.SplitN(header, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
			if ac, err := auth.ParseAccessToken(parts[1]); err == nil {
				accessClaims = ac
			}
		}

		if accessClaims == nil {
			if cookie, err := c.Request.Cookie("refresh_token"); err == nil {
				if rc, err2 := auth.ParseRefreshToken(cookie.Value); err2 == nil {
					// TODO: implement token renewal
					accessClaims = &auth.AccessClaims{
						UserID: rc.UserID,
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
						},
					}
					// review this new access token
					// newAT, _ := auth.GenerateAccessToken(rc.UserID)
				}
			}
		}

		if accessClaims != nil {
			ctx := context.WithValue(c.Request.Context(), claimsKey{}, accessClaims)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// FromContext returns (*auth.AccessClaims, true) if AuthMiddleware ran
func FromContext(ctx context.Context) (*auth.AccessClaims, bool) {
	ac, ok := ctx.Value(claimsKey{}).(*auth.AccessClaims)
	return ac, ok
}
