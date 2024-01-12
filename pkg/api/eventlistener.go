package api

type EventListener interface {
	Update([]interface{})
}
