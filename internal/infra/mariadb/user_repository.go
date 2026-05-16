package mariadb

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"

	"go-clean-api/internal/domain/errors"
	"go-clean-api/internal/domain/user"
	"go-clean-api/internal/infra/mariadb/models"
	"go-clean-api/pkg/sqldb"
)

type userRepository struct {
	client *sqldb.Client
}

// NewUserRepository -.
func NewUserRepository(client *sqldb.Client) user.Repository {
	return &userRepository{
		client: client,
	}
}

// FindByEmail -.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	tx := r.client.CurrentTx(ctx)

	// find user by email
	entity, err := models.Users(
		models.UserWhere.Email.EQ(email),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.TypeNotFound.Wrap(err, "user not found")
		}
		return nil, errors.TypeUnexpected.Wrap(err, "failed to find user by email")
	}

	return user.Reconstruct(
		user.NewID(entity.ID),
		entity.Name,
		entity.Email,
		[]byte(entity.PasswordHash),
	), nil
}

// Save -.
func (r *userRepository) Save(ctx context.Context, entity *user.User) (user.ID, error) {
	tx := r.client.CurrentTx(ctx)

	model := &models.User{
		ID:           entity.ID().Int(),
		Name:         entity.Name(),
		Email:        entity.Email(),
		PasswordHash: string(entity.PasswordHash()),
	}

	if model.ID == 0 {
		if err := model.Insert(ctx, tx, boil.Infer()); err != nil {
			// TODO: 実際はベンダエラーコードをチェックし、エラーを返却
			return 0, errors.TypeUnexpected.Wrap(err, "failed to insert user")
		}
	} else {
		if _, err := model.Update(ctx, tx, boil.Infer()); err != nil {
			// TODO: 実際はベンダエラーコードをチェックし、エラーを返却
			return 0, errors.TypeUnexpected.Wrap(err, "failed to update user")
		}
	}
	return user.NewID(model.ID), nil
}
