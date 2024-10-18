package enrollmenttransport

import (
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/enrollment/enrollmentmodel"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateEnrollmentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input enrollmentmodel.CreateEnrollment

		if err := c.ShouldBind(&input); err != nil {
			panic(err)
		}

		// requester := c.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {

			return nil
		}); err != nil {
			panic(common.ErrInternal(err))
		}
	}
}
