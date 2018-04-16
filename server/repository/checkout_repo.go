package repository

import (
	m "checkout-service/server/models"
)


type (
	CheckoutRepository interface {
		Get(id m.BasketId) (*m.Basket, error)
		Create() (*m.Basket, error)
		Delete(id m.BasketId) error
		UpdateItem(id m.BasketId, item m.CatalogItem) (*m.Basket, error)
		SetDiscounts(id m.BasketId, discounts []m.Discount) (*m.Basket, error)
	}


	CheckoutRepoImpl struct{
		baskets map[m.BasketId]*m.Basket
		m.Sequence
	}

)

func (r *CheckoutRepoImpl) Create() (*m.Basket, error) {
	id := m.NewBasketId(r.Next())
	newBasket := m.Basket{BasketId: id, Items:make(map[m.Code][]*m.CatalogItem)}
	r.baskets[id] = &newBasket
	return &newBasket,nil
}

func (r *CheckoutRepoImpl) Get(id m.BasketId) (*m.Basket, error) {
	basket,ok := r.baskets[id]

	if !ok{
		basket = &m.Basket{}
	}

	return basket, nil

}

func (r *CheckoutRepoImpl) Delete(id m.BasketId) error {

	if _, ok := r.baskets[id]; !ok {
		return m.NotFound
	}

	delete(r.baskets, id)
	return nil
}

func (r *CheckoutRepoImpl) SetDiscounts(id m.BasketId, discounts []m.Discount) (*m.Basket, error) {
	basket, ok := r.baskets[id]
	if !ok {
		return  nil,m.NotFound
	}

	basket.Discounts = discounts
	return basket,nil
}

func (r *CheckoutRepoImpl) UpdateItem(id m.BasketId, item m.CatalogItem) (*m.Basket, error) {

	basket, ok := r.baskets[id]
	if !ok {
		return  nil,m.NotFound
	}

	switch item.Code{
		case m.VOUCHER,m.TSHIRT,m.MUG:
			basket.Items[item.Code] = append(basket.Items[item.Code], &item)

		default:
			return nil,m.InvalidItem

	}
	return basket,nil
}


func NewCheckoutRepo() CheckoutRepository {
	repo := new(CheckoutRepoImpl)
	repo.baskets = make(map[m.BasketId]*m.Basket)
	return repo
}
