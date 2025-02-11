package routes

import (
	"errors"
	"net/http"

	"salon_be/apihandler"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/cache"
	"salon_be/component/genericapi/generictransport"
	"salon_be/component/ws"
	"salon_be/middleware"
	"salon_be/model/booking/bookingtransport"
	"salon_be/model/category/categorytransport"
	"salon_be/model/certificate/certificatetransport"
	"salon_be/model/comment/commenttransport"
	"salon_be/model/commission/commissiontransport"
	"salon_be/model/coupon/coupontransport"
	"salon_be/model/ekyc/ekyctransport"
	"salon_be/model/otp/otptransport"
	"salon_be/model/payment/paymenttransport"
	"salon_be/model/permission/permissiontransport"
	"salon_be/model/role/roletransport"
	"salon_be/model/service/servicetransport"
	"salon_be/model/user/usertransport"
	"salon_be/model/userdevice/userdevicetransport"
	"salon_be/model/userprofile/userprofiletransport"
	"salon_be/model/video/videotransport"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func SetupRoutes(r *gin.Engine, appCtx component.AppContext) {
	wsServer := ws.NewWebSocketServer()
	r.GET("/ws", wsServer.HandleWebSocket(appCtx))

	setupGenericRoutes(r, appCtx)
	setupPermissionRoutes(r, appCtx)
	setupRoleRoutes(r, appCtx)
	setupUserRoutes(r, appCtx)
	setupUserProfileRoutes(r, appCtx)
	setupCommentRoutes(r, appCtx)
	setupVideoRoutes(r, appCtx)
	setupPaymentRoutes(r, appCtx)
	setupCategoryRoutes(r, appCtx)
	setupServiceRoutes(r, appCtx)
	setupOTPRoutes(r, appCtx)

	setupAuthRoutes(r, appCtx)
	setupUtilityRoutes(r, appCtx)
	setupBookingRoutes(r, appCtx)
	setupUserDeviceRoutes(r, appCtx)
	setupCouponRoutes(r, appCtx)
	setupCommissionRoutes(r, appCtx)
	setupKYCRouters(r, appCtx)
	setupCertificateRouters(r, appCtx)
}

func setupGenericRoutes(r *gin.Engine, appCtx component.AppContext) {
	genTransport := generictransport.NewGenericTransport(appCtx)
	genAPIs := r.Group("/")
	{
		genAPIs.POST("search", genTransport.Search())
	}
}

func setupBookingRoutes(r *gin.Engine, appCtx component.AppContext) {
	bookingGroup := r.Group("/booking", middleware.RequiredAuth(appCtx))
	{
		bookingGroup.POST("", bookingtransport.CreateBookingHandler(appCtx))
		bookingGroup.PATCH("cancel/:id", bookingtransport.CancelBookingHandler(appCtx))
		bookingGroup.PATCH("accept/:id", bookingtransport.AcceptBookingHandler(appCtx))
		bookingGroup.PATCH("complete/:id", bookingtransport.CompleteBookingHandler(appCtx))
	}
}

func setupOTPRoutes(r *gin.Engine, appCtx component.AppContext) {
	setContextData := func(c *gin.Context) {
		c.Set("isChallengeAPI", true)
		c.Next()
	}

	otpGroup := r.Group("/otp", setContextData, middleware.RequiredAuth(appCtx))
	{
		otpGroup.POST("/verify", otptransport.VerifyOTP(appCtx))
		otpGroup.POST("/resend", otptransport.ResendOTP(appCtx))
	}
}

func setupServiceRoutes(r *gin.Engine, appCtx component.AppContext) {
	service := r.Group("/service", middleware.RequiredAuth(appCtx))
	{
		service.POST("", servicetransport.CreateServiceHandler(appCtx))
		// service.PATCH("/:id", permissiontransport.UpdatePermissionHandler(appCtx))
		service.PUT("/:id", servicetransport.UpdateServiceHandler(appCtx))
		service.PATCH("/publish/:id", servicetransport.PublishServiceHandler(appCtx))
		service.PATCH("/uploadimages", servicetransport.UploadImagesHandler(appCtx))
	}
}

func setupPermissionRoutes(r *gin.Engine, appCtx component.AppContext) {
	permission := r.Group("/permission", middleware.RequiredAuth(appCtx))
	{
		permission.POST("", permissiontransport.CreatePermissionHandler(appCtx))
		permission.PATCH("/:id", permissiontransport.UpdatePermissionHandler(appCtx))
	}
}

func setupRoleRoutes(r *gin.Engine, appCtx component.AppContext) {
	role := r.Group("/role", middleware.RequiredAuth(appCtx))
	{
		role.POST("", roletransport.CreateRoleHandler(appCtx))
		role.PATCH("/:id", roletransport.UpdateRoleHandler(appCtx))
		role.DELETE("/:id", roletransport.SoftDeleteRoleHandler(appCtx))
	}
}

func setupUserRoutes(r *gin.Engine, appCtx component.AppContext) {
	user := r.Group("/user", middleware.RequiredAuth(appCtx))
	{
		user.GET("", usertransport.GetUser(appCtx))
		user.GET("/provider-earnings", usertransport.GetProviderEarnings(appCtx))
		user.PATCH("/:id", usertransport.UpdateUser(appCtx))
	}
}

func setupUserProfileRoutes(r *gin.Engine, appCtx component.AppContext) {
	userprofile := r.Group("/profile")
	{
		userprofile.POST("", middleware.RequiredAuth(appCtx), userprofiletransport.CreateUserProfileHandler(appCtx))
		userprofile.PUT("", middleware.RequiredAuth(appCtx), userprofiletransport.UpdateProfileHandler(appCtx))
	}
}

func setupCommentRoutes(r *gin.Engine, appCtx component.AppContext) {
	commentGroup := r.Group("/comment")
	{
		commentGroup.POST("", middleware.RequiredAuth(appCtx), commenttransport.CreateCommentHandler(appCtx))
		commentGroup.PUT("/:id", middleware.RequiredAuth(appCtx), commenttransport.UpdateCommentHandler(appCtx))
	}
}

func setupVideoRoutes(r *gin.Engine, appCtx component.AppContext) {
	videoGroupInstructor := r.Group("/video",
		middleware.RequiredAuth(appCtx),
		middleware.AllowIntructorOnly(appCtx),
	)
	{
		videoGroupInstructor.POST("", videotransport.CreateVideoHandler(appCtx))
		videoGroupInstructor.PUT("/:id", videotransport.UpdateVideoHandler(appCtx))
		videoGroupInstructor.GET("/:service_slug", videotransport.ListServiceVideos(appCtx))
	}

	video := r.Group("/video", middleware.RequiredAuth(appCtx))
	{
		video.GET("/playlist/:service_slug/:video_id", apihandler.GetPlaylistHandler(appCtx))
		video.GET("/playlist/:service_slug/:video_id/:resolution/:playlistName", apihandler.GetPlaylistHandler(appCtx))
		video.GET("", apihandler.SegmentHandler(appCtx))
	}
}

func setupPaymentRoutes(r *gin.Engine, appCtx component.AppContext) {
	paymentGroup := r.Group("/payment", middleware.RequiredAuth(appCtx))
	{
		paymentGroup.POST("", paymenttransport.CreatePaymentHandler(appCtx))
	}
}

func setupCategoryRoutes(r *gin.Engine, appCtx component.AppContext) {
	categoryGroup := r.Group("/category")
	{
		categoryGroup.GET("", categorytransport.ListCategories(appCtx))
		categoryGroup.PATCH("/:id", middleware.RequiredAuth(appCtx), categorytransport.UpdateCategoryHandler(appCtx))
		categoryGroup.POST("", middleware.RequiredAuth(appCtx), categorytransport.CreateCategoryHandler(appCtx))
	}
}

func setupUserDeviceRoutes(r *gin.Engine, appCtx component.AppContext) {
	userdeviceGroup := r.Group("/userdevice")
	{
		userdeviceGroup.POST("",
			middleware.RequiredAuth(appCtx),
			userdevicetransport.CreateUserDevice(appCtx),
		)
		userdeviceGroup.PATCH("",
			middleware.RequiredAuth(appCtx),
			userdevicetransport.CreateUserDevice(appCtx),
		)
	}
}

func setupCouponRoutes(r *gin.Engine, appCtx component.AppContext) {
	couponGroup := r.Group("/coupon", middleware.RequiredAuth(appCtx))
	{
		couponGroup.POST("", coupontransport.CreateCouponHandler(appCtx))
		couponGroup.PATCH("/:id", coupontransport.UpdateCouponHandler(appCtx))
	}
}

func setupCommissionRoutes(r *gin.Engine, appCtx component.AppContext) {
	commissionGroup := r.Group("/commission", middleware.RequiredAuth(appCtx))
	{
		commissionGroup.POST("", commissiontransport.CreateCommissionHandler(appCtx))
		commissionGroup.PATCH("/:id", commissiontransport.UpdateCommissionHandler(appCtx))
	}
}

func setupAuthRoutes(r *gin.Engine, appCtx component.AppContext) {
	r.POST("/login", usertransport.Login(appCtx))
	r.POST("/register", usertransport.Register(appCtx))
	r.POST("/logout", middleware.RequiredAuth(appCtx), usertransport.Logout(appCtx))

	r.GET("/checkauth", middleware.RequiredAuth(appCtx), func(c *gin.Context) {
		requester, ok := c.Request.Context().Value(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}
		var userCached cache.CacheUser
		copier.Copy(&userCached, requester)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(userCached))
	})
}

func setupKYCRouters(r *gin.Engine, appCtx component.AppContext) {
	eKYCGroup := r.Group("/ekyc", middleware.RequiredAuth(appCtx))
	{
		eKYCGroup.POST("", ekyctransport.CreateKYCProfile(appCtx))
	}
}

func setupCertificateRouters(r *gin.Engine, appCtx component.AppContext) {
	certGroup := r.Group("/certificate", middleware.RequiredAuth(appCtx))
	{
		certGroup.POST("", certificatetransport.CreateCertificate(appCtx))
	}
}

func setupUtilityRoutes(r *gin.Engine, appCtx component.AppContext) {
	// TODO: disable these APIs if not in dev environment
	r.GET("/decode/:id", apihandler.DecodeUID(appCtx))
	r.GET("/encode/:id/:dbtype", apihandler.EncodeUID(appCtx))
}
