package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("User Not Found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, userID int64, number string) (Account, error) {
	var a Account
	err := r.db.QueryRow(ctx,
		`INSERT INTO accounts (user_id, account_number)
		 VALUES ($1, $2)
		 RETURNING id, user_id, account_number, currency, balance, created_at`,
		userID, number,
	).Scan(&a.ID, &a.UserID, &a.AccountNumber, &a.Currency, &a.Balance, &a.CreatedAt)
	if err != nil {
		return Account{}, err
	}
	return a, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID int64) ([]Account, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, account_number, currency, balance, created_at
		 FROM accounts WHERE user_id = $1 ORDER BY id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.UserID, &a.AccountNumber, &a.Currency, &a.Balance, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *Repository) GetByIDAndUser(ctx context.Context, id, userID int64) (Account, error) {
	var a Account
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, account_number, currency, balance, created_at
		 FROM accounts WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&a.ID, &a.UserID, &a.AccountNumber, &a.Currency, &a.Balance, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, ErrNotFound
		}
		return Account{}, err
	}
	return a, nil
}
