package tradeplugin

import (
	"alaninnovates.com/trade-builder/tradeplugin/trade"
	"github.com/disgoorg/snowflake/v2"
)

type State struct {
	users map[snowflake.ID]*trade.Trade
}

func NewTradeService() *State {
	return &State{users: make(map[snowflake.ID]*trade.Trade)}
}

func (s *State) CreateTrade(userID snowflake.ID) *trade.Trade {
	h := trade.NewTrade()
	s.users[userID] = h
	return h
}

func (s *State) GetTrade(userID snowflake.ID) *trade.Trade {
	return s.users[userID]
}

func (s *State) TradeCount() int {
	return len(s.users)
}

func (s *State) Trades() map[snowflake.ID]*trade.Trade {
	return s.users
}
