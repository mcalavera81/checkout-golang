package test

import (
	r "checkout-service/server/repository"
	"github.com/stretchr/testify/suite"
	s "checkout-service/server/business"
	"checkout-service/server/models"
	"testing"
)

type TestSuiteManager struct {
	suite.Suite
	s.CheckoutService
	r.CatalogRepository
}

func (t *TestSuiteManager) SetupTest() {
	t.CatalogRepository = r.NewCatalogRepo()
	t.CheckoutService = s.NewCheckoutService(r.NewCheckoutRepo(), t.CatalogRepository)
}

func (t *TestSuiteManager) TestCreateEmptyBasket(){

	basket, _ :=t.Create()
	t.NotEmpty(basket.Id,"Basket missing id")
	t.Empty(basket.Items,"Basket items should be empty")
}


func (t *TestSuiteManager) TestGetNonExistentBasket(){
	basket, err := t.Get(models.NewBasketId(-100))
	t.Empty(err, "Basket should not exist")
	t.Equal(basket, &models.EmptyBasket)
}


func (t *TestSuiteManager) TestGetBasket(){
	basket, err := t.Get(models.NewBasketId(1))
	t.Empty(err, "Basket should not exist")
	t.Equal(basket, &models.EmptyBasket)
}

func (t *TestSuiteManager) TestDeleteBasket(){
	basket, err := t.Get(models.NewBasketId(1))
	t.Empty(err, "Basket should not exist")
	t.Equal(basket, &models.EmptyBasket)

	t.Delete(basket.BasketId)
	basketDb, _ := t.Get(models.NewBasketId(1))

	t.Equal(&models.EmptyBasket,basketDb)

}


func (t *TestSuiteManager) TestScanProducts() {
	b, _ :=t.Create()
	item1, _ := t.GetCatalogItem(models.VOUCHER)
	item2, _ := t.GetCatalogItem(models.MUG)
	t.Update(b.BasketId, &item1)
	t.Update(b.BasketId, &item2)

	basketDB, _ := t.Get(b.BasketId)
	t.Equal(b.Id, basketDB.Id)
	t.Equal(1,len(basketDB.Items[models.VOUCHER]))
	t.Equal(1,len(basketDB.Items[models.MUG]))
	t.Equal(0,len(basketDB.Items[models.TSHIRT]))
}


func (t *TestSuiteManager) TestGetTotal() {
	b, _ :=t.Create()
	item1, _ := t.GetCatalogItem(models.VOUCHER)
	item2, _ := t.GetCatalogItem(models.MUG)
	t.Update(b.BasketId, &item1)
	t.Update(b.BasketId, &item2)

	total,_ := t.Total(b.BasketId)
	voucher,_:=t.GetCatalogItem(models.VOUCHER)
	mug,_:=t.GetCatalogItem(models.MUG)
	t.InDelta( voucher.Price+mug.Price, total, 0.01)
}



func (t *TestSuiteManager) TestGetTotalVoucherDiscount() {

	b, _ :=t.Create()
	voucher, _ := t.GetCatalogItem(models.VOUCHER)

	t.Update(b.BasketId, &voucher)
	t.Update(b.BasketId, &voucher)
	t.Update(b.BasketId, &voucher)

	total,_ := t.Total(b.BasketId)


	disc := t.DiscountService().Get(models.BUY2_VOUCHERS_GET1).(models.BuyXGetY)
	item, _ := t.GetCatalogItem(models.VOUCHER)

	t.InDelta( item.Price*3-item.Price*float64(disc.BuyAmount-disc.GetAmount), total, 0.01)
}

func (t *TestSuiteManager) TestGetTotalTshirtDiscount() {

	b, _ :=t.Create()
	tshirt, _ := t.GetCatalogItem(models.TSHIRT)

	t.Update(b.BasketId, &tshirt)
	t.Update(b.BasketId, &tshirt)
	t.Update(b.BasketId, &tshirt)

	total,_ := t.Total(b.BasketId)

	discount := t.DiscountService().Get(models.TSHIRT_BULK).(models.Bulk)
	item, _ := t.GetCatalogItem(models.TSHIRT)
	t.InDelta( (item.Price- discount.Discount)*3, total, 0.01)
}


func TestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteManager))
}