package user

import (
	"context"

	"cloud.google.com/go/firestore"
	v1 "github.com/flowup/petermalina/apis/go-sdk/user/v1"
)

type User struct {
	ID   string
	Name string
}

func (u User) ToAPIv1() *v1.User {
	return &v1.User{
		Id:   u.ID,
		Name: u.Name,
	}
}

type Repository interface {
	Get(ctx context.Context, id string) (User, error)
	Create(context.Context, User) (User, error)
	Update(context.Context, User) (User, error)
	Delete(ctx context.Context, id string) error
}

func NewRepository(fs *firestore.CollectionRef) *FSRepository {
	return &FSRepository{fs: fs}
}

type FSRepository struct {
	fs *firestore.CollectionRef
}

func (r *FSRepository) Get(ctx context.Context, id string) (User, error) {
	panic("implement me")
}

func (r *FSRepository) Create(ctx context.Context, user User) (User, error) {
	panic("implement me")
}

func (r *FSRepository) Update(ctx context.Context, user User) (User, error) {
	panic("implement me")
}

func (r *FSRepository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
