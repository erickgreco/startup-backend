
-- name: CreateQuote :one
INSERT INTO quotes (
    id, lead_id, quote_number, client_name, client_business, client_email, client_phone,
	project_title, project_description, project_estimated_duration_days,
	payment_method, status, created_at, subtotal, total, currency, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING *;

-- name: CreateQuoteItems :batchone
INSERT INTO quote_items (
    id, quote_id, description, quantity, unit, unit_price, subtotal, sort_order,
	created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetQuoteByID :one
SELECT * FROM quotes
WHERE id = $1;

-- name: ListQuoteItemsByQuoteID :many
SELECT * FROM quote_items
WHERE quote_id = $1
ORDER BY sort_order ASC;