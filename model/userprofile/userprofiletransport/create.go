package userprofiletransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/userprofile/userprofilebiz"
	"video_server/model/userprofile/userprofilemodel"
	"video_server/model/userprofile/userprofilerepo"
	"video_server/model/userprofile/userprofilestore"

	"github.com/gin-gonic/gin"
)

func CreateUserProfileHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input userprofilemodel.CreateUserProfile
		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		input.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		userStore := userprofilestore.NewSQLStore(db)
		repo := userprofilerepo.NewCreateUserProfileRepo(userStore)
		userProfileBusiness := userprofilebiz.NewCreateUserProfileBiz(repo)

		userProfile, err := userProfileBusiness.CreateNewUserProfile(ctx.Request.Context(), &input)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(userProfile))
	}
}
