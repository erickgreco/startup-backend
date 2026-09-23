
-- name: CreateLead :one
INSERT INTO leads (
    id, name, business_name, phone_number, email, message,
    data_quality, quality_score, status, source
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetLeadByID :one
SELECT * FROM leads
WHERE id = $1;

-- name: GetLeadsByStatus :many
SELECT * FROM leads
WHERE status = $1
ORDER BY created_at DESC;

-- name: UpdateLeadStatus :one
UPDATE leads
SET status = @status
WHERE id = @lead_id
RETURNING *;

-- name: DeleteLead :execrows
DELETE FROM leads
WHERE id = $1;

-- name: GetLeadsByEmail :many
SELECT * FROM leads
WHERE email = $1
ORDER BY created_at DESC;