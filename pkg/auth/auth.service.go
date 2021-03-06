package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ferjmc/cms/entities"
	"golang.org/x/crypto/scrypt"
)

type Auth interface {
	Scrypt(password string) ([]byte, error)
	GenerateToken(username string) (string, error)
	VerifyAuthorization(auth string) (string, string, error)
	VerifyToken(tokenString string) (string, error)
}

func New() Auth {
	return &auth{}
}

const TokenExpirationDays = 60
const PasswordKeyLength = 64

var passwordSalt = []byte("KU2YVXA7BSNExJIvemcdz61eL86IJDCC")
var jwtSecret = []byte("C92cw5od80NCWIvu4NZ8AKp5NyTbnBmG") // TODO: Generate random secrets and store in DynamoDB

type auth struct{}

func (a *auth) Scrypt(password string) ([]byte, error) {
	// https://godoc.org/golang.org/x/crypto/scrypt
	passwordHash, err := scrypt.Key([]byte(password), passwordSalt, 32768, 8, 1, PasswordKeyLength)
	if err != nil {
		return nil, err
	}

	return passwordHash, nil
}

func (a *auth) GenerateToken(username string) (string, error) {
	now := time.Now().UTC()
	exp := now.AddDate(0, 0, TokenExpirationDays).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": exp,
	})

	return token.SignedString(jwtSecret)
}

func (a *auth) VerifyAuthorization(auth string) (string, string, error) {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Token" {
		return "", "", entities.NewInputError("Authorization", "invalid format")
	}

	token := parts[1]
	username, err := a.VerifyToken(token)
	return username, token, err
}

func (a *auth) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, validateToken)

	if err != nil {
		return "", err
	}

	if token == nil || !token.Valid {
		return "", entities.NewInputError("Authorization", "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", entities.NewInputError("Authorization", "invalid claims")
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return "", entities.NewInputError("Authorization", "token expired")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return "", entities.NewInputError("Authorization", "sub missing")
	}

	return username, nil
}

func validateToken(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return jwtSecret, nil
}
