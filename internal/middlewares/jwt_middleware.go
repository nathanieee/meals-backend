package middlewares

import (
	"errors"
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func extractToken(c *gin.Context) (string, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	err := errors.New("no Authorization token detected")

	// Apple already reserved header for Authorization
	// https://developer.apple.com/documentation/foundation/nsurlrequest
	if bearerToken == "" {
		bearerToken = c.Request.Header.Get("X-Authorization")
	}

	if len(strings.Split(bearerToken, " ")) == 2 {
		bearerToken = strings.Split(bearerToken, " ")[1]
	}

	if bearerToken == "" {
		return "", err
	}

	return bearerToken, nil
}

func JWTAuthMiddleware(cfg *configs.Config, allowedLevel ...uint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		extractedToken, err := extractToken(ctx)
		if err != nil {
			utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
				Message: "Invalid extract token",
				Debug:   err,
				Errors:  err.Error(),
			})
			ctx.Abort()
			return
		}

		parsedToken, err := uttoken.ParseToken(extractedToken, cfg.App.Secret)
		if err != nil {
			utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
				Message: "Invalid parse token",
				Debug:   err,
				Errors:  err.Error(),
			})
			ctx.Abort()
			return
		}

		if !slices.Contains(allowedLevel, uint(consttypes.UR_USER)) {
			if !slices.Contains(allowedLevel, uint(parsedToken.User.Role)) || (time.Now().Unix() >= parsedToken.Expire) {
				utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
					Message: "Invalid token",
					Debug:   nil,
					Errors:  "You're not authorized to access this",
				})
				ctx.Abort()
				return
			}
		}

		if !utrequest.CheckWhitelistUrl(ctx.Request.URL.Path) {
			if parsedToken.User.ConfirmedAt == (time.Time{}) && !strings.Contains(ctx.Request.URL.Path, "verify") {
				utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
					Message: "Invalid token",
					Debug:   nil,
					Errors:  "This account is not verified",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Set("user", *parsedToken.User)
		ctx.Next()
	}
}
