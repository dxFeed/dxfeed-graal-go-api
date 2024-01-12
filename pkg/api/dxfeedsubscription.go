package api

type DXFeedSubscription struct {
	eventListenerList []EventListener
}

func (s *DXFeedSubscription) IsClosed() bool {
	return true
}

func (s *DXFeedSubscription) AddListener(listener EventListener) {
	s.eventListenerList = append(s.eventListenerList, listener)
}

func (s *DXFeedSubscription) RemoveListener(listener EventListener) {
	s.eventListenerList = removeFromSlice(s.eventListenerList, listener)
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

func removeFromSlice(list []EventListener, observerToRemove EventListener) []EventListener {
	listLength := len(list)
	for i, observer := range list {
		if observerToRemove == observer {
			list[listLength-1], list[i] = list[i], list[listLength-1]
			return list[:listLength-1]
		}
	}
	return list
}
