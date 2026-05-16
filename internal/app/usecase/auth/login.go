package auth

import (
	"context"

	"go-clean-api/internal/domain/errors"
	"go-clean-api/internal/domain/user"
)

// NewLoginUsecase はログインユースケースを依存関係とともに生成します。
func NewLoginUsecase(
	userRepo user.Repository,
	tokenGenerator TokenGenerator,
	passwordVerifier PasswordVerifier) LoginUsecase {
	return &loginUsecase{
		userRepo:  userRepo,
		generator: tokenGenerator,
		verifier:  passwordVerifier,
	}
}

// LoginInput はログイン処理の入力値を表します。
type LoginInput struct {
	Email         string `validate:"required"`
	PlainPassword string `validate:"required"`
}

// LoginOutput はログイン成功時の出力値を表します。
type LoginOutput struct {
	Token string
}

// loginUsecase -.
type loginUsecase struct {
	userRepo  user.Repository
	generator TokenGenerator
	verifier  PasswordVerifier
}

// Execute はメールアドレスとパスワードで認証し、トークンを発行します。
func (uc *loginUsecase) Execute(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	entity, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.TypeInvalidCredentials.Wrap(err, "invalid email or password")
	}

	if err := uc.verifier.Verify(entity.PasswordHash(), input.PlainPassword); err != nil {
		return nil, errors.TypeInvalidCredentials.Wrap(err, "invalid email or password")
	}

	token, err := uc.generator.Generate(entity.ID().String())
	if err != nil {
		return nil, errors.TypeUnexpected.Wrap(err, "failed to generate token")
	}

	return &LoginOutput{Token: token}, nil
}
