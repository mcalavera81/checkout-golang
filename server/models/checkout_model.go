package models

import (
	"reflect"
	"errors"
	"strconv"
	"sync/atomic"
	"encoding/json"
)

type CatalogItem struct {
	Code   Code	`json:"code"`
	Name   string 	`json:"-"`
	Price  float64 	`json:"-"`
}

type BasketId struct {
	Id string
}

func NewBasketId(id int) BasketId{
	return BasketId{strconv.Itoa(id)}
}


var EmptyBasket = Basket{}


type Basket struct {
	BasketId
	Items     map[Code][]*CatalogItem
	Discounts []Discount
}

func (b *Basket) IsEmpty() bool {
	return reflect.DeepEqual(*b, EmptyBasket)
}


func (b *Basket) AddDiscount(discount Discount){
	b.Discounts = append(b.Discounts, discount)
}

func (b Basket) MarshalJSON() ([]byte, error) {

	var items = make(map[Code]int )
	for k,v := range b.Items {
		items[k] = len(v)
	}

	return json.Marshal(&struct {
		BasketId
		Items map[Code]int
	}{
		b.BasketId,
		items,
	})
}

//Add ISempty idiom
type Catalog map[Code]CatalogItem

type Code string
const (
	VOUCHER		Code = "VOUCHER"
	TSHIRT  	Code = "TSHIRT"
	MUG   	 	Code = "MUG"
)


type Sequence struct {
	sequence int32
}

func (s* Sequence) Next() int{
	next := atomic.AddInt32(&s.sequence, 1)
	return int(next)
}

var NotFound = errors.New("Basket Not found!")
var InvalidItem = errors.New("Invalid Item!")