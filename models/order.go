package models

import (
	"time"

	"github.com/Skyvko6607/go-api-example/enums"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Persist Products list, so if original product changes, it won't reflect the changes in the order
// ProductIDs are stored, to do statistics later on. Ex. do people buy some products for price trends
type Order struct {
	ID          bson.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID     `bson:"user_id,omitempty" json:"userId"`
	ProductIDs  []bson.ObjectID   `bson:"product_ids,omitempty" json:"productIds"`
	Products    []Product         `bson:"products,omitempty" json:"products"`
	TotalPaid   float32           `bson:"total_paid,omitempty" json:"totalPaid"`
	CreatedAt   time.Time         `bson:"created_at,omitempty" json:"createdAt"`
	FinalizedAt time.Time         `bson:"finalized_at,omitempty" json:"finalizedAt"`
	PaymentType enums.PaymentType `bson:"payment_type,omitempty" json:"paymentType"`
	BillingData BillingData       `bson:"billing_data,omitempty" json:"billingData"`
}

type OrderDTO struct {
	ID          bson.ObjectID     `json:"id"`
	Products    []ProductDTO      `json:"products"`
	CreatedAt   time.Time         `bson:"created_at,omitempty" json:"createdAt"`
	FinalizedAt time.Time         `bson:"finalized_at,omitempty" json:"finalizedAt"`
	TotalPaid   float32           `bson:"total_paid,omitempty" json:"totalPaid"`
	PaymentType enums.PaymentType `bson:"payment_type,omitempty" json:"paymentType"`
}

type CreateOrderDTO struct {
	Products    []bson.ObjectID   `json:"productIds"`
	PaymentType enums.PaymentType `json:"paymentType"`
	BillingData BillingData       `json:"billingData"`
}

func (o *Order) AsDTO() OrderDTO {
	var productsDTO []ProductDTO
	for _, product := range o.Products {
		productsDTO = append(productsDTO, product.AsDTO())
	}

	return OrderDTO{
		ID:          o.ID,
		Products:    productsDTO,
		TotalPaid:   o.TotalPaid,
		CreatedAt:   o.CreatedAt,
		FinalizedAt: o.FinalizedAt,
		PaymentType: o.PaymentType,
	}
}
