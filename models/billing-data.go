package models

type BillingData struct {
	FirstName string `bson:"first_name,omitempty" json:"firstName"`
	LastName  string `bson:"last_name,omitempty" json:"lastName"`
	Country   string `bson:"country,omitempty" json:"country"`
	Region    string `bson:"region,omitempty" json:"region"`
	City      string `bson:"city,omitempty" json:"city"`
	ZipCode   string `bson:"zip_code,omitempty" json:"zipCode"`
	Address1  string `bson:"address1,omitempty" json:"address1"`
	Address2  string `bson:"address2,omitempty" json:"address2"`
}
