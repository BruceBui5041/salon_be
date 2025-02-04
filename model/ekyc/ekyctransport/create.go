package ginkyc

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/ekyc/ekycbiz"
	"salon_be/model/ekyc/ekycmodel"
	"salon_be/model/ekyc/ekycrepo"
	"salon_be/model/ekyc/ekycstore"

	"github.com/gin-gonic/gin"
)

func CreateKYCProfile(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var request ekycmodel.CreateKYCProfileRequest
		if err := c.ShouldBind(&request); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Get current user
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsUser() {
			panic(common.ErrNoPermission(errors.New("only users can create KYC profiles")))
		}

		request.UserID = requester.GetUserId()

		store := ekycstore.NewSQLStore(appCtx.GetMainDBConnection())
		repo := ekycrepo.NewCreateKYCRepo(store, appCtx.GetEKYCClient())
		business := ekycbiz.NewCreateKYCBiz(repo)

		if err := business.CreateKYCProfile(c.Request.Context(), &request); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
