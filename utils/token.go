package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/thuongnn/clst-mgt-api/models"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}

func DecodeOauth2Token[T any](tokenString string, out *T) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWT token")
	}

	payload, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode JWT payload: %v", err)
	}

	if err := json.Unmarshal(payload, out); err != nil {
		return fmt.Errorf("failed to parse JWT claims: %v", err)
	}

	return nil
}

func FetchWellKnownConfig(wellKnownConfigUrl string) (*models.WellKnownConfig, error) {
	resp, err := http.Get(wellKnownConfigUrl)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d", resp.StatusCode)
	}

	var config models.WellKnownConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("parse JSON error: %v", err)
	}

	return &config, nil
}

func ParseOAuth2Config(oauth2Info models.OAuth2Config) (*oauth2.Config, error) {
	wellKnownConfig, err := FetchWellKnownConfig(oauth2Info.WellKnownConfigURL)
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     oauth2Info.ClientID,
		ClientSecret: oauth2Info.ClientSecret,
		RedirectURL:  oauth2Info.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  wellKnownConfig.AuthURL,
			TokenURL: wellKnownConfig.TokenURL,
		},
		Scopes: oauth2Info.Scopes,
	}, nil
}
