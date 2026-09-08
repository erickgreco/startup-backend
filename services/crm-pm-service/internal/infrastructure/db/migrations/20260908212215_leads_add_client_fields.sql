-- +goose Up
ALTER TABLE leads
ADD COLUMN phone_number VARCHAR(20) NOT NULL,
ADD COLUMN business_name VARCHAR(255) NOT NULL;

-- +goose Down
ALTER TABLE leads
DROP COLUMN business_name,
DROP COLUMN phone_number;
