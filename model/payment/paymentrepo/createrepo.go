package paymentrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/payment/paymentmodel"
)

type CreatePaymentStore interface {
	Create(
		ctx context.Context,
		newPayment *models.Payment,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Payment, error)
}

type createPaymentRepo struct {
	store CreatePaymentStore
}

func NewCreatePaymentRepo(store CreatePaymentStore) *createPaymentRepo {
	return &createPaymentRepo{store: store}
}

func (repo *createPaymentRepo) CreateNewPayment(
	ctx context.Context,
	input *paymentmodel.CreatePayment,
) (*models.Payment, error) {
	newPayment := &models.Payment{
		UserID:            input.UserID,
		Amount:            input.Amount,
		Currency:          input.Currency,
		PaymentMethod:     input.PaymentMethod,
		TransactionID:     input.TransactionID,
		TransactionStatus: "completed",
	}

	paymentId, err := repo.store.Create(ctx, newPayment)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.PaymentEntityName, err)
	}

	payment, err := repo.store.FindOne(
		ctx,
		map[string]interface{}{"id": paymentId},
	)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.PaymentEntityName, err)
	}

	return payment, nil
}

func (repo *createPaymentRepo) CheckDuplicatePayment(ctx context.Context, userID uint32, transactionID string) (bool, error) {
	payment, err := repo.store.FindOne(ctx, map[string]interface{}{
		"user_id":        userID,
		"transaction_id": transactionID,
	})
	if err != nil {
		if err == common.RecordNotFound {
			return false, nil
		}
		return false, err
	}
	return payment != nil, nil
}
func (repo *createPaymentRepo) FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Payment, error) {
	return repo.store.FindOne(ctx, conditions, moreInfo...)
}
