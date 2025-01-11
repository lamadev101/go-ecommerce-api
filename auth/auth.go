package auth

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/lamadev101/ecommerce-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtWrapper struct {
	SecretKey      string
	Issuer         string
	ExpirationTime int64
}

type JwtClaim struct {
	UserId   primitive.ObjectID
	Email    string
	UserType string
	jwt.StandardClaims
}

// generate a token
func (j *JwtWrapper) GenerateToken(id primitive.ObjectID, email, userType string) (token string, err error) {
	claims := &JwtClaim{
		UserId:   id,
		UserType: userType,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: time.Now().Add(time.Hour * time.Duration(j.ExpirationTime))},
			Issuer:    j.Issuer,
		},
	}

	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = token1.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// validate token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("could not parse claims")
		return
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token is expired")
		return
	}
	return
}

// Check authorization
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondWithUnauthorizedError(c, "Authorization header is missing")
			return
		}

		// Validate Bearer token format
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			utils.RespondWithUnauthorizedError(c, "Invalid token format")
			return
		}

		// Trim whitespace from the token
		token := strings.TrimSpace(tokenParts[1])
		if token == "" {
			utils.RespondWithUnauthorizedError(c, "Token is empty")
			return
		}

		// Create JWT wrapper
		jwtWrapper := JwtWrapper{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			Issuer:    os.Getenv("JWT_ISSUER"),
		}

		// Validate the token and extract claims
		claims, err := jwtWrapper.ValidateToken(token)
		if err != nil {
			utils.RespondWithUnauthorizedError(c, "Invalid or expired token")
			return
		}

		// Set user details in the context for further use
		c.Set("user_id", claims.ID)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}
