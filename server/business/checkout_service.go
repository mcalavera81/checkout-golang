package business

import (
	m "checkout-service/server/models"
	"sync"
	r "checkout-service/server/repository"
)



type (
	CheckoutService interface {
		Get(id m.BasketId) (*m.Basket, error)
		Create() (*m.Basket, error)
		Delete(id m.BasketId) error
		Update(id m.BasketId, item *m.CatalogItem) (*m.Basket, error)
		Total(id m.BasketId) (float64, error)
		DiscountService() DiscountService
	}

	CheckoutServiceImpl struct {
		mutex        sync.RWMutex
		checkoutRepo r.CheckoutRepository
		catalogRepo r.CatalogRepository
		discounts DiscountService
	}
)

func NewCheckoutService(checkoutRepo r.CheckoutRepository, catalogRepo r.CatalogRepository) CheckoutService {
	return &CheckoutServiceImpl{
		checkoutRepo:checkoutRepo,
		catalogRepo:catalogRepo,
		discounts:NewDiscountService(),
	}
}

func (s *CheckoutServiceImpl) DiscountService() DiscountService{
	return s.discounts
}

func (s *CheckoutServiceImpl) Create() (*m.Basket, error) {
	s.mutex.Lock()
	s.mutex.Unlock()
	return s.checkoutRepo.Create();
}

func (s *CheckoutServiceImpl) Get(id m.BasketId) (*m.Basket, error) {
	s.mutex.RLock()
	s.mutex.RUnlock()
	return s.checkoutRepo.Get(id)
}
func (s *CheckoutServiceImpl) Delete(id m.BasketId) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.checkoutRepo.Delete(id);
}

func (s *CheckoutServiceImpl) Update(id m.BasketId, item *m.CatalogItem) (*m.Basket, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	basket, _ := s.checkoutRepo.Get(id)

	templateItem, _:= s.catalogRepo.GetCatalogItem(item.Code)
	basket, err := s.checkoutRepo.UpdateItem(id, templateItem)

	if err != nil {
		return nil, err
	}

	s.discounts.Apply(basket)
	s.checkoutRepo.SetDiscounts(id, basket.Discounts)

	return basket, nil
}


func (s *CheckoutServiceImpl) Total(id m.BasketId) (float64, error) {

	basket, err := s.checkoutRepo.Get(id)
	if err !=nil {
		return  -1, err
	}


	var totalDiscount float64

	for _, disc := range basket.Discounts {
		totalDiscount += disc.GetDiscount(basket)
	}

	var total float64
	for _, itemSlice := range basket.Items {
		for _,item := range itemSlice {
			total += item.Price
		}
	}

	total -= totalDiscount

	return total,nil

}
