package repository


import (
	m "checkout-service/server/models"
)


type CatalogRepository interface {
	GetCatalogItem(code m.Code) (m.CatalogItem, error)
}


type CatalogRepoImpl struct{
	catalog m.Catalog
}

func (r *CatalogRepoImpl) GetCatalogItem(code m.Code) (m.CatalogItem, error){
	item, ok:= r.catalog[code]
	if !ok {
		return m.CatalogItem{}, m.NotFound
	}else{
		item.Name = ""
		return item,nil
	}
}


func NewCatalogRepo() *CatalogRepoImpl{
	repo := new(CatalogRepoImpl)
	repo.catalog = 	map[m.Code]m.CatalogItem{
		m.VOUCHER: {Code: m.VOUCHER, Name: "Cabify Voucher", Price: 5.00},
		m.TSHIRT:  {Code: m.TSHIRT, Name: "Cabify T-Shirt", Price: 20.00},
		m.MUG:     {Code: m.MUG, Name: "Cabify Coffee Mug", Price: 7.50},
	}
	return repo
}

