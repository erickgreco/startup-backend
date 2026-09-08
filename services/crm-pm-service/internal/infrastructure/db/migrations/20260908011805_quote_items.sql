-- +goose Up
CREATE TABLE quote_items(
    id                  UUID PRIMARY KEY,
    quote_id            UUID NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    description         VARCHAR(500) NOT NULL,
    quantity            DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
    unit_price_cents    BIGINT NOT NULL CHECK (unit_price_cents > 0),
    subtotal_cents      BIGINT NOT NULL CHECK (subtotal_cents >= 0),
    sort_order          INTEGER NOT NULL DEFAULT 0,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT quote_items_subtotal_consistency CHECK (
        subtotal_cents = (quantity * unit_price_cents)::BIGINT
    )
);

CREATE INDEX idx_quote_items_quote_id ON quote_items(quote_id);
CREATE INDEX idx_quote_items_sort_order ON quote_items(quote_id, sort_order);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_quote_items_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_quote_items_updated_at
BEFORE UPDATE ON quote_items
FOR EACH ROW
EXECUTE FUNCTION update_quote_items_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_quote_items_updated_at;

DROP FUNCTION IF EXISTS update_quote_items_updated_at();

DROP INDEX IF EXISTS idx_quote_items_sort_order;
DROP INDEX IF EXISTS idx_quote_items_quote_id;

DROP TABLE IF EXISTS quote_items;
