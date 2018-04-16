package test

import (
	"testing"
	"net/http"
	"checkout-service/server/business"
	m "checkout-service/server/models"
	"checkout-service/server/api"
	"checkout-service/server/repository"
	"fmt"
	"github.com/stretchr/testify/mock"
)


func TestCheckout(t *testing.T) {

	router := newRouter()
	fakeRepo := new(FakeRepo)
	sequence := m.Sequence{}
	catalogRepo := repository.NewCatalogRepo()
	basketGen := repository.BasketGenerator(catalogRepo)

	api.ServeCheckoutResource(router, business.NewCheckoutService(fakeRepo, catalogRepo))

	Id1 :=m.NewBasketId(sequence.Next())
	Id2 :=m.NewBasketId(sequence.Next())
	Id3 :=m.NewBasketId(sequence.Next())
	Id4 :=m.NewBasketId(sequence.Next())
	Id5 :=m.NewBasketId(sequence.Next())
	Id6 :=m.NewBasketId(sequence.Next())

	fakeRepo.On("Get", Id5).Return(&m.EmptyBasket,nil)
	fakeRepo.On("Get", Id1).Return(basketGen(Id1, map[m.Code]int{m.MUG:3}),nil)
	fakeRepo.On("Get", Id2).Return(basketGen(Id2, map[m.Code]int{m.VOUCHER:2}),nil)
	fakeRepo.On("Get", Id4).Return(basketGen(Id4, map[m.Code]int{m.VOUCHER:2, m.TSHIRT:2}),nil)
	fakeRepo.On("Create").Return(basketGen(Id6, nil),nil)
	fakeRepo.On("Delete", Id1).Return(nil)
	fakeRepo.On("Delete", Id3).Return(m.NotFound)
	fakeRepo.On("UpdateItem", Id1, matchItemType(m.VOUCHER)).
		Return(basketGen(Id1, map[m.Code]int{m.MUG:3,m.VOUCHER:1}),nil)
	fakeRepo.On("SetDiscounts",Id1, mock.Anything).Return(basketGen(Id1, map[m.Code]int{m.MUG:3,m.VOUCHER:1}),nil)

	fakeRepo.On("UpdateItem", Id2, matchItemType(m.VOUCHER) ).
		Return(basketGen(Id2, map[m.Code]int{m.VOUCHER:3}),nil)
	fakeRepo.On("SetDiscounts",Id2, mock.Anything).Return(basketGen(Id2, map[m.Code]int{m.VOUCHER:3}),nil)


	runAPITests(t, router, []apiTestCase{
		{tag: "t1 - get a nonexistent basket", method: "GET", url: "/api/checkout/5", status: http.StatusNotFound, response: `{"message":"Basket 5 not found", "Details":"Basket Not found!", "status":404}`},
		{tag: "t2 - get an existing basket", method: "GET", url: "/api/checkout/1", status: http.StatusOK, response: `{"Id":"1", "Items":{"MUG":3}}`},
		{tag: "t3 - create a basket", method: "POST", url: "/api/checkout/", status: http.StatusCreated, response: fmt.Sprintf(`{"Id": "%v", "Items":{}}`,Id6.Id)},
		{tag: "t4 - delete a nonexistent basket", method: "DELETE", url: "/api/checkout/3", status: http.StatusNotFound},
		{tag: "t4 - delete a basket", method: "DELETE", url: "/api/checkout/1", status: http.StatusOK},
		{tag: "t5 - update a basket; new item", method: "PUT", url: "/api/checkout/1", body: `{"code":"VOUCHER"}`, status: http.StatusOK, response: `{"Id":"1", "Items":{"MUG":3,"VOUCHER":1}}`},
		{tag: "t6 - update a basket; inc count", method: "PUT", url: "/api/checkout/2", body: `{"code":"VOUCHER"}`, status: http.StatusOK, response: `{"Id":"2", "Items":{"VOUCHER":3}}`},
		{tag: "t7 - get total", method: "GET", url: "/api/checkout/4/total", status: http.StatusOK, response: `{"Id":"4", "Total":50}`},

	})
}




