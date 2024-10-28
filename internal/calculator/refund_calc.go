package calculator

import (
	"github.com/cockroachdb/apd/v3"
	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/sportsbook"
)

type Selection interface {
	GetStatus() sportsbook.OddStatus
	GetValue() *apd.Decimal
}

type FormExpressesFunc func([]Selection, []int) [][]Selection

type RefundCalc struct {
	formExpresses FormExpressesFunc
	log           *zap.Logger
}

func NewRefundCalc(log *zap.Logger, formExpresses FormExpressesFunc) *RefundCalc {
	return &RefundCalc{log: log, formExpresses: formExpresses}
}

func (c *RefundCalc) Calc(
	betType callback.BetType,
	betSize []int,
	refundBase *apd.Decimal,
	selections []Selection,
) (*apd.Decimal, error) {
	refund, err := c.CalcRefund(betType, betSize, refundBase, selections)
	if err != nil {
		return nil, err
	}

	ctx := c.newApdCtx()
	_, err = ctx.Quantize(refund, refund, -6)

	return refund, err
}

func (c *RefundCalc) CalcRefund(
	betType callback.BetType,
	betSize []int,
	refundBase *apd.Decimal,
	selections []Selection,
) (*apd.Decimal, error) {
	if betType != callback.SystemBetType {
		return c.calcExpress(refundBase, selections)
	}

	return c.calcSystem(refundBase, betSize, selections)
}

func (c *RefundCalc) calcSystem(
	refundBase *apd.Decimal,
	systemSizes []int,
	selections []Selection,
) (*apd.Decimal, error) {
	refund := apd.New(0, 0)

	if len(selections) == 0 {
		return refund, nil
	}

	ctx := c.newApdCtx()

	expresses := c.formExpresses(selections, systemSizes)

	expressStakes, err := c.expressStakes(refundBase, len(expresses))
	if err != nil {
		return nil, err
	}

	for index, express := range expresses {
		expressRefund, err := c.calcExpress(expressStakes[index], express)
		if err != nil {
			return nil, err
		}

		_, err = ctx.Add(refund, refund, expressRefund)
		if err != nil {
			return nil, err
		}
	}

	return refund, nil
}

/**
 * Why it is incorrect to calculate the expresses stakes
 * simply by dividing the total stake by the of expresses count?
 *
 * Let's say your place system - 3 odds, 2 - size, totalStake - 1,
 * if divide 1/3 - you`ll get 0.333333, 0.333333, 0.333333, but where did the other 0.000001?
 * You must place a reminder for one of the expresses  example:
 * 1 - 0.333333
 * 2 - 0.333333
 * 3 - 0.333334
 */
func (c *RefundCalc) expressStakes(refundBase *apd.Decimal, expressesCount int) ([]*apd.Decimal, error) {
	ctx := c.newApdCtx()

	expressStake := new(apd.Decimal)
	expressesCountDec := apd.New(int64(expressesCount), 0)

	_, err := ctx.Quo(expressStake, refundBase, expressesCountDec)
	if err != nil {
		return nil, err
	}

	expressStakeWithoutReminder := new(apd.Decimal)

	_, err = ctx.Mul(expressStakeWithoutReminder, expressStake, expressesCountDec)
	if err != nil {
		return nil, err
	}

	expressStakeReminder := new(apd.Decimal)

	_, err = ctx.Sub(expressStakeReminder, refundBase, expressStakeWithoutReminder)
	if err != nil {
		return nil, err
	}

	expressStakes := make([]*apd.Decimal, expressesCount)
	for i := 0; i < expressesCount; i++ {
		expressStakes[i] = new(apd.Decimal).Set(expressStake)
	}

	lastExpressStakeIndex := expressesCount - 1

	_, err = ctx.Add(expressStakes[lastExpressStakeIndex], expressStakes[lastExpressStakeIndex], expressStakeReminder)
	if err != nil {
		return nil, err
	}

	return expressStakes, nil
}

func (c *RefundCalc) apdFromString(value string) *apd.Decimal {
	res, _, err := apd.BaseContext.WithPrecision(100).NewFromString(value)
	if err != nil {
		c.log.Warn("cannot convert string to decimal", zap.String("value", value))
	}

	return res
}

func (c *RefundCalc) calcExpress(
	refundBase *apd.Decimal,
	selections []Selection,
) (*apd.Decimal, error) {
	ctx := c.newApdCtx()

	r := c.apdFromString("1")

	for _, s := range selections {
		switch s.GetStatus() {
		case sportsbook.OddStatusLoss:
			return c.apdFromString("0"), nil

		case sportsbook.OddStatusWin:
			_, err := ctx.Mul(r, r, s.GetValue())
			if err != nil {
				return nil, err
			}

		case sportsbook.OddStatusHalfWin:
			n := c.apdFromString("0")
			_, err := ctx.Mul(n, s.GetValue(), c.apdFromString("1"))

			if err != nil {
				return nil, err
			}
			_, err = ctx.Quo(n, n, c.apdFromString("2"))

			if err != nil {
				return nil, err
			}

			_, err = ctx.Mul(r, r, n)
			if err != nil {
				return nil, err
			}

		case sportsbook.OddStatusHalfLoss:
			_, err := ctx.Quo(r, r, c.apdFromString("2"))
			if err != nil {
				return nil, err
			}

		case sportsbook.OddStatusRefunded:
			_, err := ctx.Mul(r, r, c.apdFromString("1"))
			if err != nil {
				return nil, err
			}

		default:
			r = c.apdFromString("0")
		}
	}

	res := new(apd.Decimal)

	_, err := ctx.Mul(res, r, refundBase)

	return res, err
}

func (c *RefundCalc) newApdCtx() *apd.Context {
	ctx := apd.BaseContext.WithPrecision(100)
	ctx.Rounding = apd.RoundDown

	return ctx
}
