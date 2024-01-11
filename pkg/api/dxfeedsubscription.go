package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

type DXFeedSubscription struct {
	sub               *native.DXFeedSubscription
	eventListenerList []common.EventListener
}

func (s *DXFeedSubscription) IsClosed() bool {
	return true
}

func (s *DXFeedSubscription) AddListener(listener common.EventListener) error {
	s.eventListenerList = append(s.eventListenerList, listener)
	return s.sub.AttachListener(listener)
}

func (s *DXFeedSubscription) RemoveListener(listener common.EventListener) {
	s.eventListenerList = removeFromSlice(s.eventListenerList, listener)
}

func (s *DXFeedSubscription) AddSymbol(symbol any) error {
	return s.sub.AddSymbol(symbol)
}

func (s *DXFeedSubscription) AddSymbols(symbols ...any) error {
	return s.sub.AddSymbols(symbols...)
}

func (s *DXFeedSubscription) RemoveSymbol(symbol any) error {
	return s.sub.RemoveSymbol(symbol)
}

func (s *DXFeedSubscription) RemoveSymbols(symbols ...any) error {
	return s.sub.RemoveSymbols(symbols...)
}

func (s *DXFeedSubscription) Clear() {
	s.sub.Clear()
}

func (s *DXFeedSubscription) Close() {
	s.sub.Close()
}

func removeFromSlice(list []common.EventListener, observerToRemove common.EventListener) []common.EventListener {
	listLength := len(list)
	for i, observer := range list {
		if observerToRemove == observer {
			list[listLength-1], list[i] = list[i], list[listLength-1]
			return list[:listLength-1]
		}
	}
	return list
}
