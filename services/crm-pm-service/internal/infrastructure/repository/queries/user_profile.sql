
-- name: CreateUserProfile :one
INSERT INTO user_profile (
    id, first_name, last_name, email, phone_number, profile_image_url, business_name, industry,
    logo_url, default_currency, timezone, language, website_url, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: GetUserProfileByID :one
SELECT * FROM user_profile
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE user_profile
SET 
    first_name          = sqlc.arg(first_name),
    last_name           = sqlc.arg(last_name),
    email               = sqlc.arg(email),
    phone_number        = sqlc.arg(phone_number),
    profile_image_url   = sqlc.narg(profile_image_url),
    business_name       = sqlc.narg(business_name),
    industry            = sqlc.arg(industry),
    logo_url            = sqlc.narg(logo_url),
    default_currency    = sqlc.arg(default_currency),
    language            = sqlc.arg(language),
    website_url         = sqlc.narg(website_url)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUserProfile :execrows
DELETE FROM user_profile
WHERE id = $1;
