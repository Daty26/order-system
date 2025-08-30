package service

import "github.com/Daty26/order-system/payment-service/internal/model"

type mockPaymnetRepo struct {
	SaveFN func(payment model.Payment) (model.Payment, error)
	GetAll func() ([]model.Payment, error)
}

func (m *mockPaymnetRepo) Save(payment model.Payment) (model.Payment, error) {
	return m.SaveFN(payment)
}
