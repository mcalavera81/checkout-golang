package test

import (
	"github.com/stretchr/testify/mock"
	m "checkout-service/server/models"
	"checkout-service/server/repository"
)

type FakeRepo struct {
	mock.Mock
}

func (fake *FakeRepo) UpdateItem(id m.BasketId, item m.CatalogItem) (*m.Basket, error) {
	args := fake.Called(id, item)
	return args.Get(0).(*m.Basket),args.Error(1)
}

func (fake *FakeRepo) SetDiscounts(id m.BasketId, discounts []m.Discount) (*m.Basket, error) {
	args := fake.Called(id, discounts)
	return args.Get(0).(*m.Basket),args.Error(1)
}

func (fake *FakeRepo) GetCatalogItem(code m.Code) (m.CatalogItem, error) {
	args := fake.Called(code)
	return args.Get(0).(m.CatalogItem),args.Error(1)
}

func (fake *FakeRepo) Get(id m.BasketId) (*m.Basket, error) {
	args := fake.Called(id)
	return args.Get(0).(*m.Basket),args.Error(1)
}

func (fake *FakeRepo) Create() (*m.Basket, error) {
	args := fake.Called()
	return args.Get(0).(*m.Basket),args.Error(1)
}

func (fake *FakeRepo) Delete(id m.BasketId) error {
	args := fake.Called(id)
	return args.Error(0)
}


func (fake *FakeRepo) Total(id m.BasketId) (float64, error) {
	args := fake.Called(id)
	return args.Get(0).(float64),args.Error(1)
}





var catalogRepo repository.CatalogRepository

func init()  {
	catalogRepo = repository.NewCatalogRepo()
}