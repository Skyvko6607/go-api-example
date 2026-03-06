package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Category struct {
	Name       string          `bson:"name,omitempty" json:"name"`
	ProductIDs []bson.ObjectID `bson:"product_ids,omitempty" json:"productIds"`
	Products   []Product       `bson:"products,omitempty" json:"products"`
	Tags       []string        `bson:"tags,omitempty" json:"tags"`
}
