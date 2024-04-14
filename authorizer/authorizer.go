package authorizer

import (
	"crypto"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/theleeeo/thor/models"
)

type Authorizer struct {
	privateKey    crypto.PrivateKey
	publicKey     crypto.PublicKey
	rawPublicKey  []byte
	validDuration time.Duration
	appUrl        string

	parser *jwt.Parser
}

func New(cfg *Config) (*Authorizer, error) {
	pub, err := jwt.ParseEdPublicKeyFromPEM(cfg.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	priv, err := jwt.ParseEdPrivateKeyFromPEM(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Authorizer{
		privateKey:    priv,
		publicKey:     pub,
		rawPublicKey:  cfg.PublicKey,
		validDuration: cfg.ValidDuration,
		appUrl:        cfg.AppUrl,
		parser:        jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}), jwt.WithExpirationRequired()),
	}, nil
}

func (a *Authorizer) PublicKey() []byte {
	return a.rawPublicKey
}

func (a *Authorizer) Decode(token string) (*Claims, error) {
	t, err := a.parser.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func (a *Authorizer) CreateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA,
		&Claims{
			Issuer:    a.appUrl,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(a.validDuration),
			Role:      user.Role,
		},
	)

	tokenString, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
