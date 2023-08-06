package usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/auth"
)

type UserUsecase struct {
	auth     *auth.Firebase
	userRepo domain.UserRepository
}

func NewUserUsecase(auth *auth.Firebase, userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecase{
		auth:     auth,
		userRepo: userRepo,
	}
}

// GenerateIDToken implements domain.UserUsecase.
func (u *UserUsecase) GenerateIDToken(ctx context.Context, uid string) (string, error) {
	return u.auth.CustomToken(ctx, uid)
}

func (u *UserUsecase) VerifyIDToken(ctx context.Context, authorization string) (*domain.User, error) {

	token, err := u.auth.ParseIDToken(authorization)
	if err != nil {
		return nil, err
	}
	log.Print("token:", token)
	authToken, err := u.auth.VerifyIDToken(ctx, token)
	if err != nil {

		if e := log.Debug(); e.Enabled() {
			e.Msg(err.Error())

		}
		return nil, err
	}

	return &domain.User{
		Email:          authToken.Firebase.Identities["email"].([]any)[0].(string),
		SignInProvider: authToken.Firebase.SignInProvider,
		Uid:            authToken.UID,
	}, nil
}
func (u *UserUsecase) Sync(ctx context.Context, user domain.User) (*domain.User, error) {
	_user, err := u.userRepo.FetchOrCreate(ctx, user)
	if len(_user.Locations) == 0 {

	}
	return _user, err
}
func (u *UserUsecase) GetAuthUserByUid(ctx context.Context, uid string) (*domain.User, error) {
	_uid, err := u.auth.ParseIDToken(uid)
	if err != nil {
		return nil, err
	}
	userRecord, err := u.auth.GetUser(ctx, _uid)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Uid:            userRecord.UID,
		Email:          userRecord.Email,
		SignInProvider: userRecord.ProviderID,
	}, nil
}
