package auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims for access tokens
type AccessClaims struct {
	ID          string   `json:"id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"`
	jwt.RegisteredClaims
}

// RefreshClaims for refresh tokens
type RefreshClaims struct {
	ID        string `json:"id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenPair holds both access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(id string, roles []string, permissions []string, accessTTL, refreshTTL time.Duration) (TokenPair, error) {
	accessToken, err := GenerateAccessToken(id, roles, permissions, accessTTL)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := GenerateRefreshToken(id, refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateAccessToken creates an access token
func GenerateAccessToken(id string, roles []string, permissions []string, ttl time.Duration) (string, error) {
	claims := &AccessClaims{
		ID:          id,
		Permissions: permissions,
		Roles:       roles,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return signToken(claims, "JWT_SECRET")
}

// GenerateRefreshToken creates a refresh token
func GenerateRefreshToken(id string, ttl time.Duration) (string, error) {
	claims := &RefreshClaims{
		ID:        id,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return signToken(claims, "REFRESH_SECRET")
}

// signToken signs a JWT
func signToken(claims jwt.Claims, secretEnv string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv(secretEnv)
	if secret == "" {
		return "", fmt.Errorf("%s not set", secretEnv)
	}
	return token.SignedString([]byte(secret))
}

// VerifyToken verifies a token
func VerifyToken(tokenStr, tokenType string) (interface{}, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token missing")
	}

	secretEnv := "JWT_SECRET"
	if tokenType == "refresh" {
		secretEnv = "REFRESH_SECRET"
	}
	secret := os.Getenv(secretEnv)
	if secret == "" {
		return nil, fmt.Errorf("%s not set", secretEnv)
	}

	var claims jwt.Claims
	if tokenType == "access" {
		claims = &AccessClaims{}
	} else {
		claims = &RefreshClaims{}
	}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithStrictDecoding())

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func ExtractClaimsFromContext(ctx context.Context) (*AccessClaims, error) {
	// First, try getting from context
	if claims, ok := ctx.Value("claims").(*AccessClaims); ok {
		return claims, nil
	}

	// Fallback: extract from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := VerifyToken(token, "access")
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	accessClaims, ok := claims.(*AccessClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return accessClaims, nil
}
