package main

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/nutreet/common"
	proto "github.com/nutreet/common/gen/user"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	Register(ctx context.Context, data *proto.RegisterRequest) (*auth.UserRecord, error)
	GetAutenticatedUser(ctx context.Context, token string) (string, error)
}

type userService struct {
	logger *log.Logger
}

func NewUserService() UserService {
	logger := *log.New()
	logger.WithFields(log.Fields{
		"service": "User",
	})

	return &userService{
		logger: &logger,
	}
}

func (s *userService) Register(ctx context.Context, data *proto.RegisterRequest) (*auth.UserRecord, error) {
	user, err := common.Firebase.Auth.GetUserByEmail(ctx, data.Email)

	if err != nil {
		s.logger.Errorf("no firebase auth user related: %e", err)
		return nil, fmt.Errorf("no firebase auth user related")
	}

	_, err = common.Firebase.Firestore.Doc(fmt.Sprintf("users/%s", user.UID)).Create(ctx, map[string]interface{}{
		"uid":   user.UID,
		"email": data.Email,
	})

	if err != nil {
		s.logger.Errorf("error creating user in firestore: %e", err)
		return nil, fmt.Errorf("error creating user in firestore")
	}

	return user, nil
}

func (s *userService) GetAutenticatedUser(ctx context.Context, token string) (string, error) {
	u, err := common.Firebase.Auth.VerifyIDToken(ctx, token)

	if err != nil {
		s.logger.Errorf("erorr verifying id token: %e", err)
		return "", err
	}

	return u.UID, nil
}
