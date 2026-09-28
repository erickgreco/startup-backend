-- +goose Up
CREATE TABLE leads(
    id              UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    business_name   VARCHAR(255) NOT NULL,
    phone_number    VARCHAR(20) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    message         TEXT NOT NULL,
    data_quality    VARCHAR(20) NOT NULL DEFAULT 'low'
                    CHECK (data_quality IN ('high', 'medium', 'low')),
    quality_score   FLOAT NOT NULL DEFAULT 0.0
                    CHECK (quality_score >= 0 AND quality_score <= 100),
    status          VARCHAR(20) NOT NULL DEFAULT 'new'
                    CHECK (status IN ('new', 'quoted', 'converted', 'lost', 'merged')),
    source          VARCHAR(20) NOT NULL
                    CHECK(source IN ('landing_page', 'email', 'whatsapp', 'manual')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leads_email ON leads(email);
CREATE INDEX idx_leads_status ON leads(status);

CREATE TRIGGER trg_leads_updated_at
    BEFORE UPDATE ON leads
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_leads_updated_at ON leads;
DROP TABLE IF EXISTS leads;
