package calculator

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/internal/calculator/former"
	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/sportsbook"
)

//nolint:govet,funlen //it's ok for test
func TestRefundCalc_Calc(t *testing.T) {
	var table = []struct {
		name       string
		betType    callback.BetType
		betSize    []int
		stake      *apd.Decimal
		selections []Selection
		refund     *apd.Decimal
	}{
		{
			"single_win",
			callback.SingleBetType,
			[]int{1},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{

					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("12.5"),
		},
		{
			"single_loss",
			callback.SingleBetType,
			[]int{1},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{

					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusLoss,
				},
			}),
			toDecimal("0"),
		},
		{
			"single_manually_refund",
			callback.SingleBetType,
			[]int{1},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{

					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusRefundedManually,
				},
			}),
			toDecimal("5"),
		},
		{
			"single_refund",
			callback.SingleBetType,
			[]int{1},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{

					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusRefunded,
				},
			}),
			toDecimal("5"),
		},
		{
			"express_win",
			callback.ExpressBetType,
			[]int{3},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("78.125000"),
		},
		{
			"express_win_3",
			callback.ExpressBetType,
			[]int{3},
			toDecimal("1.93"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("1.93"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1.93"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1.93"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1.93"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("26.778518"),
		},
		{
			"express_loss",
			callback.ExpressBetType,
			[]int{3},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusLoss,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("0"),
		},
		{
			"system_2_win",
			callback.SystemBetType,
			[]int{2},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("31.249999"),
		},
		{
			"system_7_6_5_win",
			callback.SystemBetType,
			[]int{5},
			toDecimal("7"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{OddRatio: toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("7"),
		},
		{
			"system_1_3_2_win",
			callback.SystemBetType,
			[]int{2},
			toDecimal("1"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("1"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("1"),
		},
		{
			"system_2_with_one_loss",
			callback.SystemBetType,
			[]int{2},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusLoss,
				},
			}),
			toDecimal("10.416666"),
		},
		{
			"system_2/6_with_all",
			callback.SystemBetType,
			[]int{2},
			toDecimal("5"),
			toSelections([]*callback.Odd{
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
				{
					OddRatio:  toDecimal("2.5"),
					OddStatus: sportsbook.OddStatusWin,
				},
			}),
			toDecimal("31.249999"),
		},
	}

	calculator := NewRefundCalc(zap.NewNop(), former.FormExpresses[Selection])

	for _, test := range table {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			refund, err := calculator.Calc(test.betType, test.betSize, test.stake, test.selections)
			if err != nil {
				t.Errorf("failed to calculate refund: %s", err)
				return
			}

			if refund.Cmp(test.refund) != 0 {
				t.Errorf("invalid refund sum. expected: %s. got: %s", test.refund.String(), refund.String())
			}
		})
	}
}

func toDecimal(str string) *apd.Decimal {
	d, _, err := apd.NewFromString(str)
	if err != nil {
		panic(err)
	}

	return d
}

func toSelections[V Selection](values []V) []Selection {
	result := make([]Selection, len(values))
	for i := range values {
		result[i] = values[i]
	}

	return result
}
