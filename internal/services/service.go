package services

type ServiceContext struct {
	HealthService HealthService
	UsersService  UsersService
	OrdersService OrdersService
}
