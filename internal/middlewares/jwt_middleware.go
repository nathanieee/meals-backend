package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func extractToken(c *gin.Context) (string, error) {
	var bearerToken = c.Request.Header.Get("Authorization")
	var err = errors.New("no Authorization token detected")

	// Apple already reserved header for Authorization
	// https://developer.apple.com/documentation/foundation/nsurlrequest
	if bearerToken == "" {
		bearerToken = c.Request.Header.Get("X-Authorization")
	}

	if len(strings.Split(bearerToken, " ")) == 2 {
		bearerToken = strings.Split(bearerToken, " ")[1]

		accessToken, err := c.Cookie("access_token")
		if err != nil {
			return "", err
		} else if accessToken == "" {
			return "", fmt.Errorf("you are not logged in")
		}
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
				Status:  consttypes.RST_ERROR,
				Message: "Invalid extract token",
				Data: utresponse.ErrorData{
					Debug:  err,
					Errors: err.Error(),
				},
			})
			ctx.Abort()
			return
		}

		parsedToken, err := uttoken.ParseToken(extractedToken, cfg.AccessToken.PublicKey)
		if err != nil {
			utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
				Status:  consttypes.RST_ERROR,
				Message: "Invalid parse token",
				Data: utresponse.ErrorData{
					Debug:  err,
					Errors: err.Error(),
				},
			})
			ctx.Abort()
			return
		}

		if !slices.Contains(allowedLevel, uint(consttypes.UR_USER)) || !slices.Contains(allowedLevel, uint(consttypes.UR_ADMIN)) {
			if !slices.Contains(allowedLevel, uint(parsedToken.User.Role)) || (consttypes.DateNow.Unix() >= parsedToken.Expires.Unix()) {
				utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
					Status:  consttypes.RST_ERROR,
					Message: "Invalid token",
					Data: utresponse.ErrorData{
						Debug:  nil,
						Errors: "You're not authorized to access this",
					},
				})
				ctx.Abort()
				return
			}
		}

		ctx.Set("user", *parsedToken.User)
		ctx.Set("access_token_uuid", parsedToken.TokenUUID.String())
		ctx.Next()
	}
}
