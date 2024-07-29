package middlewares

import (
	"project-skbackend/configs"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func extractToken(ctx *gin.Context) (string, error) {
	tbearer := ctx.GetHeader("Authorization")
	if tbearer == "" {
		tbearer = ctx.GetHeader("X-Authorization")
	}

	if tbearer == "" {
		return "", consttypes.ErrTokenNotFound
	}

	splitToken := strings.Split(tbearer, " ")
	if len(splitToken) != 2 {
		return "", consttypes.ErrTokenInvalidFormat
	}

	tbearer = splitToken[1]

	taccess := ctx.GetHeader(consttypes.T_ACCESS)
	if taccess == "" {
		return "", consttypes.ErrUserNotSignedIn
	}

	return tbearer, nil
}

func JWTAuthMiddleware(cfg *configs.Config, allowedlevel ...consttypes.UserRole) gin.HandlerFunc {
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

		tparsed, err := uttoken.ParseToken(textract, cfg.JWT.JWTAccessToken.PublicKey)
		if err != nil {
			utresponse.GeneralUnauthorized(
				ctx,
				err,
			)
			ctx.Abort()
			return
		}

		if !slices.Contains(allowedlevel, tparsed.User.Role) || (consttypes.TimeNow().Unix() >= tparsed.Expires.Unix()) {
			utresponse.GeneralUnauthorized(
				ctx,
				consttypes.ErrUnauthorized,
			)
			ctx.Abort()
			return
		}

		if !utrequest.CheckWhitelistUrl(ctx.Request.URL.Path) {
			if tparsed.User.ConfirmedAt.IsZero() && !strings.Contains(ctx.Request.URL.Path, "verify") {
				utresponse.GeneralUnauthorized(
					ctx,
					consttypes.ErrAccountIsNotVerified,
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
