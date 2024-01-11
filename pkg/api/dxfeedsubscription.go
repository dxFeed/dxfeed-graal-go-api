package api

type DXFeedSubscription struct {
}

func (s *DXFeedSubscription) IsClosed() bool {
	return true
}

func (s *DXFeedSubscription) AddSymbol(symbol any) {
}

func (s *DXFeedSubscription) AddSymbols(symbols ...any) {
}

func (s *DXFeedSubscription) RemoveSymbol(symbol any) {
}

func (s *DXFeedSubscription) RemoveSymbols(symbol ...any) {
}

func (s *DXFeedSubscription) Clear() {
}

func (s *DXFeedSubscription) Close() {
}
