package service

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"slices"
	"strconv"
	"time"

	"github.com/cockroachdb/apd/v3"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/databet-cloud/callback-test-tool/internal/balance"
	"github.com/databet-cloud/callback-test-tool/internal/calculator"
	"github.com/databet-cloud/callback-test-tool/internal/callback"
	"github.com/databet-cloud/callback-test-tool/internal/sportsbook"
	"github.com/databet-cloud/callback-test-tool/internal/storage"
)

type Service struct {
	playerID         string
	playerBalance    *balance.Service
	sportsBookClient *sportsbook.Client
	callbackClient   *callback.Client
	calculator       *calculator.Calculator

	// list of bets in actual state
	bets     *storage.Storage[*callback.Data]
	cashOuts *storage.Storage[*callback.Data]
	// all sent requests
	sentRequests *storage.Storage[*callback.Data]

	log *zap.Logger
}

func NewService(
	playerID string,
	playerBalance *balance.Service,
	sportsBookClient *sportsbook.Client,
	callbackClient *callback.Client,
	calc *calculator.Calculator,
	log *zap.Logger,
) *Service {
	return &Service{
		playerID:         playerID,
		playerBalance:    playerBalance,
		sportsBookClient: sportsBookClient,
		callbackClient:   callbackClient,
		calculator:       calc,
		bets:             storage.New[*callback.Data](100),
		cashOuts:         storage.New[*callback.Data](100),
		sentRequests:     storage.New[*callback.Data](400),
		log:              log,
	}
}

func (s *Service) PlaceBet(ctx context.Context, betType callback.BetType, amount float64) {
	sportEventsCount := 0

	switch betType {
	case callback.SingleBetType:
		sportEventsCount = 1
	case callback.ExpressBetType:
		sportEventsCount = 2
	case callback.SystemBetType:
		sportEventsCount = 3
	}

	if sportEventsCount == 0 {
		s.log.Error("invalid bet type")
		return
	}

	decimalAmount, err := apd.New(0, 0).SetFloat64(amount)
	if err != nil {
		s.log.Error("invalid amount", zap.Float64("amount", amount), zap.Error(err))
		return
	}

	sportEvents, err := s.sportsBookClient.SportEventsByFilter(ctx, 0, sportEventsCount)
	if err != nil {
		s.log.Error("failed to get sport events", zap.Error(err))
		return
	}

	data := s.generatePlaceBetData(betType, decimalAmount, sportEvents)

	s.log.Info("Player balance before request", zap.Any("balance", s.PlayerBalance()))

	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to send bet place", zap.Error(err))
		return
	}

	s.playerBalance.Hold(decimalAmount)
	s.sentRequests.Insert(data)
	s.bets.Insert(data)

	s.processResponse(response)

	s.log.Info("Expect balance after request", zap.Any("balance", s.PlayerBalance()))
}

func (s *Service) AcceptBet(ctx context.Context, betID string) {
	placedBetFunc := func(d *callback.Data) bool {
		return d.BetID == betID && d.RequestType == callback.BetPlaceRequestType
	}

	bet, ok := s.bets.Get(placedBetFunc)
	if !ok {
		s.log.Error("failed to find placed bet", zap.String("id", betID))
		return
	}

	acceptedBet := bet.WithRequestType(callback.BetAcceptRequestType).WithRequestID(uuid.New().String())

	response, err := s.callbackClient.SendCallback(ctx, acceptedBet)
	if err != nil {
		s.log.Error("failed to send bet accept", zap.Error(err))
		return
	}

	s.sentRequests.Insert(acceptedBet)

	if ok := s.bets.Replace(acceptedBet, placedBetFunc); !ok {
		s.log.Error("failed to replace placed bet", zap.String("id", betID))
	}

	s.processResponse(response)
}

