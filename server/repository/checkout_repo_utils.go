package repository

import m "checkout-service/server/models"

type BasketGen func(basketId m.BasketId, items map[m.Code]int) *m.Basket

func BasketGenerator(catalogRepo CatalogRepository) BasketGen{
	return func(basketId m.BasketId, items map[m.Code]int) *m.Basket{
		basket := new(m.Basket)
		basket.BasketId = basketId

		basket.Items = make(map[m.Code][](*m.CatalogItem))

		for k,v := range items {
			for i := 0;  i<v; i++ {
				item, _ := catalogRepo.GetCatalogItem(k)
				basket.Items[k] = append(basket.Items[k] , &item)
			}
		}
		return basket
	}

}

