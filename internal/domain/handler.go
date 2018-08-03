package domain

import "net/http"

type HandlerHelper interface {
	URLParamInt(*http.Request, string) (int, error)
}
