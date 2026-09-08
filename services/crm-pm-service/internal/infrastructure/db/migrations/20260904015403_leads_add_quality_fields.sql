-- +goose Up
ALTER TABLE leads
ADD COLUMN data_quality VARCHAR(20) NOT NULL DEFAULT 'low'
    CHECK (data_quality IN ('high', 'medium', 'low')),
ADD COLUMN quality_score FLOAT NOT NULL DEFAULT 0.0
    CHECK (quality_score >= 0 AND quality_score <= 100);

CREATE INDEX idx_leads_data_quality ON leads(data_quality);
CREATE INDEX idx_leads_quality_score ON leads(quality_score DESC);

CREATE INDEX idx_leads_quality_status ON leads(data_quality, status);

-- +goose Down
DROP INDEX IF EXISTS idx_leads_quality_status;
DROP INDEX IF EXISTS idx_leads_quality_score;
DROP INDEX IF EXISTS idx_leads_data_quality;

ALTER TABLE leads
DROP COLUMN quality_score,
DROP COLUMN data_quality;
