package userutils

import (
	"context"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/model/user/userrepo"
	"salon_be/model/user/userstore"

	"go.uber.org/zap"
)

func UpdateUserCache(
	ctx context.Context,
	appCtx component.AppContext,
	userId string,
) error {
	appCtx.GetAppCache().DeleteUserCache(ctx, userId)

	store := userstore.NewSQLStore(appCtx.GetMainDBConnection())
	repo := userrepo.NewGetUserRepo(store)

	uid, err := common.FromBase58(userId)
	if err != nil {
		logger.AppLogger.Error(ctx,
			"Failed to parse uid",
			zap.Error(err),
			zap.Any("userId", userId),
		)
		return err
	}

	user, err := repo.GetUser(ctx, uid.GetLocalID())
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to get user", zap.Error(err))
		return err
	}

	user.Mask(false)
	for _, role := range user.Roles {
		role.Mask(false)
	}

	for _, enrollment := range user.Enrollments {
		enrollment.Mask(false)
		enrollment.Service.Mask(false)
	}

	err = appCtx.GetAppCache().SetUserCache(ctx, user)
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to cache user", zap.Error(err))
		return err
	}

	return nil
}
