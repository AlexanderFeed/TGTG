package store

// This file is the offers persistence layer. It contains parameterized SQL and
// converts database rows into Offer values; HTTP details stay in httpapi.

import (
	"context"
	"database/sql"
	"errors"
)

// offerColumns gives SELECT and RETURNING exactly the same column order. That
// order must match scanOffer below or values would be scanned into wrong fields.
const offerColumns = `
	id, title, merchant, category, description, contents, image_url,
	price_kopecks, original_price_kopecks, pickup_start, pickup_end,
	quantity, address, district, latitude, longitude, delivery, status,
	created_by, created_at, updated_at`

// rowScanner is a tiny interface satisfied by both *sql.Row and *sql.Rows.
// It lets one scanOffer helper work for a single-row query and a list query.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanOffer translates one SQL row into an Offer. Nullable database columns use
// sql.Null* during scanning, then become pointers/zero values in normal Go data.
func scanOffer(row rowScanner) (Offer, error) {
	var offer Offer
	var latitude, longitude sql.NullFloat64
	var createdBy sql.NullString
	err := row.Scan(
		&offer.ID, &offer.Title, &offer.Merchant, &offer.Category,
		&offer.Description, &offer.Contents, &offer.ImageURL,
		&offer.PriceKopecks, &offer.OriginalPriceKopecks,
		&offer.PickupStart, &offer.PickupEnd, &offer.Quantity,
		&offer.Address, &offer.District, &latitude, &longitude,
		&offer.Delivery, &offer.Status, &createdBy, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if latitude.Valid {
		offer.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		offer.Longitude = &longitude.Float64
	}
	if createdBy.Valid {
		offer.CreatedBy = createdBy.String
	}
	return offer, err
}

// ListOffers returns active offers matching optional text/category filters.
func (s *Store) ListOffers(ctx context.Context, filter OfferFilter) ([]Offer, error) {
	filter.Query = normalizeSearch(filter.Query)
	filter.Category = normalizeSearch(filter.Category)
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	// $1...$4 are PostgreSQL parameters. Passing their values separately from
	// the SQL text prevents user input from becoming executable SQL.
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+offerColumns+`
		FROM offers
		WHERE status = 'active'
		  AND ($1 = '' OR title ILIKE '%' || $1 || '%' OR merchant ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR category = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`, filter.Query, filter.Category, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	// Always return the connection to the pool, including on an early error.
	defer rows.Close()

	offers := make([]Offer, 0)
	// QueryContext returns a cursor. Next advances one row at a time.
	for rows.Next() {
		offer, err := scanOffer(rows)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}
	return offers, rows.Err()
}

// OfferByID loads one non-deleted offer. It converts database/sql's generic
// sql.ErrNoRows into the store package's domain-level ErrNotFound.
func (s *Store) OfferByID(ctx context.Context, id string) (Offer, error) {
	offer, err := scanOffer(s.db.QueryRowContext(ctx, `
		SELECT `+offerColumns+`
		FROM offers
		WHERE id = $1 AND status <> 'deleted'`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Offer{}, ErrNotFound
	}
	return offer, err
}

// CreateOffer inserts an offer and uses RETURNING to retrieve database defaults
// such as created_at without a second round trip.
func (s *Store) CreateOffer(ctx context.Context, offer Offer) (Offer, error) {
	return scanOffer(s.db.QueryRowContext(ctx, `
		INSERT INTO offers (
			id, title, merchant, category, description, contents, image_url,
			price_kopecks, original_price_kopecks, pickup_start, pickup_end,
			quantity, address, district, latitude, longitude, delivery, status, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19
		)
		RETURNING `+offerColumns,
		offer.ID, offer.Title, offer.Merchant, offer.Category, offer.Description,
		offer.Contents, offer.ImageURL, offer.PriceKopecks, offer.OriginalPriceKopecks,
		offer.PickupStart, offer.PickupEnd, offer.Quantity, offer.Address, offer.District,
		offer.Latitude, offer.Longitude, offer.Delivery, offer.Status, offer.CreatedBy))
}

// UpdateOffer changes editable fields while preserving ID/creator/created_at.
func (s *Store) UpdateOffer(ctx context.Context, offer Offer) (Offer, error) {
	updated, err := scanOffer(s.db.QueryRowContext(ctx, `
		UPDATE offers SET
			title = $2, merchant = $3, category = $4, description = $5,
			contents = $6, image_url = $7, price_kopecks = $8,
			original_price_kopecks = $9, pickup_start = $10, pickup_end = $11,
			quantity = $12, address = $13, district = $14, latitude = $15,
			longitude = $16, delivery = $17, status = $18, updated_at = now()
		WHERE id = $1 AND status <> 'deleted'
		RETURNING `+offerColumns,
		offer.ID, offer.Title, offer.Merchant, offer.Category, offer.Description,
		offer.Contents, offer.ImageURL, offer.PriceKopecks, offer.OriginalPriceKopecks,
		offer.PickupStart, offer.PickupEnd, offer.Quantity, offer.Address, offer.District,
		offer.Latitude, offer.Longitude, offer.Delivery, offer.Status))
	if errors.Is(err, sql.ErrNoRows) {
		return Offer{}, ErrNotFound
	}
	return updated, err
}

// DeleteOffer is a soft delete. Existing references can remain valid in the
// database while public queries hide rows whose status is "deleted".
func (s *Store) DeleteOffer(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE offers SET status = 'deleted', updated_at = now()
		WHERE id = $1 AND status <> 'deleted'`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
