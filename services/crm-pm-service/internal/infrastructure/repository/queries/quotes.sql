
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

-- name: UpdateQuote :one
UPDATE quotes
SET
    client_name                     = sqlc.arg(client_name),
    client_business                 = sqlc.narg(client_business),
    client_email                    = sqlc.arg(client_email),
    client_phone                    = sqlc.arg(client_phone),
    project_title                   = sqlc.arg(project_title),
    project_description             = sqlc.arg(project_description),
    project_resume                  = sqlc.narg(project_resume),
    project_estimated_duration_days = sqlc.arg(project_estimated_duration_days),
    functionality_example           = sqlc.narg(functionality_example),
    payment_method                  = sqlc.arg(payment_method),
    subtotal                        = sqlc.arg(subtotal),
    total                           = sqlc.arg(total),
    currency                        = sqlc.arg(currency),
    notes                           = sqlc.narg(notes)
WHERE id = sqlc.arg(id) AND status = 'draft'
RETURNING *;

-- name: DeleteQuoteItemsByQuoteID :execrows
DELETE FROM quote_items
WHERE quote_id = $1;