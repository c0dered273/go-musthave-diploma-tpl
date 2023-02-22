package models

import (
	"time"

	"github.com/shopspring/decimal"
)

//go:generate easyjson withdrawal.go

type Withdrawal struct {
	OrderID     string
	Username    string
	Amount      *decimal.Decimal
	ProcessedAt time.Time
}

func (w Withdrawal) toWithdrawalDTO() WithdrawalDTO {
	return WithdrawalDTO{
		OrderID:     w.OrderID,
		Amount:      w.Amount.InexactFloat64(),
		ProcessedAt: w.ProcessedAt,
	}
}

type Withdrawals []Withdrawal

func ToWithdrawalsDTO(w Withdrawals) WithdrawalsDTO {
	l := len(w)
	withdrawalsDTO := make([]WithdrawalDTO, l)
	for i := range w {
		withdrawalsDTO[i] = w[i].toWithdrawalDTO()
	}

	return withdrawalsDTO
}

//easyjson:json
type WithdrawalDTO struct {
	OrderID     string    `json:"order"`
	Amount      float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

//easyjson:json
type WithdrawalsDTO []WithdrawalDTO

//easyjson:json
type WithdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}
