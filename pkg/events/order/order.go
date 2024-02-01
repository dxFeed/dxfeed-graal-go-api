package order

type Order struct {
	Base
	marketMaker *string
}

func NewOrder(eventSymbol string) *Order {
	return &Order{Base: Base{eventSymbol: &eventSymbol}}
}
