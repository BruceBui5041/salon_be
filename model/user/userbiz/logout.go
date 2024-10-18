package userbiz

import (
	"context"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"

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
