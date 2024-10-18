package generictransport

import "video_server/component"

type GenericTransport struct {
	AppContext component.AppContext
}

func NewGenericTransport(appCtx component.AppContext) *GenericTransport {
	return &GenericTransport{
		AppContext: appCtx,
	}
}
