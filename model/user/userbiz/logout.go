package userbiz

import (
	"context"
	"video_server/common"
	"video_server/component"
	"video_server/component/logger"

	"go.uber.org/zap"
)

type logoutBusiness struct {
	appCache component.AppCache
}

func NewLogoutBusiness(
	appCache component.AppCache,
) *logoutBusiness {
	return &logoutBusiness{
		appCache: appCache,
	}
}

func (business *logoutBusiness) Logout(ctx context.Context) error {
	user := ctx.Value(common.CurrentUser).(common.Requester)
	err := business.appCache.DeleteUserCache(ctx, user.GetFakeId())
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to delete user cache", zap.Error(err))
		return common.ErrInternal(err)
	}

	return nil
}
