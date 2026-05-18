package auth

import (
	"context"

	"go-clean-api/internal/domain/errors"
	"go-clean-api/internal/domain/user"
)

// NewRegisterUsecase -.
func NewRegisterUsecase(
	userRepo user.Repository,
	hasher PasswordHasher) RegisterUsecase {
	return &registerUsecase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

// RegisterInput -.
type RegisterInput struct {
	Name          string `validate:"required"`
	Email         string `validate:"required"`
	PlainPassword string `validate:"required"`
}

// RegisterOutput -.
type RegisterOutput struct {
	ID user.ID
}

// registerUsecase -.
type registerUsecase struct {
	userRepo user.Repository
	hasher   PasswordHasher
}

// Execute -.
func (uc *registerUsecase) Execute(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	passwordHash, err := uc.hasher.Hash(input.PlainPassword)
	if err != nil {
		return nil, errors.TypeUnexpected.Wrap(err, "failed to hash password")
	}

	entity := user.New(input.Name, input.Email, passwordHash)
	newID, err := uc.userRepo.Save(ctx, entity)
	if err != nil {
		if errors.TypeAlreadyExists.Is(err) {
			return nil, err
		}
		return nil, errors.TypeUnexpected.Wrap(err, "failed to save user")
	}

	return &RegisterOutput{ID: newID}, nil
}
