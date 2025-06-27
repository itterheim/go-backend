package handler

type Handler interface {
	CreateRouter() []Route
}