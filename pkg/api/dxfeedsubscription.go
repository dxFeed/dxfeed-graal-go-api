package api

import (
	"dxfeed-graal-go-api/pkg/native"
)

type DXFeedSubscription struct {
	sub               *native.DXFeedSubscription
	eventListenerList []native.EventListener
}

func (s *DXFeedSubscription) IsClosed() bool {
	return true
}

func (s *DXFeedSubscription) AddListener(listener native.EventListener) {
	s.eventListenerList = append(s.eventListenerList, listener)
	s.sub.AttachListener(listener)
}

func (s *DXFeedSubscription) RemoveListener(listener native.EventListener) {
	s.eventListenerList = removeFromSlice(s.eventListenerList, listener)
}

func (s *DXFeedSubscription) AddSymbol(symbol any) {
	strSymbol := symbol.(string)
	s.sub.AddSymbol(strSymbol)
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

func removeFromSlice(list []native.EventListener, observerToRemove native.EventListener) []native.EventListener {
	listLength := len(list)
	for i, observer := range list {
		if observerToRemove == observer {
			list[listLength-1], list[i] = list[i], list[listLength-1]
			return list[:listLength-1]
		}
	}
	return list
}
