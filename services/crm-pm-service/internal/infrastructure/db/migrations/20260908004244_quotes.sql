-- +goose Up
CREATE TABLE quotes(
    id                  UUID PRIMARY KEY,
    lead_id             UUID NOT NULL REFERENCES leads(id) ON DELETE RESTRICT,
    quote_number        VARCHAR(50) NOT NULL UNIQUE,
    status              VARCHAR(30) NOT NULL DEFAULT 'draft' 
                        CHECK (status IN ('draft', 'sent', 'accepted', 'rejected', 'expired', 'converted')),
    issued_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    valid_until         TIMESTAMP WITH TIME ZONE NOT NULL,
    sent_at             TIMESTAMP WITH TIME ZONE,
    accepted_at         TIMESTAMP WITH TIME ZONE,
    rejected_at         TIMESTAMP WITH TIME ZONE,
    subtotal_cents      BIGINT NOT NULL CHECK (subtotal_cents >= 0),
    tax_amount_cents    BIGINT NOT NULL CHECK (tax_amount_cents >= 0),
    total_cents         BIGINT NOT NULL CHECK (total_cents >= 0),
    currency            VARCHAR(3) NOT NULL DEFAULT 'MXN',
    tax_rate_percent    DECIMAL(5, 2) DEFAULT 0.00 CHECK (tax_rate_percent >= 0),
    notes               TEXT,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT quotes_valid_dates CHECK (valid_until > issued_at),
    CONSTRAINT quotes_total_consistency CHECK (
        total_cents = subtotal_cents + tax_amount_cents
    )
);

CREATE INDEX idx_quotes_lead_id ON quotes(lead_id);
CREATE INDEX idx_quotes_status ON quotes(status);
CREATE INDEX idx_quotes_created_at ON quotes(created_at DESC);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_quotes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_quotes_updated_at
BEFORE UPDATE ON quotes
FOR EACH ROW
EXECUTE FUNCTION update_quotes_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_quotes_updated_at ON quotes;

DROP FUNCTION IF EXISTS update_quotes_updated_at();

DROP INDEX IF EXISTS idx_quotes_created_at;
DROP INDEX IF EXISTS idx_quotes_status;
DROP INDEX IF EXISTS idx_quotes_lead_id;

DROP TABLE IF EXISTS quotes;