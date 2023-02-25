package services

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/clients"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/middleware"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/repositories"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/validators"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

var (
	ErrUserAlreadyExist = models.NewErrConflict(nil, "USER_ERROR", "Username already exists")
	ErrWrongLoginPasswd = models.NewErrUnauthorized(nil, "USER_ERROR", "Invalid login or password")
	ErrInvalidToken     = models.NewErrUnauthorized(nil, "AUTH_ERROR", "Access token invalid")
	ErrInvalidOrderID   = models.NewStatusError(
		nil,
		http.StatusUnprocessableEntity,
		"ORDER_ERROR",
		"Invalid order number format",
	)
	ErrOrderLoadedByAnotherUser = models.NewStatusError(
		nil,
		http.StatusConflict,
		"ORDER_ERROR",
		"Order loaded by another user",
	)
	ErrPaymentRequired   = models.NewErrPaymentRequired(nil, "USER_ERROR", "Balance not enough")
	ErrRequestValidation = models.NewErrBadRequest(nil, "SERVER_ERROR", "Request validation failed")
	ErrInternal          = models.NewErrInternal(nil, "SERVER_ERROR", "Internal error")
)

type UsersService interface {
	NewUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error)
	LoginUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error)
	CreateOrders(ctx context.Context, orderNumber string) error
	GetOrders(ctx context.Context) (models.OrdersDTO, error)
	GetWithdrawals(ctx context.Context) (models.WithdrawalsDTO, error)
	GetBalance(ctx context.Context) (models.UserBalanceDTO, error)
	WithdrawBalance(ctx context.Context, orderID string, amount decimal.Decimal) error
}

type UsersServiceImpl struct {
	validator      *validator.Validate
	cfg            *configs.ServerConfig
	userRepo       repositories.UserRepository
	orderRepo      repositories.OrderRepository
	withdrawalRepo repositories.WithdrawalRepository
	accrualClient  *resty.Client
	logger         zerolog.Logger
}

func NewUsersService(
	logger zerolog.Logger,
	cfg *configs.ServerConfig,
	validator *validator.Validate,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	withdrawalRepo repositories.WithdrawalRepository,
	accrualClient *resty.Client,
) UsersService {
	return &UsersServiceImpl{
		validator:      validator,
		cfg:            cfg,
		userRepo:       userRepo,
		orderRepo:      orderRepo,
		withdrawalRepo: withdrawalRepo,
		accrualClient:  accrualClient,
		logger:         logger,
	}
}

func (us *UsersServiceImpl) NewUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error) {
	err := us.loginValidation(login)
	if err != nil {
		return models.AuthResponseDTO{}, ErrRequestValidation
	}

	newUser := login.ToUser()
	err = us.userRepo.Save(ctx, newUser)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			us.logger.Error().Err(err).Send()
			return models.AuthResponseDTO{}, ErrUserAlreadyExist
		}
		us.logger.Error().Err(err).Send()
		return models.AuthResponseDTO{}, ErrInternal
	}

	return us.authResponse(newUser)
}

func (us *UsersServiceImpl) LoginUser(ctx context.Context, login *models.LoginRequestDTO) (models.AuthResponseDTO, error) {
	err := us.loginValidation(login)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return models.AuthResponseDTO{}, ErrRequestValidation
	}

	user, err := us.userRepo.FindByNameAndPasswd(ctx, login.Login, login.Password)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			us.logger.Error().Msgf("server: username or passwd invalid, login: %s", login.Login)
			return models.AuthResponseDTO{}, ErrWrongLoginPasswd
		}
		us.logger.Error().Err(err).Send()
		return models.AuthResponseDTO{}, ErrInternal
	}

	return us.authResponse(user)
}

