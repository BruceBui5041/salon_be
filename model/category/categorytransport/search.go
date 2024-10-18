package categorytransport

import (
	"video_server/component"
	"video_server/component/genericapi/generictransport"
)

type CategoryTransport struct {
	generictransport.GenericTransport
}

func NewCourseTransport(appCtx component.AppContext) *CategoryTransport {
	return &CategoryTransport{
		GenericTransport: generictransport.GenericTransport{
			AppContext: appCtx,
		},
	}
}
