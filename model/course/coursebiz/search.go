package coursebiz

import (
	"video_server/component/genericapi/genericbiz"
	"video_server/component/genericapi/genericstore"
)

type CourseBiz struct {
	genericbiz.GenericBiz
	// Add any course-specific fields here if needed
}

func NewSearchCourseBiz(store genericstore.GenericStoreInterface) *CourseBiz {
	return &CourseBiz{
		GenericBiz: genericbiz.NewGenericBiz(store),
	}
}
