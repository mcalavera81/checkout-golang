package models


type Discount interface {
	Apply(basket *Basket) bool
	GetDiscount(basket *Basket) float64
}

type BuyXGetY struct {
	Type Code
	BuyAmount int
	GetAmount int
}

func (d BuyXGetY) GetDiscount(basket *Basket) float64 {
	items := basket.Items[d.Type]
	if len(items) > 0 {
		return items[0].Price
	}
	return 0
}

func (d BuyXGetY) Apply(basket *Basket) (applied bool)  {
	applied = false

	items := basket.Items[d.Type]
	for i:=0;  i< len(items)/(d.BuyAmount+d.GetAmount);i++ {
		applied = true
		for j :=0; j < d.GetAmount; j++ {
			basket.AddDiscount(d)
		}
	}

	return applied
}

type Bulk struct {
	Type Code
	Amount int
	Discount float64
}

func (d Bulk) GetDiscount(basket *Basket) float64 {
	items := basket.Items[d.Type]
	num := len(items)
	if num >= d.Amount {
		return float64(num)* d.Discount
	}
	return 0
}

func (d Bulk) Apply(basket *Basket) (applied bool) {
	applied=false

	items := basket.Items[d.Type]
	if len(items) == d.Amount {
		applied = true
		basket.AddDiscount(d)
	}

	return
}


type DiscountId string
const (
	BUY2_VOUCHERS_GET1 	DiscountId = "BUY2VOUCHERSGET1"
	TSHIRT_BULK 		DiscountId = "TSHIRTBULK"

)




