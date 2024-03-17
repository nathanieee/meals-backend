package middlewares

import (
	"project-skbackend/configs"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func extractToken(c *gin.Context) (string, error) {
	tbearer := c.Request.Header.Get("Authorization")
	if tbearer == "" {
		tbearer = c.Request.Header.Get("X-Authorization")
	}

	if tbearer == "" {
		return "", utresponse.ErrTokenNotFound
	}

	splitToken := strings.Split(tbearer, " ")
	if len(splitToken) != 2 {
		return "", utresponse.ErrTokenInvalidFormat
	}

	tbearer = splitToken[1]

	taccess, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	} else if taccess == "" {
		return "", utresponse.ErrUserNotSignedIn
	}

	return tbearer, nil
}

func JWTAuthMiddleware(cfg *configs.Config, allowedlevel ...uint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		textract, err := extractToken(ctx)
		if err != nil {
			utresponse.GeneralUnauthorized(
				ctx,
				err,
			)
			ctx.Abort()
			return
		}

		tparsed, err := uttoken.ParseToken(textract, cfg.AccessToken.PublicKey)
		if err != nil {
			utresponse.GeneralUnauthorized(
				ctx,
				err,
			)
			ctx.Abort()
			return
		}

		if !slices.Contains(allowedlevel, uint(consttypes.UR_USER)) || !slices.Contains(allowedlevel, uint(consttypes.UR_ADMIN)) {
			if !slices.Contains(allowedlevel, uint(tparsed.User.Role)) || (consttypes.DateNow.Unix() >= tparsed.Expires.Unix()) {
				utresponse.GeneralUnauthorized(
					ctx,
					err,
				)
				ctx.Abort()
				return
			}
		}

		if !utrequest.CheckWhitelistUrl(ctx.Request.URL.Path) {
			if !tparsed.User.ConfirmedAt.IsZero() && !strings.Contains(ctx.Request.URL.Path, "verify") {
				utresponse.GeneralUnauthorized(
					ctx,
					err,
				)
				ctx.Abort()
				return
			}
		}

		ctx.Set("user", *tparsed.User)
		ctx.Set("access_token_uuid", tparsed.TokenUUID.String())
		ctx.Next()
	}
}
