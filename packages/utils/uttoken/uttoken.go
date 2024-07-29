package uttoken

import (
	"encoding/base64"
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type (
	TokenHeader struct {
		AccessToken         string
		AccessTokenExpires  time.Time
		RefreshToken        string
		RefreshTokenExpires time.Time
	}

	TokenClaims struct {
		jwt.StandardClaims
		Authorized bool            `json:"authorized"`
		User       *responses.User `json:"user"`
		Expire     int64           `json:"expire"`
		TokenUUID  uuid.UUID       `json:"token_uuid"`
	}

	Token struct {
		Token     *string         `json:"token"`
		TokenUUID uuid.UUID       `json:"token_uuid"`
		Expires   *time.Time      `json:"expires"`
		User      *responses.User `json:"user"`
	}
)

func NewTokenHeader(access Token, refresh Token) *TokenHeader {
	return &TokenHeader{
		AccessToken:         *access.Token,
		AccessTokenExpires:  *access.Expires,
		RefreshToken:        *refresh.Token,
		RefreshTokenExpires: *refresh.Expires,
	}
}

func (token *TokenHeader) ToAuthResponse(user responses.User) *responses.Auth {
	return &responses.Auth{
		ID:                 user.ID,
		Email:              user.Email,
		Role:               user.Role,
		ConfirmationSentAt: user.ConfirmationSentAt,
		ConfirmedAt:        user.ConfirmedAt,
		CreatedAt:          *user.CreatedAt,
		UpdatedAt:          *user.UpdatedAt,
		Token:              token.AccessToken,
		Expires:            token.AccessTokenExpires,
	}
}

func GenerateToken(
	ures *responses.User,
	lifespan int,
	timeunit,
	privatekey string,
) (*Token, error) {
	exptime := consttypes.TimeNow().Add(getDuration(lifespan, timeunit))
	tuuid := uuid.New()

	claims := TokenClaims{
		Authorized: true,
		User:       ures,
		Expire:     exptime.Unix(),
		TokenUUID:  tuuid,
	}

	token := &Token{
		TokenUUID: tuuid,
		Expires:   &exptime,
		User:      ures,
	}

	decodedprivatekey, err := base64.StdEncoding.DecodeString(privatekey)
	if err != nil {
		utlogger.Error(err)
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedprivatekey)
	if err != nil {
		utlogger.Error(err)
		return nil, fmt.Errorf("could not parse token private key: %w", err)
	}

	tclaim, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		utlogger.Error(err)
		return nil, fmt.Errorf("could not create token: %w", err)
	}

	token.Token = &tclaim

	return token, nil
}

func getDuration(lifespan int, timeunit string) time.Duration {
	switch timeunit {
	case "minutes":
		return time.Minute * time.Duration(lifespan)
	case "seconds":
		return time.Second * time.Duration(lifespan)
	default:
		return time.Hour * time.Duration(lifespan)
	}
}

func ParseToken(token string, pubkey string) (*Token, error) {
	decodedpubkey, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		utlogger.Error(err)
		return nil, fmt.Errorf("could not decode token public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedpubkey)
	if err != nil {
		utlogger.Error(err)
		return nil, fmt.Errorf("could not parse token public key: %w", err)
	}

	claims := &TokenClaims{}
	tparsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	claims, ok := tparsed.Claims.(*TokenClaims)
	if !ok || !tparsed.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	expires := time.Unix(claims.Expire, 0)

	return &Token{
		Token:     &token,
		TokenUUID: claims.TokenUUID,
		Expires:   &expires,
		User:      claims.User,
	}, nil
}

func GetToken(ctx *gin.Context) (string, string, error) {
	taccess := ctx.GetHeader(consttypes.T_ACCESS)
	if taccess == "" {
		taccess = ctx.Request.Header.Get("X-Authorization")
	}

	if taccess == "" {
		return "", "", fmt.Errorf("access token not found")
	}

	trefresh := ctx.GetHeader(consttypes.T_REFRESH)
	if trefresh == "" {
		return "", "", fmt.Errorf("refresh token not found")
	}

	return taccess, trefresh, nil
}

func DeleteToken(ctx *gin.Context) {
	ctx.Request.Header.Del(consttypes.T_REFRESH)
	ctx.Request.Header.Del(consttypes.T_ACCESS)
}

func GetUser(ctx *gin.Context) (*responses.User, error) {
	userctx, exists := ctx.Get("user")
	if !exists {
		return nil, consttypes.ErrUserIDNotFound
	}

	userres, ok := userctx.(responses.User)
	if !ok {
		return nil, consttypes.ErrUserNotFound
	}

	return &userres, nil
}
