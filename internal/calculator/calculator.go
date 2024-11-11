package calculator

import (
	"github.com/cockroachdb/apd/v3"
	"go.uber.org/zap"

	"github.com/databet-cloud/callback-test-tool/internal/callback"
)

type Calculator struct {
	refundCalculator *RefundCalc
	log              *zap.Logger
}

func NewCalculator(refundCalculator *RefundCalc, log *zap.Logger) *Calculator {
	return &Calculator{
		refundCalculator: refundCalculator,
		log:              log,
	}
}

func (c *Calculator) Settle(
	betType callback.BetType,
	sizes []int,
	betStake *apd.Decimal,
	odds []*callback.Odd,
) (*apd.Decimal, callback.SettleType, error) {
	selections := make([]Selection, len(odds))
	for i, odd := range odds {
		selections[i] = odd
	}

	amount, err := c.refundCalculator.Calc(betType, sizes, betStake, selections)
	if err != nil {
		return apd.New(0, 0), 0, err
	}

	switch amount.Cmp(betStake) {
	case 1:
		return amount, callback.WinSettleType, nil
	case 0:
		return amount, callback.RefundSettleType, nil
	default:
		return amount, callback.LossSettleType, nil
	}
}