func (s *Service) DeclineBet(ctx context.Context, betID string, restrictionType callback.RestrictionType) {
	betFindFunc := func(d *callback.Data) bool {
		return d.BetID == betID &&
			(d.RequestType == callback.BetPlaceRequestType || d.RequestType == callback.BetAcceptRequestType)
	}

	bet, ok := s.bets.Get(betFindFunc)
	if !ok {
		s.log.Error("failed to find bet", zap.String("id", betID))
		return
	}

	restriction, err := generateRestriction(restrictionType, bet)
	if err != nil {
		s.log.Error("failed to generate restriction", zap.Error(err))
		return
	}

	data := &callback.Data{
		RequestType:           callback.BetDeclineRequestType,
		PrivateStake:          bet.PrivateStake,
		PrivateOdds:           bet.PrivateOdds,
		PrivateBetType:        bet.PrivateBetType,
		PrivateBetSystemSizes: bet.PrivateBetSystemSizes,
		PrivateCashOutAmount:  bet.PrivateCashOutAmount,

		RequestID:    uuid.NewString(),
		BetID:        bet.BetID,
		BetPlayerID:  s.playerID,
		Restrictions: []callback.Restriction{restriction},
	}

	s.log.Info("Player balance before request", zap.Any("balance", s.PlayerBalance()))

	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to send bet decline", zap.Error(err))
		return
	}

	s.playerBalance.UnHold(bet.PrivateStake)
	s.log.Info("Expect balance after request", zap.Any("balance", s.PlayerBalance()))

	s.sentRequests.Insert(data)
	s.bets.Replace(data, betFindFunc)

	s.processResponse(response)
}

// nolint:funlen // extended limit of lines to handle all possible ways in the single function
func (s *Service) SettleBet(ctx context.Context, betID string, odds []*callback.Odd) {
	betFindFunc := func(d *callback.Data) bool {
		return d.BetID == betID &&
			(d.RequestType == callback.BetAcceptRequestType ||
				d.RequestType == callback.BetUnSettleRequestType ||
				d.RequestType == callback.BetCashOutOrdersAcceptedRequestType ||
				d.RequestType == callback.BetCashOutOrdersDeclinedRequestType)
	}

	bet, ok := s.bets.Get(betFindFunc)
	if !ok {
		s.log.Error("failed to find bet", zap.String("id", betID))
		return
	}

	settleAmount, settleType, err := s.calculator.Settle(
		bet.PrivateBetType,
		bet.PrivateBetSystemSizes,
		bet.PrivateStake,
		odds,
	)
	if err != nil {
		s.log.Error("failed to settle bet", zap.String("id", betID), zap.Error(err))
		return
	}

	// patch values after cash-out
	if bet.RequestType == callback.BetCashOutOrdersAcceptedRequestType {
		settleType = callback.LossSettleType
		settleAmount = apd.New(0, 0)
	}

	data := &callback.Data{
		RequestType:           callback.BetSettleRequestType,
		PrivateStake:          bet.PrivateStake,
		PrivateOdds:           bet.PrivateOdds,
		PrivateBetType:        bet.PrivateBetType,
		PrivateBetSystemSizes: bet.PrivateBetSystemSizes,
		PrivateCashOutAmount:  bet.PrivateCashOutAmount,

		RequestID:    uuid.NewString(),
		BetID:        bet.BetID,
		BetPlayerID:  s.playerID,
		BetOdds:      odds,
		SettleAmount: formatApd(settleAmount),
		SettleType:   settleType,
	}

	s.log.Info("Player balance before request", zap.Any("balance", s.PlayerBalance()))

	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to send bet settle win", zap.Error(err))
		return
	}

	switch {
	case bet.RequestType == callback.BetCashOutOrdersAcceptedRequestType:
		// do nothing
	case settleType == callback.WinSettleType:
		s.playerBalance.WithdrawHold(bet.PrivateStake) // remove stake from hold
		s.playerBalance.Deposit(settleAmount)          // accrual
	case settleType == callback.RefundSettleType:
		s.playerBalance.UnHold(bet.PrivateStake) // return stake
	case settleType == callback.LossSettleType:
		s.playerBalance.WithdrawHold(bet.PrivateStake) // remove stake
	}

	s.log.Info("Expect balance after request", zap.Any("balance", s.PlayerBalance()))

	s.sentRequests.Insert(data)
	s.bets.Replace(data, betFindFunc)

	s.processResponse(response)
}

