package database

import (
	"web_server/internal/entity"

	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{
		DB: db,
	}
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) GetByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := p.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.GetByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *Product) Delete(id string) error {
	_, err := p.GetByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(&entity.Product{}, "id = ?", id).Error
}

func (p *Product) GetAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = p.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error

	} else {
		err = p.DB.Find(&products).Error
	}

	return products, err
}
