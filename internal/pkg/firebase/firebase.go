package firebase

import (
	"context"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type Firebase struct {
	*auth.Client
}

func NewFirebase(firebaseCredentials []byte) func() (*Firebase, error) {

	return func() (*Firebase, error) {
		opt := option.WithCredentialsJSON(firebaseCredentials)

		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			return nil, err
		}

		auth, err := app.Auth(context.Background())
		if err != nil {
			return nil, err
		}

		return &Firebase{auth}, nil
	}
}

func (f *Firebase) ParseIDToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", domain.ErrParseIDToken
	}
	token := strings.TrimPrefix(
		tokenString, "Bearer ",
	)
	return token, nil
}
