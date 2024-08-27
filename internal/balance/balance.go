package balance

import (
	"sync"

	"github.com/cockroachdb/apd/v3"
	"go.uber.org/zap"
)

type Balance struct {
	Available *apd.Decimal `json:"available"`
	Hold      *apd.Decimal `json:"hold"`
}

type Service struct {
	mu     sync.RWMutex
	active *apd.Decimal
	hold   *apd.Decimal
	log    *zap.Logger
}

func NewService(log *zap.Logger) *Service {
	return &Service{
		active: apd.New(0, 0),
		hold:   apd.New(0, 0),
		log:    log,
	}
}

func (s *Service) Deposit(amount *apd.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.newApdCtx().Add(s.active, s.active, amount)
	if err != nil {
		s.log.Error("failed to deposit", zap.Any("amount", amount), zap.Error(err))
	}
}

func (s *Service) DepositHold(amount *apd.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.newApdCtx().Add(s.hold, s.hold, amount)
	if err != nil {
		s.log.Error("failed to deposit", zap.Any("amount", amount), zap.Error(err))
	}
}

func (s *Service) DepositFloat(amount float64) error {
	decimal, err := apd.New(0, 0).SetFloat64(amount)
	if err != nil {
		return err
	}

	s.Deposit(decimal)

	return nil
}

func (s *Service) Withdraw(amount *apd.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.newApdCtx().Sub(s.active, s.active, amount)
	if err != nil {
		s.log.Error("failed to withdraw", zap.Any("amount", amount), zap.Error(err))
	}
}

func (s *Service) WithdrawHold(amount *apd.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.newApdCtx().Sub(s.hold, s.hold, amount)
	if err != nil {
		s.log.Error("failed to withdraw hold", zap.Any("amount", amount), zap.Error(err))
	}
}

func (s *Service) Hold(amount *apd.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.newApdCtx()

	_, err := ctx.Sub(s.active, s.active, amount)
	if err != nil {
		s.log.Error("failed to hold", zap.Any("amount", amount), zap.Error(err))

		return
	}

	_, err = ctx.Add(s.hold, s.hold, amount)
	if err != nil {
		s.log.Error("failed to hold", zap.Any("amount", amount), zap.Error(err))

		return
	}
}

func (s *Service) UnHold(amount *apd.Decimal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.newApdCtx()

	_, err := ctx.Add(s.active, s.active, amount)
	if err != nil {
		s.log.Error("failed to hold", zap.Any("amount", amount), zap.Error(err))

		return
	}

	_, err = ctx.Sub(s.hold, s.hold, amount)
	if err != nil {
		s.log.Error("failed to hold", zap.Any("amount", amount), zap.Error(err))

		return
	}
}

func (s *Service) State() Balance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return Balance{
		Available: s.active,
		Hold:      s.hold,
	}
}

func (s *Service) newApdCtx() *apd.Context {
	ctx := apd.BaseContext.WithPrecision(100)
	ctx.Rounding = apd.RoundDown

	return ctx
}
