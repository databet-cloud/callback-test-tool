package callback

type Restriction struct {
	Type    RestrictionType        `json:"type"`
	Context map[string]interface{} `json:"context"`
}

type RestrictionType string

func (t RestrictionType) String() string {
	return string(t)
}

const (
	MaxBetRestriction                 RestrictionType = "max_bet"
	BetTypeRestriction                RestrictionType = "bet_type"
	BetIntervalRestriction            RestrictionType = "bet_interval"
	SelectionValueRestriction         RestrictionType = "selection_value"
	SportEventStatusRestriction       RestrictionType = "sport_event_status"
	SportEventExistenceRestriction    RestrictionType = "sport_event_existence"
	SportEventBetStopRestriction      RestrictionType = "sport_event_bet_stop"
	MarketStatusRestriction           RestrictionType = "market_status"
	MarketExistenceRestriction        RestrictionType = "market_existence"
	MarketDefectiveRestriction        RestrictionType = "market_defective"
	OddStatusRestriction              RestrictionType = "odd_status"
	OddExistenceRestriction           RestrictionType = "odd_existence"
	PlayerLimitRestriction            RestrictionType = "player_limit"
	FreebetNotApplicableRestriction   RestrictionType = "freebet_not_applicable"
	FreebetStatusRestriction          RestrictionType = "freebet_status"
	FreebetAmountRestriction          RestrictionType = "freebet_amount"
	InsuranceNotApplicableRestriction RestrictionType = "insurance_not_applicable"
	InsuranceStatusRestriction        RestrictionType = "insurance_status"
	InsuranceValueRestriction         RestrictionType = "insurance_value"
	InternalErrorRestriction          RestrictionType = "internal_error"
	MinBetRestriction                 RestrictionType = "min_bet"
	NotEnoughBalanceRestriction       RestrictionType = "not_enough_balance"
	WLDefinedRestriction              RestrictionType = "wl_defined"
)

func GetAllBetRestrictions() []RestrictionType {
	return []RestrictionType{
		MaxBetRestriction,
		BetTypeRestriction,
		BetIntervalRestriction,
		SelectionValueRestriction,
		SportEventStatusRestriction,
		SportEventExistenceRestriction,
		SportEventBetStopRestriction,
		MarketStatusRestriction,
		MarketExistenceRestriction,
		MarketDefectiveRestriction,
		OddStatusRestriction,
		OddExistenceRestriction,
		PlayerLimitRestriction,
		FreebetNotApplicableRestriction,
		FreebetStatusRestriction,
		FreebetAmountRestriction,
		InsuranceNotApplicableRestriction,
		InsuranceStatusRestriction,
		InsuranceValueRestriction,
		InternalErrorRestriction,
		MinBetRestriction,
		NotEnoughBalanceRestriction,
		WLDefinedRestriction,
	}
}
