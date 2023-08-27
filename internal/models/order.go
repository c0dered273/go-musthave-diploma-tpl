package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type OrderStatus int

const (
	NEW OrderStatus = 1 + iota
	PROCESSING
	INVALID
	PROCESSED
)

var (
	statusString = []string{
		"",
		"NEW",
		"PROCESSING",
		"INVALID",
		"PROCESSED",
	}
)

func (s OrderStatus) String() string {
	return statusString[s]
}

func ParseStatus(name string) (OrderStatus, error) {
	switch strings.ToUpper(name) {
	case "NEW":
		return NEW, nil
	case "PROCESSING":
		return PROCESSING, nil
	case "INVALID":
		return INVALID, nil
	case "PROCESSED":
		return PROCESSED, nil
	default:
		return 0, errors.New("failed to parse OrderStatus")
	}
}

//go:generate easyjson order.go

type Order struct {
	ID         uint64
	Status     OrderStatus
	Username   string
	Amount     *decimal.Decimal
	UploadedAt time.Time
}

func (o Order) toOrderDTO() OrderDTO {
	if o.Amount == nil {
		o.Amount = &decimal.Zero
	}

	return OrderDTO{
		ID:         strconv.FormatUint(o.ID, 10),
		Status:     o.Status.String(),
		Amount:     o.Amount.InexactFloat64(),
		UploadedAt: o.UploadedAt,
	}
}

type Orders []Order

func ToOrdersDTO(o Orders) OrdersDTO {
	l := len(o)
	ordersDTO := make([]OrderDTO, l)
	for i := range o {
		ordersDTO[i] = o[i].toOrderDTO()
	}

	return ordersDTO
}

//easyjson:json
type OrderDTO struct {
	ID         string    `json:"number"`
	Status     string    `json:"status"`
	Amount     float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

//easyjson:json
type OrdersDTO []OrderDTO

//easyjson:json
type AccrualOrderDTO struct {
	ID      string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (au AccrualOrderDTO) ToOrder() (Order, error) {
	orderID, err := strconv.ParseUint(au.ID, 10, 64)
	if err != nil {
		return Order{}, err
	}

	status, err := ParseStatus(au.Status)
	if err != nil {
		return Order{}, err
	}

	amount := decimal.NewFromFloat(au.Accrual)

	return Order{
		ID:         orderID,
		Status:     status,
		Amount:     &amount,
		UploadedAt: time.Now(),
	}, nil
}