func (s *Service) UnSettleBet(ctx context.Context, betID string) {
	betFindFunc := func(d *callback.Data) bool {
		return d.BetID == betID && d.RequestType == callback.BetSettleRequestType
	}

	bet, ok := s.bets.Get(betFindFunc)
	if !ok {
		s.log.Error("failed to find bet", zap.String("id", betID))
		return
	}

	settleAmount, _, err := apd.NewFromString(bet.SettleAmount)
	if err != nil {
		s.log.Error("failed to parse settle amount", zap.String("id", betID))
		return
	}

	data := &callback.Data{
		RequestType:           callback.BetUnSettleRequestType,
		PrivateStake:          bet.PrivateStake,
		PrivateOdds:           bet.PrivateOdds,
		PrivateBetType:        bet.PrivateBetType,
		PrivateBetSystemSizes: bet.PrivateBetSystemSizes,
		PrivateCashOutAmount:  nil,
		RequestID:             uuid.NewString(),
		BetID:                 bet.BetID,
		BetPlayerID:           s.playerID,
		UnSettleAmount:        bet.SettleAmount,
	}

	s.log.Info("Player balance before request", zap.Any("balance", s.PlayerBalance()))

	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to send bet unsettle", zap.Error(err))
		return
	}

	switch {
	case bet.PrivateCashOutAmount != nil:
		s.playerBalance.DepositHold(bet.PrivateStake)
		s.playerBalance.Withdraw(bet.PrivateCashOutAmount)
	case bet.SettleType == callback.WinSettleType:
		s.playerBalance.Withdraw(settleAmount)
		s.playerBalance.DepositHold(bet.PrivateStake)
	case bet.SettleType == callback.RefundSettleType:
		s.playerBalance.Hold(bet.PrivateStake)
	case bet.SettleType == callback.LossSettleType:
		s.playerBalance.DepositHold(bet.PrivateStake)
	}

	s.log.Info("Expect balance after request", zap.Any("balance", s.PlayerBalance()))

	s.sentRequests.Insert(data)
	s.bets.Replace(data, betFindFunc)

	s.processResponse(response)
}

func (s *Service) AcceptBetCashOut(ctx context.Context, betID string) {
	betFindFunc := func(d *callback.Data) bool {
		return d.BetID == betID &&
			(d.RequestType == callback.BetAcceptRequestType || d.RequestType == callback.BetUnSettleRequestType)
	}

	bet, ok := s.bets.Get(betFindFunc)
	if !ok {
		s.log.Error("failed to find bet", zap.String("id", betID))
		return
	}

	// settle bet as half win
	cashOutAmount, err := s.settleOddsAs(bet, sportsbook.OddStatusHalfWin)
	if err != nil {
		s.log.Error(
			"failed to settle odds as half win",
			zap.String("id", betID),
			zap.Error(err),
		)

		return
	}

	data := &callback.Data{
		RequestType:           callback.BetCashOutOrdersAcceptedRequestType,
		PrivateStake:          bet.PrivateStake,
		PrivateOdds:           bet.PrivateOdds,
		PrivateBetType:        bet.PrivateBetType,
		PrivateBetSystemSizes: bet.PrivateBetSystemSizes,
		PrivateCashOutAmount:  cashOutAmount,

		RequestID:      uuid.NewString(),
		BetID:          bet.BetID,
		CashOutOrderID: uuid.NewString(),
		Amount:         formatApd(bet.PrivateStake),
		RefundAmount:   formatApd(cashOutAmount),
	}

	s.log.Info("Player balance before request", zap.Any("balance", s.PlayerBalance()))

	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to send accept bet cash out win", zap.Error(err))
		return
	}

	s.playerBalance.WithdrawHold(bet.PrivateStake) // remove from hold
	s.playerBalance.Deposit(cashOutAmount)         // deposit
	s.log.Info("Expect balance after request", zap.Any("balance", s.PlayerBalance()))

	s.sentRequests.Insert(data)
	s.bets.Replace(bet, betFindFunc)
	s.cashOuts.Insert(data)

	s.processResponse(response)
}

