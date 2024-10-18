package categorymodel

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	models "salon_be/model"

	"github.com/samber/lo"
)

func init() {
	modelhelper.RegisterResponseType(models.Category{}.TableName(), CategoryResponse{})
}

type CategoryResponse struct {
	common.SQLModel `json:",inline"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Services        []models.ServiceVersion `json:"service,omitempty"`
	ServiceCount    int                     `json:"service_count"`
}

func (cr *CategoryResponse) CountService() {
	cr.ServiceCount = lo.CountBy(cr.Services, func(service models.ServiceVersion) bool {
		return service.Status == "active"
	})
}

func (cr *CategoryResponse) RemoveServicesResponse() {
	cr.Services = []models.ServiceVersion{}
}

// func (cr *CategoryResponse) Mask(isAdmin bool) {
// 	cr.GenUID(common.DBTypeCategory)
// }
