package certificaterepo

import (
	"context"
	models "salon_be/model"
)

type CreateCertificateStore interface {
	Create(ctx context.Context, data *models.Certificate) error
}

type createRepo struct {
	store CreateCertificateStore
}

func NewCreateRepo(store CreateCertificateStore) *createRepo {
	return &createRepo{store: store}
}

func (r *createRepo) Create(ctx context.Context, data *models.Certificate) error {
	if err := r.store.Create(ctx, data); err != nil {
		return err
	}
	return nil
}
