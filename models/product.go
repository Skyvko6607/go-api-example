package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name,omitempty" json:"name"`
	Description string        `bson:"description,omitempty" json:"description"`
	ImageUrl    string        `bson:"image_url,omitempty" json:"imageUrl"`
	Price       float32       `bson:"price,omitempty" json:"price"`
	Discount    float32       `bson:"discount,omitempty" json:"discount"`
}

type ProductDTO struct {
	ID          bson.ObjectID `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ImageUrl    string        `json:"imageUrl"`
	Price       float32       `json:"price"`
}

func (p *ProductDTO) AsBase() Product {
	return Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageUrl:    p.ImageUrl,
		Price:       p.Price,
		Discount:    0,
	}
}

func (p *Product) AsDTO() ProductDTO {
	return ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageUrl:    p.ImageUrl,
		Price:       p.Price * (1. - p.Discount),
	}
}
