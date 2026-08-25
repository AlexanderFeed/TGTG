package store

// This file defines the Go values exchanged between the store and HTTP layers.
// They are not active-record objects: methods that execute SQL remain on Store.

import "time"

// Impact contains profile summary values. These are currently zero because the
// order/rescue subsystem has not been implemented yet.
type Impact struct {
	Rescued     int     `json:"rescued"`
	SavedRubles int64   `json:"savedRubles"`
	CO2Kg       float64 `json:"co2Kg"`
}

// User is the safe public representation returned by the API. Notice that it
// contains no session token, challenge hash, or other authentication secret.
type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	City       string    `json:"city"`
	Role       string    `json:"role"`
	VerifiedAt time.Time `json:"verifiedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	Impact     Impact    `json:"impact"`
}

// Challenge is internal data used while requesting/verifying an email code.
// Its fields have no JSON tags because handlers never serialize it directly.
type Challenge struct {
	ID          string
	Email       string
	Purpose     string
	PendingName string
	CodeHash    string
	MaxAttempts int
	ExpiresAt   time.Time
}

// Offer is the database/domain representation of one rescue listing. JSON tags
// define the field names used when a handler encodes it for the frontend.
type Offer struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	Merchant             string    `json:"merchant"`
	Category             string    `json:"category"`
	Description          string    `json:"description"`
	Contents             string    `json:"contents"`
	ImageURL             string    `json:"imageUrl"`
	PriceKopecks         int64     `json:"priceKopecks"`
	OriginalPriceKopecks int64     `json:"originalPriceKopecks"`
	PickupStart          time.Time `json:"pickupStart"`
	PickupEnd            time.Time `json:"pickupEnd"`
	Quantity             int       `json:"quantity"`
	Address              string    `json:"address"`
	District             string    `json:"district"`
	Latitude             *float64  `json:"latitude,omitempty"`
	Longitude            *float64  `json:"longitude,omitempty"`
	Delivery             bool      `json:"delivery"`
	Status               string    `json:"status"`
	// CreatedBy is needed for authorization but hidden from public JSON.
	CreatedBy string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// OfferFilter carries optional list/search parameters from HTTP to SQL.
type OfferFilter struct {
	Query    string
	Category string
	Limit    int
	Offset   int
}
