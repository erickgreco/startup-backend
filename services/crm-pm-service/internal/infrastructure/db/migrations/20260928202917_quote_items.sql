-- +goose Up
CREATE TABLE quote_items(
    id                  UUID PRIMARY KEY,
    quote_id            UUID NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    description         VARCHAR(500) NOT NULL,
    quantity            DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
    unit_price          BIGINT NOT NULL CHECK (unit_price >= 0),
    unit                VARCHAR(50) NOT NULL DEFAULT 'unidad',
    subtotal            BIGINT NOT NULL CHECK (subtotal >= 0),
    sort_order          INTEGER NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT quote_items_subtotal_consistency CHECK (
        subtotal = (quantity * unit_price)::BIGINT
    )
);

CREATE INDEX idx_quote_items_quote_id ON quote_items(quote_id);
CREATE INDEX idx_quote_items_sort_order ON quote_items(quote_id, sort_order);

CREATE TRIGGER trg_quote_items_updated_at
    BEFORE UPDATE ON quote_items
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_quote_items_updated_at ON quote_items;
DROP TABLE IF EXISTS quote_items;