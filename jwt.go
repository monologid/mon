package mon

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type IJwt interface {
	// Encrypt returns encrypted json web token
	Encrypt(data jwt.MapClaims) (string, error)

	// Decrypt returns map-claim after successfully parses and validate json web token
	// based on secret and signing method
	Decrypt(token string) (jwt.MapClaims, error)
}

type Jwt struct {
	Secret        string
	SigningMethod string
}

func (j *Jwt) Encrypt(data jwt.MapClaims) (string, error) {
	data["exp"] = time.Now().Add(time.Hour * 24).Unix()
	secret := []byte(j.Secret)
	return jwt.NewWithClaims(jwt.GetSigningMethod(j.SigningMethod), data).SignedString(secret)
}

func (j *Jwt) Decrypt(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.GetSigningMethod(j.SigningMethod) {
			return nil, errors.New("invalid signing method")
		}

		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func NewJwt(secret string, signingMethod string) IJwt {
	return &Jwt{
		Secret:        secret,
		SigningMethod: signingMethod,
	}
}
