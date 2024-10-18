package generictransport

import "salon_be/component"

type GenericTransport struct {
	AppContext component.AppContext
}

func NewGenericTransport(appCtx component.AppContext) *GenericTransport {
	return &GenericTransport{
		AppContext: appCtx,
	}
}
