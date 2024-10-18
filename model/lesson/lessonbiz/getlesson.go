package lessonbiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lesson/lessonmodel"

	"github.com/jinzhu/copier"
)

type GetLessonStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Lesson, error)
}

type getLessonBiz struct {
	store GetLessonStore
}

func NewGetLessonBiz(store GetLessonStore) *getLessonBiz {
	return &getLessonBiz{store: store}
}

func (biz *getLessonBiz) GetLessonByID(ctx context.Context, id string) (*lessonmodel.LessonResponse, error) {
	uid, err := common.FromBase58(id)
	if err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	lesson, err := biz.store.FindOne(ctx, map[string]interface{}{"id": uid.GetLocalID()})
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.LessonEntityName, err)
	}

	var lessoRes lessonmodel.LessonResponse
	copier.Copy(&lessoRes, lesson)

	return &lessoRes, nil
}
