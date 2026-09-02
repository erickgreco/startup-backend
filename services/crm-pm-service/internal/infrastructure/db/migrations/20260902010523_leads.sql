-- +goose Up
CREATE TABLE leads (
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL, 
    message     TEXT NOT NULL, 
    status      VARCHAR(20) NOT NULL DEFAULT 'new'
                CHECK (status IN ('new', 'quoted', 'converted', 'lost', 'merged')),
    source      VARCHAR(20) NOT NULL
                CHECK (source IN ('landing_page', 'email', 'whatsapp', 'manual')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_leads_email ON leads (email);

CREATE INDEX idx_leads_status ON leads (status);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_leads_updated_at
BEFORE UPDATE ON leads
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


-- +goose Down
DROP TRIGGER IF EXISTS trg_leads_updated_at ON leads;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS leads;