func (s *Service) DeclineBetCashOut(ctx context.Context, betID, cashOutOrderID string) {
	betFindFunc := func(d *callback.Data) bool {
		return d.BetID == betID
	}
	cashOutFindFunc := func(d *callback.Data) bool {
		return d.CashOutOrderID == cashOutOrderID
	}

	bet, ok := s.bets.Get(betFindFunc)
	if !ok {
		s.log.Error("failed to find bet", zap.String("id", betID))
		return
	}

	cashOut, ok := s.cashOuts.Get(cashOutFindFunc)
	if !ok {
		s.log.Error("failed to find cash-out", zap.String("id", cashOutOrderID))
		return
	}

	data := &callback.Data{
		RequestType:           callback.BetCashOutOrdersDeclinedRequestType,
		PrivateStake:          bet.PrivateStake,
		PrivateOdds:           bet.PrivateOdds,
		PrivateBetType:        bet.PrivateBetType,
		PrivateBetSystemSizes: bet.PrivateBetSystemSizes,
		PrivateCashOutAmount:  nil,

		RequestID:       uuid.NewString(),
		BetID:           bet.BetID,
		CashOutOrderIDs: []string{cashOut.CashOutOrderID},
	}

	s.log.Info("Player balance before request", zap.Any("balance", s.PlayerBalance()))

	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to send decline bet cash out win", zap.Error(err))
		return
	}

	s.playerBalance.DepositHold(cashOut.PrivateStake)
	s.playerBalance.Withdraw(cashOut.PrivateCashOutAmount)
	s.log.Info("Expect balance after request", zap.Any("balance", s.PlayerBalance()))

	data.CashOutOrderID = cashOutOrderID
	s.sentRequests.Insert(data)
	s.bets.Replace(bet, betFindFunc)
	s.cashOuts.Replace(data, cashOutFindFunc)

	s.processResponse(response)
}

func (s *Service) PlayerBalance() balance.Balance {
	return s.playerBalance.State()
}

func (s *Service) PlayerID() string {
	return s.playerID
}

func (s *Service) PlayerToken() string {
	return s.sportsBookClient.Token()
}

func (s *Service) SentRequests(types ...callback.RequestType) []*storage.Document[*callback.Data] {
	requests := s.sentRequests.GetDocuments(func(data *callback.Data) bool {
		return len(types) == 0 || slices.Contains(types, data.RequestType)
	})

	slices.Reverse(requests)

	return requests
}

func (s *Service) Bets(types ...callback.RequestType) []*storage.Document[*callback.Data] {
	docs := s.bets.GetDocuments(func(data *callback.Data) bool {
		return len(types) == 0 || slices.Contains(types, data.RequestType)
	})

	slices.Reverse(docs)

	return docs
}

func (s *Service) CashOuts(types ...callback.RequestType) []*storage.Document[*callback.Data] {
	docs := s.cashOuts.GetDocuments(func(data *callback.Data) bool {
		return len(types) == 0 || slices.Contains(types, data.RequestType)
	})

	slices.Reverse(docs)

	return docs
}

func (s *Service) ReplayCallback(ctx context.Context, data *callback.Data) {
	response, err := s.callbackClient.SendCallback(ctx, data)
	if err != nil {
		s.log.Error("failed to replay callback", zap.Any("data", data), zap.Error(err))
	}

	s.sentRequests.Replace(data, func(d *callback.Data) bool {
		return d.BetID == data.BetID && d.RequestID == data.RequestID
	})

	s.processResponse(response)
}

