-- +goose Up
CREATE TABLE quotes(
    id                              UUID PRIMARY KEY,
    lead_id                         UUID REFERENCES leads(id) ON DELETE RESTRICT,
    quote_number                    VARCHAR(50) NOT NULL UNIQUE,
    client_name                     VARCHAR(255) NOT NULL,
    client_business                 VARCHAR(255),
    client_email                    VARCHAR(255) NOT NULL,
    client_phone                    VARCHAR(20) NOT NULL,
    project_title                   VARCHAR(500) NOT NULL,
    project_description             TEXT NOT NULL,
    project_resume                  TEXT,
    project_estimated_duration_days INT NOT NULL,
    functionality_example           TEXT,
    payment_method                  VARCHAR(50) NOT NULL
                                    CHECK (payment_method IN ('upfront', 'milestone', 'upon_delivery', 'split')),
    status                          VARCHAR(30) NOT NULL DEFAULT 'draft'
                                    CHECK (status IN ('draft', 'sent', 'accepted', 'rejected', 'expired', 'converted')),
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until                     TIMESTAMPTZ,
    sent_at                         TIMESTAMPTZ,
    accepted_at                     TIMESTAMPTZ,
    rejected_at                     TIMESTAMPTZ,
    subtotal                        BIGINT NOT NULL CHECK (subtotal >= 0),
    total                           BIGINT NOT NULL CHECK (total >= 0),
    currency                        VARCHAR(3) NOT NULL
                                    CHECK (currency ~ '^[A-Z]{3}$'),
    notes                           TEXT,
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT quotes_valid_dates CHECK (
        valid_until IS NULL OR sent_at IS NULL OR valid_until > sent_at
    ),
    CONSTRAINT quotes_total_consistency CHECK (
        total >= subtotal
    )
);

CREATE INDEX idx_quotes_lead_id ON quotes(lead_id);
CREATE INDEX idx_quotes_client_email ON quotes(client_email);
CREATE INDEX idx_quotes_status ON quotes(status);
CREATE INDEX idx_quotes_created_at ON quotes(created_at DESC);
CREATE INDEX idx_quotes_payment_method ON quotes(payment_method);

CREATE TRIGGER trg_quotes_updated_at
    BEFORE UPDATE ON quotes
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_quotes_updated_at ON quotes;
DROP TABLE IF EXISTS quotes;
