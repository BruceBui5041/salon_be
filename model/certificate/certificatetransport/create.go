package certificatetransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/certificate/certificatebiz"
	"salon_be/model/certificate/certificatemodel"
	"salon_be/model/certificate/certificaterepo"
	"salon_be/model/certificate/certificatestore"

	"github.com/gin-gonic/gin"
)

func CreateCertificate(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input certificatemodel.CreateCertificateInput

		// Handle multipart form
		file, err := c.FormFile("file")
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		input.File = file
		input.Type = c.PostForm("type")

		if input.Type == "" {
			panic(common.ErrInvalidRequest(errors.New("type is required")))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsProvider() && !requester.IsAdmin() {
			panic(common.ErrNoPermission(errors.New("no permission to create certificate")))
		}

		input.CreatorID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		store := certificatestore.NewSQLStore(db)
		repo := certificaterepo.NewCreateRepo(store)
		biz := certificatebiz.NewCreateBiz(repo)

		if err := biz.CreateCertificate(c.Request.Context(), &input); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
