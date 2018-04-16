package business

import (
	m "checkout-service/server/models"
)

func NewDiscountService() DiscountService {
	return &DiscountServiceImpl{
		discounts: map[m.DiscountId]m.Discount{
			m.BUY2_VOUCHERS_GET1: m.BuyXGetY{Type: m.VOUCHER, BuyAmount:2,GetAmount:1},
			m.TSHIRT_BULK: m.Bulk{Type: m.TSHIRT, Amount:3, Discount: 1},
		},
	}
}


type DiscountService interface {
	Apply(basket *m.Basket)
	Get(id m.DiscountId)m.Discount
}

type DiscountServiceImpl struct {
	discounts map[m.DiscountId]m.Discount
}

func (s *DiscountServiceImpl) Get(id m.DiscountId) m.Discount {
	return s.discounts[id]
}

func (s *DiscountServiceImpl) Apply(basket *m.Basket){
	basket.Discounts = []m.Discount{}
	for _,d := range s.discounts {
		d.Apply(basket)
	}
}

