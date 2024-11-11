package callback

import (
	"slices"
	"time"

	"github.com/cockroachdb/apd/v3"

	"github.com/databet-cloud/callback-test-tool/internal/sportsbook"
)

type (
	BetType    int
	SettleType int
)

const (
	SingleBetType  BetType = 1
	ExpressBetType BetType = 2
	SystemBetType  BetType = 3

	WinSettleType    SettleType = 1
	RefundSettleType SettleType = 2
	LossSettleType   SettleType = 3
)

type Competitor struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
}

type OddMeta struct {
	MarketType                 string       `json:"market_type"`
	ProviderID                 string       `json:"provider_id"`
	SportID                    string       `json:"sport_id"`
	TournamentID               string       `json:"tournament_id"`
	SportEventInfoProviderId   string       `json:"sport_event_info_provider_id"`
	SportEventInfoSportId      string       `json:"sport_event_info_sport_id"`
	SportEventInfoTournamentId string       `json:"sport_event_info_tournament_id"`
	SportEventInfoMarketType   string       `json:"sport_event_info_market_type"`
	SportEventInfoState        string       `json:"sport_event_info_state"`
	SportEventInfoCompetitors  []Competitor `json:"sport_event_info_competitors,omitempty"`
}

type Odd struct {
	OddId        string               `json:"odd_id"`
	OddRatio     *apd.Decimal         `json:"odd_ratio"`
	OddStatus    sportsbook.OddStatus `json:"odd_status"`
	MatchId      string               `json:"match_id"`
	MatchStatus  int                  `json:"match_status"`
	MarketId     string               `json:"market_id"`
	OddUpdatedAt time.Time            `json:"odd_updated_at"`
	Meta         OddMeta              `json:"meta"`
	StatusReason string               `json:"status_reason"`
}

func (o *Odd) GetStatus() sportsbook.OddStatus {
	return o.OddStatus
}

func (o *Odd) GetValue() *apd.Decimal {
	return o.OddRatio
}

func (o *Odd) Clone() *Odd {
	return &Odd{
		OddId:        o.OddId,
		OddRatio:     o.OddRatio,
		OddStatus:    o.OddStatus,
		MatchId:      o.MatchId,
		MatchStatus:  o.MatchStatus,
		MarketId:     o.MarketId,
		OddUpdatedAt: o.OddUpdatedAt,
		Meta: OddMeta{
			MarketType:                 o.Meta.MarketType,
			ProviderID:                 o.Meta.ProviderID,
			SportID:                    o.Meta.SportID,
			TournamentID:               o.Meta.TournamentID,
			SportEventInfoProviderId:   o.Meta.SportEventInfoProviderId,
			SportEventInfoSportId:      o.Meta.SportEventInfoSportId,
			SportEventInfoTournamentId: o.Meta.SportEventInfoTournamentId,
			SportEventInfoMarketType:   o.Meta.SportEventInfoMarketType,
			SportEventInfoState:        o.Meta.SportEventInfoState,
			SportEventInfoCompetitors:  slices.Clone(o.Meta.SportEventInfoCompetitors),
		},
	}
}

func (o *Odd) WithStatus(status sportsbook.OddStatus) *Odd {
	odd := o.Clone()
	odd.OddStatus = status
	odd.OddUpdatedAt = time.Now()

	return odd
}

type Data struct {
	RequestType           RequestType  `json:"-"`
	PrivateStake          *apd.Decimal `json:"-"`
	PrivateOdds           []*Odd       `json:"-"`
	PrivateBetType        BetType      `json:"-"`
	PrivateBetSystemSizes []int        `json:"-"`
	PrivateCashOutAmount  *apd.Decimal `json:"-"`

	RequestID       string        `json:"request_id"`
	BetID           string        `json:"bet_id"`
	BetPlayerID     string        `json:"bet_player_id,omitempty"`
	BetType         BetType       `json:"bet_type,omitempty"`
	BetStake        string        `json:"bet_stake,omitempty"`
	BetFreeBetID    string        `json:"bet_freebet_id,omitempty"`
	BetInsuranceID  string        `json:"bet_insurance_id,omitempty"`
	BetOdds         []*Odd        `json:"bet_odds,omitempty"`
	BetSystemSizes  []int         `json:"bet_system_sizes,omitempty"`
	BetCreatedAt    *time.Time    `json:"bet_created_at,omitempty"`
	Competitors     []Competitor  `json:"competitors,omitempty"`
	SettleAmount    string        `json:"settle_amount,omitempty"`
	UnSettleAmount  string        `json:"unsettle_amount,omitempty"`
	SettleType      SettleType    `json:"settle_type,omitempty"`
	Restrictions    []Restriction `json:"restrictions,omitempty"`
	CashOutOrderID  string        `json:"cash_out_order_id,omitempty"`
	CashOutOrderIDs []string      `json:"cash_out_order_ids,omitempty"`
	Amount          string        `json:"amount,omitempty"`
	RefundAmount    string        `json:"refund_amount,omitempty"`
}

func (d *Data) WithRequestType(t RequestType) *Data {
	data := d.Clone()
	data.RequestType = t

	return data
}

func (d *Data) WithRequestID(id string) *Data {
	data := d.Clone()
	data.RequestID = id

	return data
}

func (d *Data) Clone() *Data {
	odds := make([]*Odd, len(d.BetOdds))
	for i, odd := range d.BetOdds {
		odds[i] = odd.Clone()
	}

	return &Data{
		RequestType:           d.RequestType,
		PrivateStake:          d.PrivateStake,
		PrivateOdds:           d.PrivateOdds,
		PrivateBetType:        d.PrivateBetType,
		PrivateBetSystemSizes: d.PrivateBetSystemSizes,
		PrivateCashOutAmount:  d.PrivateCashOutAmount,

		RequestID:       d.RequestID,
		BetID:           d.BetID,
		BetPlayerID:     d.BetPlayerID,
		BetType:         d.BetType,
		BetStake:        d.BetStake,
		BetFreeBetID:    d.BetFreeBetID,
		BetInsuranceID:  d.BetInsuranceID,
		BetOdds:         odds,
		BetSystemSizes:  slices.Clone(d.BetSystemSizes),
		BetCreatedAt:    d.BetCreatedAt,
		Competitors:     slices.Clone(d.Competitors),
		SettleAmount:    d.SettleAmount,
		UnSettleAmount:  d.UnSettleAmount,
		SettleType:      d.SettleType,
		Restrictions:    slices.Clone(d.Restrictions),
		CashOutOrderID:  d.CashOutOrderID,
		CashOutOrderIDs: slices.Clone(d.CashOutOrderIDs),
		Amount:          d.Amount,
		RefundAmount:    d.RefundAmount,
	}
}
