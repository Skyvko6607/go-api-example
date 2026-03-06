package models

import (
	"time"

	"TestAPI/enums"

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
}