func (us *UsersServiceImpl) CreateOrders(ctx context.Context, orderString string) error {
	claim, err := claimFromCtx(ctx)
	if err != nil {
		return err
	}

	orderID, err := strconv.ParseUint(orderString, 10, 64)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return ErrInvalidOrderID
	}

	err = orderIDValidate(orderID)
	if err != nil {
		return ErrInvalidOrderID
	}

	existsOrder, err := us.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			us.logger.Error().Err(err).Send()
			return ErrInternal
		}
	}

	if existsOrder != nil {
		if existsOrder.Username == claim.ID {
			return models.NewStatusCreated("Already processing")
		} else {
			return ErrOrderLoadedByAnotherUser
		}
	}

	newOrder := &models.Order{
		ID:         orderID,
		Status:     models.NEW,
		Username:   claim.ID,
		UploadedAt: time.Now(),
	}

	err = us.orderRepo.Save(ctx, newOrder)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return models.NewStatusCreated("Order already exists")
		}
		us.logger.Error().Err(err).Send()
		return ErrInternal
	}

	// TODO("Отрефакторить, вынести работу с базой в хранимку с транзакцией")

	go func() {
		get, err := us.accrualClient.R().Get(clients.AccrualURL + orderString)
		if err != nil {
			us.logger.Error().Err(err).Send()
			return
		}

		accrualOrderResponse := models.AccrualOrderDTO{}
		err = easyjson.Unmarshal(get.Body(), &accrualOrderResponse)
		if err != nil {
			us.logger.Error().Err(err).Send()
			return
		}

		order, err := accrualOrderResponse.ToOrder()
		if err != nil {
			us.logger.Error().Err(err).Send()
			return
		}

		err = us.orderRepo.UpdateByID(context.Background(), order.ID, order.Status, *order.Amount)
		if err != nil {
			us.logger.Error().Err(err).Send()
			return
		}

		err = us.userRepo.AccrueBalance(context.Background(), claim.ID, *order.Amount)
		if err != nil {
			us.logger.Error().Err(err).Send()
			return
		}

	}()

	return nil
}

func (us *UsersServiceImpl) GetOrders(ctx context.Context) (models.OrdersDTO, error) {
	claim, err := claimFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	orders, err := us.orderRepo.FindByUsername(ctx, claim.ID)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return nil, ErrInternal
	}

	return models.ToOrdersDTO(orders), nil
}

func (us *UsersServiceImpl) GetWithdrawals(ctx context.Context) (models.WithdrawalsDTO, error) {
	claim, err := claimFromCtx(ctx)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return nil, ErrInternal
	}

	withdrawals, err := us.withdrawalRepo.FindByUsername(ctx, claim.ID)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return nil, ErrInternal
	}

	return models.ToWithdrawalsDTO(withdrawals), nil
}

func (us *UsersServiceImpl) GetBalance(ctx context.Context) (models.UserBalanceDTO, error) {
	claim, err := claimFromCtx(ctx)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return models.UserBalanceDTO{}, ErrInternal
	}

	balance, err := us.userRepo.GetBalance(ctx, claim.ID)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return models.UserBalanceDTO{}, ErrInternal
	}

	allWithdrawals, err := us.withdrawalRepo.GetAllWithdrawalByUsername(ctx, claim.ID)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return models.UserBalanceDTO{}, ErrInternal
	}

	return models.UserBalanceDTO{
		Current:   balance.InexactFloat64(),
		Withdrawn: allWithdrawals.InexactFloat64(),
	}, nil
}

func (us *UsersServiceImpl) WithdrawBalance(ctx context.Context, orderID string, amount decimal.Decimal) error {
	claim, err := claimFromCtx(ctx)
	if err != nil {
		us.logger.Error().Err(err).Send()
		return ErrInternal
	}

	err = us.userRepo.Withdrawing(ctx, claim.ID, orderID, amount)
	if err != nil {
		if errors.Is(err, repositories.ErrBalanceNotEnough) {
			return ErrPaymentRequired
		}
		us.logger.Error().Err(err).Send()
		return ErrInternal
	}

	return nil
}

func (us *UsersServiceImpl) loginValidation(login *models.LoginRequestDTO) error {
	err := validators.ValidateStructWithLogger(login, us.logger, us.validator)
	if err != nil {
		return err
	}
	return nil
}

func (us *UsersServiceImpl) authResponse(user *models.User) (models.AuthResponseDTO, error) {
	tokenString, err := generateToken(user, us.cfg.APISecret)
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

func claimFromCtx(ctx context.Context) (*models.AuthClaim, error) {
	claim, ok := ctx.Value(middleware.ClaimCtxKey{}).(*models.AuthClaim)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claim, nil
}

func orderIDValidate(id uint64) error {
	check := (id%10 + luhnChecksum(id/10)) % 10
	if check != 0 {
		return errors.New("invalid order id")
	}

	return nil
}

func luhnChecksum(number uint64) uint64 {
	var luhn uint64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
