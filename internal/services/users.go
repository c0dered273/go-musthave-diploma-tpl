package services

import (
	"context"
	"errors"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/middleware"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/repositories"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/validators"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

var (
	ErrUserAlreadyExist  = models.NewErrConflict(nil, "USER_ERROR", "Username already exists")
	ErrWrongLoginPasswd  = models.NewErrUnauthorized(nil, "USER_ERROR", "Wrong login or password")
	ErrRequestValidation = models.NewErrBadRequest(nil, "SERVER_ERROR", "Request validation failed")
	ErrInternal          = models.NewErrInternal(nil, "SERVER_ERROR", "Internal error")
)

type UsersService interface {
	NewUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error)
	LoginUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error)
	GetWithdrawals(ctx context.Context) (string, error)
}

type UsersServiceImpl struct {
	validator *validator.Validate
	cfg       *configs.ServerConfig
	repo      repositories.UsersRepository
	logger    zerolog.Logger
}

func NewUsersService(
	logger zerolog.Logger,
	cfg *configs.ServerConfig,
	validator *validator.Validate,
	repo repositories.UsersRepository,
) UsersService {
	return &UsersServiceImpl{
		validator: validator,
		cfg:       cfg,
		repo:      repo,
		logger:    logger,
	}
}

func (us *UsersServiceImpl) NewUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error) {
	err := us.loginValidation(login)
	if err != nil {
		return models.AuthResponseDTO{}, ErrRequestValidation
	}

	err = us.repo.Save(ctx, login.ToUser())
	if err != nil {
		if errors.Is(err, store.ErrUserAlreadyExists) {
			us.logger.Error().Err(err).Send()
			return models.AuthResponseDTO{}, ErrUserAlreadyExist
		}
		us.logger.Error().Err(err).Send()
		return models.AuthResponseDTO{}, ErrInternal
	}

	return us.getAuthResponse(login)
}

func (us *UsersServiceImpl) LoginUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error) {
	err := us.loginValidation(login)
	if err != nil {
		return models.AuthResponseDTO{}, ErrRequestValidation
	}

	user, err := us.repo.FindUserByNameAndPasswd(ctx, login.Login, login.Password)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return models.AuthResponseDTO{}, ErrInternal
	}
	if user.Username == "" {
		us.logger.Error().Msgf("server: username or passwd invalid, login: %s", login.Login)
		return models.AuthResponseDTO{}, ErrWrongLoginPasswd
	}

	return us.getAuthResponse(login)
}

func (us *UsersServiceImpl) GetWithdrawals(ctx context.Context) (string, error) {
	// TODO("Implement")
	claim, _ := ctx.Value(middleware.ClaimCtxKey{}).(*models.AuthClaim)
	return claim.ID, nil
}

func (us *UsersServiceImpl) loginValidation(login *models.LoginRequestDTO) error {
	err := validators.ValidateStructWithLogger(login, us.logger, us.validator)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return err
	}
	return nil
}

func (us *UsersServiceImpl) getAuthResponse(login *models.LoginRequestDTO) (models.AuthResponseDTO, error) {
	tokenString, err := generateToken(login.ToUser(), us.cfg.ApiSecret)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return models.AuthResponseDTO{}, ErrInternal
	}

	return models.AuthResponseDTO{
		AccessToken: tokenString,
	}, nil
}

func generateToken(user *models.User, secret string) (string, error) {
	claim := models.AuthClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
