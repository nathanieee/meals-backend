package uttoken

import (
	"encoding/base64"
	"errors"
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
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
		Token:     token.AccessToken,
		Expires:   token.AccessTokenExpires,
	}
}

func GenerateToken(
	ures *responses.User,
	lifespan int,
	timeunit,
	privatekey string,
) (*Token, error) {
	exptime := consttypes.DateNow.Add(getDuration(lifespan, timeunit))
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

func ParseToken(token string, publicKey string) (*Token, error) {
	decodedpublickey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		utlogger.Error(err)
		return nil, fmt.Errorf("could not decode token public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedpublickey)
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
			return nil, errors.New("invalid token signature")
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

func GetUser(ctx *gin.Context) (*responses.User, error) {
	claims := ctx.MustGet("claims").(*TokenClaims)
	return claims.User, nil
}
