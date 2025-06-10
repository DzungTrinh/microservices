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

type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(id uint, email, role string) (string, error) {
	claims := &Claims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	return token.SignedString([]byte(jwtSecret))
}

func VerifyJWT(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token missing")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithStrictDecoding())

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

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

		claims, err := VerifyJWT(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

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

	claims, err := VerifyJWT(parts[1])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	newCtx := context.WithValue(ctx, "claims", claims)
	return handler(newCtx, req)
}
