package account

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID int64) (Account, error) {
	return s.repo.Create(ctx, userID, generateNumber())
}

func (s *Service) List(ctx context.Context, userID int64) ([]Account, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) Get(ctx context.Context, id, userID int64) (Account, error) {
	return s.repo.GetByIDAndUser(ctx, id, userID)
}

func (s *Service) Deposit(ctx context.Context, accountID, userID, amount int64) (Account, error) {
	return s.repo.Deposit(ctx, accountID, userID, amount)
}
