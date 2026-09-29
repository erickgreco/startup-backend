-- +goose Up
CREATE TABLE user_profile (
    id                  UUID PRIMARY KEY,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    phone_number        VARCHAR(20) NOT NULL,
    profile_image_url   TEXT,
    business_name       VARCHAR(255),
    industry            VARCHAR(100) NOT NULL,
    logo_url            TEXT,
    default_currency    VARCHAR(3) NOT NULL,
    timezone            VARCHAR(50) NOT NULL,
    language            VARCHAR(5) NOT NULL,
    website_url         TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_user_profile_updated_at
    BEFORE UPDATE ON user_profile
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_user_profile_updated_at ON user_profile;
DROP TABLE IF EXISTS user_profile;
