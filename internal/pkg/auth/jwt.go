package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AccessClaims for access tokens
type AccessClaims struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
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
func GenerateTokenPair(id, role string, accessTTL, refreshTTL time.Duration) (TokenPair, error) {
	accessToken, err := GenerateAccessToken(id, role, accessTTL)
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
func GenerateAccessToken(id, role string, ttl time.Duration) (string, error) {
	claims := &AccessClaims{
		ID:        id,
		Role:      role,
		TokenType: "access",
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

// JWTVerifyMiddleware for Gin
func JWTVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := VerifyToken(parts[1], "access")
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// JWTVerifyInterceptor for gRPC
func JWTVerifyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata missing")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header required")
	}

	parts := strings.Split(authHeaders[0], " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header")
	}

	claims, err := VerifyToken(parts[1], "access")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	newCtx := context.WithValue(ctx, "claims", claims)
	return handler(newCtx, req)
}

// AdminOnlyMiddleware for Gin
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(401, gin.H{"error": "Claims not found"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*AccessClaims)
		if !ok || userClaims.Role != "admin" {
			c.JSON(403, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
