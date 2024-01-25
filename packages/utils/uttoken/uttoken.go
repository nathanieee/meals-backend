package uttoken

import (
	"encoding/base64"
	"errors"
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/packages/utils/utlogger"
	"time"

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
		Authorized bool                    `json:"authorized"`
		User       *responses.UserResponse `json:"user"`
		Expire     int64                   `json:"expire"`
		TokenUUID  uuid.UUID               `json:"token_uuid"`
	}

	Token struct {
		Token     *string                 `json:"token"`
		TokenUUID uuid.UUID               `json:"token_uuid"`
		Expires   *time.Time              `json:"expires"`
		User      *responses.UserResponse `json:"user"`
	}
)

func GenerateToken(
	ures *responses.UserResponse,
	lifespan int,
	timeunit,
	privateKey string,
) (*Token, error) {
	expTime := time.Now().Add(getDuration(lifespan, timeunit))
	tokenUUID := uuid.New()

	claims := TokenClaims{
		Authorized: true,
		User:       ures,
		Expire:     expTime.Unix(),
		TokenUUID:  tokenUUID,
	}

	token := &Token{
		TokenUUID: tokenUUID,
		Expires:   &expTime,
		User:      ures,
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		utlogger.LogError(err)
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		utlogger.LogError(err)
		return nil, fmt.Errorf("could not parse token private key: %w", err)
	}

	newToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		utlogger.LogError(err)
		return nil, fmt.Errorf("could not create token: %w", err)
	}

	token.Token = &newToken

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
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		utlogger.LogError(err)
		return nil, fmt.Errorf("could not decode token public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		utlogger.LogError(err)
		return nil, fmt.Errorf("could not parse token public key: %w", err)
	}

	claims := &TokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &Token{
		TokenUUID: claims.TokenUUID,
		User:      claims.User,
	}, nil
}
