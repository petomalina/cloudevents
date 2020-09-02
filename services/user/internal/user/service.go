package user

import (
	"context"

	v1 "github.com/flowup/petermalina/apis/go-sdk/user/v1"
	"go.uber.org/zap"
)

type Service struct {
	log        *zap.Logger
	repository Repository
	v1.UnimplementedUserServiceServer
}

func NewService(log *zap.Logger, repository Repository) *Service {
	return &Service{
		log:        log,
		repository: repository,
	}
}

func (s *Service) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.User, error) {
	panic("implement me")
}

func (s *Service) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.User, error) {
	panic("implement me")
}

func (s *Service) ListUsers(ctx context.Context, request *v1.ListUsersRequest) (*v1.ListUserResponse, error) {
	panic("implement me")
}

func (s *Service) SetUserPassword(ctx context.Context, request *v1.SetUserRequest) (*v1.SetUserPasswordResponse, error) {
	panic("implement me")
}
