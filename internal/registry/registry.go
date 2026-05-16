package registry

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"go-clean-api/config"
	userquery "go-clean-api/internal/app/query/user"
	querymemory "go-clean-api/internal/app/query/user/provider/memory"
	authuc "go-clean-api/internal/app/usecase/auth"
	"go-clean-api/internal/infra/mariadb"
	"go-clean-api/pkg/bcrypt"
	"go-clean-api/pkg/jwt"
	"go-clean-api/pkg/log"
	"go-clean-api/pkg/sqldb"
	"go-clean-api/pkg/validate"
)

// TODO: workflow の場合は別のコンテナを作成する
func New(cfg *config.Server) (*Registry, error) {
	// setup Logger
	logger, err := log.New(os.Stdout,
		log.WithLevel(cfg.Log.Level),
		log.WithGlobalAttrs(
			"app", cfg.APP.Name,
			"version", cfg.APP.Version,
		),
		log.WithPretty(cfg.Log.Pretty),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	// setup mariadb client
	dbClient, err := sqldb.NewClient(cfg.DB.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to setup mariadb client: %w", err)
	}

	// setup validator
	validator, err := validate.New()
	if err != nil {
		return nil, err
	}

	// setup jwt
	jwtMng, err := jwt.New(cfg.JWT.Secret,
		jwt.WithTokenExpiry(cfg.JWT.TokenExpiry),
	)
	if err != nil {
		return nil, err
	}

	// setup bcrypt hasher & verifier
	bcryptVerifier := bcrypt.NewVerifier()
	bcryptHasher := bcrypt.NewHasher()

	// setup repositories
	userRepo := mariadb.NewUserRepository(dbClient)

	// setup dependencies for wire
	deps := &deps{
		dbClient:  dbClient,
		validator: validator,
		logger:    logger,
		nowFunc:   time.Now,
	}

	// setup query services
	// ユーザ検索はトランザクションを必要としないので applyStandard を使用
	userQuery := applyStandard("user-search-query",
		querymemory.NewUserQuery(),
		deps,
	)

	// setup usecases
	// ログイン認証はトランザクションを必要としないので applyStandard を使用
	authLoginUsecase := applyStandard("auth-login",
		authuc.NewLoginUsecase(userRepo, jwtMng, bcryptVerifier),
		deps,
	)
	// ユーザ登録はトランザクションを必要とするので applyStandardWithRequiredTx を使用
	authRegisterUsecase := applyStandardWithRequiredTx("auth-register",
		authuc.NewRegisterUsecase(userRepo, bcryptHasher),
		deps,
	)

	// setup registry
	registry := Registry{Logger: logger, Validator: validator}
	registry.UsecaseSet.AuthLogin = authLoginUsecase
	registry.UsecaseSet.AuthRegister = authRegisterUsecase
	registry.QuerySet.UserSearch = userQuery

	return &registry, nil
}

type Registry struct {
	Logger     *slog.Logger
	Validator  *validate.Validator
	UsecaseSet struct {
		AuthLogin    authuc.LoginUsecase
		AuthRegister authuc.RegisterUsecase
	}
	QuerySet struct {
		UserSearch userquery.QueryService
	}
}
