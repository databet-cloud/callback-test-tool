package sportsbook

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

type (
	OddStatus      string
	MatchStatus    string
	CompetitorType string
)

const (
	OddStatusNotResulted      OddStatus = "NOT_RESULTED"
	OddStatusWin              OddStatus = "WIN"
	OddStatusLoss             OddStatus = "LOSS"
	OddStatusHalfWin          OddStatus = "HALF_WIN"
	OddStatusHalfLoss         OddStatus = "HALF_LOSS"
	OddStatusRefunded         OddStatus = "REFUNDED"
	OddStatusCancelled        OddStatus = "CANCELLED"
	OddStatusRefundedManually OddStatus = "REFUNDED_MANUALLY"

	MatchStatusNotStarted MatchStatus = "NOT_STARTED"
	MatchStatusLive       MatchStatus = "LIVE"
	MatchStatusSuspended  MatchStatus = "SUSPENDED"
	MatchStatusEnded      MatchStatus = "ENDED"
	MatchStatusClosed     MatchStatus = "CLOSED"
	MatchStatusCancelled  MatchStatus = "CANCELLED"
	MatchStatusAbandoned  MatchStatus = "ABANDONED"
	MatchStatusDelayed    MatchStatus = "DELAYED"
	MatchStatusUnknown    MatchStatus = "UNKNOWN"
)

func (s MatchStatus) Int() int {
	switch s {
	case MatchStatusNotStarted:
		return 0
	case MatchStatusLive:
		return 1
	case MatchStatusSuspended:
		return 2
	case MatchStatusEnded:
		return 3
	case MatchStatusClosed:
		return 4
	case MatchStatusCancelled:
		return 5
	case MatchStatusAbandoned:
		return 6
	case MatchStatusDelayed:
		return 7
	}

	return 8 // Default case for UNKNOWN
}

func (s OddStatus) Int() int {
	switch s {
	case OddStatusNotResulted:
		return 0
	case OddStatusWin:
		return 1
	case OddStatusLoss:
		return 2
	case OddStatusHalfWin:
		return 3
	case OddStatusHalfLoss:
		return 4
	case OddStatusRefunded:
		return 5
	case OddStatusCancelled:
		return 6
	case OddStatusRefundedManually:
		return 7
	}

	panic("invalid oddStatus")
}

func (s OddStatus) String() string {
	return string(s)
}

func (s OddStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Int())
}

func (s MatchStatus) String() string {
	return strings.ToLower(string(s))
}

func (t CompetitorType) Int() int {
	switch string(t) {
	case "PERSON":
		return 1
	case "TEAM":
		return 2
	}

	return 0
}

type Competitor struct {
	Id   string         `json:"id"`
	Type CompetitorType `json:"type"`
}

type Fixture struct {
	StartTime  time.Time   `json:"startTime"`
	SportId    string      `json:"sportId"`
	Status     MatchStatus `json:"status"`
	Tournament struct {
		Id      string `json:"id"`
		SportId string `json:"sportId"`
	} `json:"tournament"`
	Competitors []*Competitor `json:"competitors"`
}

func (f *Fixture) GetCompetitor(id string) (*Competitor, bool) {
	for _, competitor := range f.Competitors {
		if competitor.Id == id {
			return competitor, true
		}
	}

	return nil, false
}

type Market struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	TypeId int    `json:"typeId"`
	Odds   []Odd  `json:"odds"`
}

type Odd struct {
	ID            string       `json:"id"`
	Value         *apd.Decimal `json:"value"`
	Status        OddStatus    `json:"status"`
	CompetitorIds []string     `json:"competitorIds"`
}

type SportEvent struct {
	ID         string   `json:"ID"`
	ProviderId string   `json:"providerId"`
	Fixture    Fixture  `json:"fixture"`
	Markets    []Market `json:"markets"`
}

type SportEventListByFilters struct {
	SportEventListByFilters struct {
		SportEvents []SportEvent `json:"sportEvents"`
	} `json:"sportEventListByFilters"`
}