func (s *Service) processResponse(response *http.Response) {
	dumpResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		s.log.Error("failed to dump response", zap.Error(err))
		return
	}

	s.log.Debug("Callback response", zap.String("response", string(dumpResponse)))
}

// nolint:funlen // extended limit of lines to handle all bet types in the single function
func (s *Service) generatePlaceBetData(
	betType callback.BetType,
	amount *apd.Decimal,
	sportEvents []sportsbook.SportEvent,
) *callback.Data {
	odds := make([]*callback.Odd, 0, len(sportEvents))
	allCompetitors := make([]callback.Competitor, 0, len(sportEvents))

	for _, sportEvent := range sportEvents {
		market := randSelect(sportEvent.Markets)
		odd := randSelect(market.Odds)

		competitors := make([]callback.Competitor, 0, len(sportEvent.Fixture.Competitors))
		for _, cmp := range sportEvent.Fixture.Competitors {
			competitors = append(competitors, callback.Competitor{
				Id:   cmp.Id,
				Type: cmp.Type.Int(),
			})
		}

		sportEventInfoState := "prematch"
		if sportEvent.Fixture.Status == sportsbook.MatchStatusLive {
			sportEventInfoState = "live"
		}

		odds = append(odds, &callback.Odd{
			OddId:        odd.ID,
			OddRatio:     odd.Value,
			OddStatus:    odd.Status,
			MatchId:      sportEvent.ID,
			MatchStatus:  sportEvent.Fixture.Status.Int(),
			MarketId:     market.ID,
			OddUpdatedAt: time.Now().UTC().Add(-time.Hour),
			Meta: callback.OddMeta{
				MarketType:                 strconv.Itoa(market.TypeId),
				ProviderID:                 sportEvent.ProviderId,
				SportID:                    sportEvent.Fixture.SportId,
				TournamentID:               sportEvent.Fixture.Tournament.Id,
				SportEventInfoProviderId:   sportEvent.ProviderId,
				SportEventInfoSportId:      sportEvent.Fixture.SportId,
				SportEventInfoTournamentId: sportEvent.Fixture.Tournament.Id,
				SportEventInfoMarketType:   strconv.Itoa(market.TypeId),
				SportEventInfoState:        sportEventInfoState,
				SportEventInfoCompetitors:  competitors,
			},
		})

		allCompetitors = append(allCompetitors, competitors...)
	}

	systemSize := len(sportEvents)
	if betType == callback.SystemBetType {
		systemSize--
	}

	betCreatedAt := time.Now().UTC()

	return &callback.Data{
		RequestType:           callback.BetPlaceRequestType,
		PrivateStake:          amount,
		PrivateOdds:           odds,
		PrivateBetType:        betType,
		PrivateBetSystemSizes: []int{systemSize},

		RequestID:      uuid.NewString(),
		BetID:          xid.New().String(),
		BetPlayerID:    s.playerID,
		BetType:        betType,
		BetStake:       formatApd(amount),
		BetOdds:        odds,
		BetSystemSizes: []int{systemSize},
		BetCreatedAt:   &betCreatedAt,
		Competitors:    allCompetitors,
	}
}

func (s *Service) settleOddsAs(bet *callback.Data, oddStatus sportsbook.OddStatus) (*apd.Decimal, error) {
	odds := make([]*callback.Odd, len(bet.PrivateOdds))
	for i, odd := range bet.PrivateOdds {
		odds[i] = odd.WithStatus(oddStatus)
	}

	settleAmount, _, err := s.calculator.Settle(bet.PrivateBetType, bet.PrivateBetSystemSizes, bet.PrivateStake, odds)
	if err != nil {
		return nil, err
	}

	return settleAmount, err
}

func formatApd(v *apd.Decimal) string {
	return v.Text('f')
}

func randSelect[T any](values []T) T {
	if len(values) == 0 {
		var t T

		return t
	}

	return values[rand.Intn(len(values))]
}
