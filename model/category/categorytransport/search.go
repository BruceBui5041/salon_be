package categorytransport

import (
	"salon_be/component"
	"salon_be/component/genericapi/generictransport"
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
