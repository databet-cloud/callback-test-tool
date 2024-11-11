package service

import (
	"errors"

	"github.com/cockroachdb/apd/v3"

	"github.com/databet-cloud/callback-test-tool/internal/callback"
)

// nolint:funlen,gocyclo // its ok, because we generate all types of restrictions
func generateRestriction(t callback.RestrictionType, bet *callback.Data) (restriction callback.Restriction, err error) {
	odd := randSelect(bet.PrivateOdds)
	ctx := apd.BaseContext.WithPrecision(100)

	maxBet := apd.New(bet.PrivateStake.Coeff.Int64(), bet.PrivateStake.Exponent)
	for _, betOdd := range bet.BetOdds {
		_, err = ctx.Mul(maxBet, maxBet, betOdd.OddRatio)
		if err != nil {
			return restriction, err
		}
	}

	switch t {
	case callback.MaxBetRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(bet.PrivateStake),
				"sport_event_id": odd.MatchId,
			},
		}
	case callback.BetTypeRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"bet_type":       bet.BetType,
			},
		}
	case callback.BetIntervalRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"time_to_wait":   "10.00",
			},
		}
	case callback.SelectionValueRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"market_id":      odd.MarketId,
				"odd_id":         odd.OddId,
				"value":          odd.OddRatio.Text('f'),
			},
		}
	case callback.SportEventStatusRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"status":         odd.MatchStatus,
			},
		}
	case callback.SportEventExistenceRestriction, callback.SportEventBetStopRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
			},
		}
	case callback.MarketStatusRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"market_id":      odd.MarketId,
				"status":         odd.MatchStatus,
			},
		}
	case callback.MarketExistenceRestriction, callback.MarketDefectiveRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"market_id":      odd.MarketId,
			},
		}
	case callback.OddStatusRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"market_id":      odd.MarketId,
				"odd_id":         odd.OddId,
				"status":         odd.OddStatus,
				"is_active":      true,
			},
		}
	case callback.OddExistenceRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"market_id":      odd.MarketId,
				"odd_id":         odd.OddId,
			},
		}
	case callback.PlayerLimitRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"reason":         "limit_exceeded",
			},
		}
	case callback.FreebetNotApplicableRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"reason":         "not_found", // Змінити відповідно до даних
			},
		}
	case callback.FreebetStatusRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"status":         3,
			},
		}
	case callback.FreebetAmountRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":          formatApd(maxBet),
				"sport_event_id":   odd.MatchId,
				"freebet_amount":   "1.2",
				"freebet_currency": "USD",
				"bet_stake":        bet.PrivateStake.Text('f'),
				"bet_currency":     "EUR",
			},
		}
	case callback.InsuranceNotApplicableRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"reason":         "not_found",
			},
		}
	case callback.InsuranceStatusRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"status":         "used",
			},
		}
	case callback.InsuranceValueRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":            formatApd(maxBet),
				"sport_event_id":     odd.MatchId,
				"given_version":      "sfjgfd89844234h23l",
				"actual_version":     "eqwrqw89844234h23l",
				"bet_currency":       "USD",
				"insurance_currency": "EUR",
			},
		}
	case callback.InternalErrorRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"code":           "place_retry_limit_reached",
			},
		}
	case callback.MinBetRestriction:
		minBet := apd.New(0, 0)
		if _, err := ctx.Add(minBet, maxBet, bet.PrivateStake); err != nil {
			return restriction, err
		}

		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"min_bet":        formatApd(minBet),
			},
		}
	case callback.NotEnoughBalanceRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"actual_balance": "1.5",
			},
		}
	case callback.WLDefinedRestriction:
		restriction = callback.Restriction{
			Type: t,
			Context: map[string]interface{}{
				"max_bet":        formatApd(maxBet),
				"sport_event_id": odd.MatchId,
				"code":           "player_limit_reached",
			},
		}
	default:
		err = errors.New("unknown restriction type")
	}

	return restriction, err
}
