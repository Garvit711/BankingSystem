package account

import "time"

type Account struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	Currency      string    `json:"currency"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}
