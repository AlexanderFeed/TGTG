package store

import "time"

type Impact struct {
	Rescued     int     `json:"rescued"`
	SavedRubles int64   `json:"savedRubles"`
	CO2Kg       float64 `json:"co2Kg"`
}

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

type Challenge struct {
	ID          string
	Email       string
	Purpose     string
	PendingName string
	CodeHash    string
	MaxAttempts int
	ExpiresAt   time.Time
}

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
	CreatedBy            string    `json:"-"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type OfferFilter struct {
	Query    string
	Category string
	Limit    int
	Offset   int
}
