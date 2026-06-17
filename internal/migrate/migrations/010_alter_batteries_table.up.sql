
ALTER TABLE batteries DROP COLUMN IF EXISTS details;

ALTER TABLE batteries
    ADD COLUMN IF NOT EXISTS img VARCHAR(500),
    ADD COLUMN IF NOT EXISTS brand VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS voltage INT NOT NULL DEFAULT 12,
    ADD COLUMN IF NOT EXISTS polarity VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS capacity DECIMAL(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS standart VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS technology VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS size_type VARCHAR(50) NOT NULL DEFAULT '';

DROP INDEX IF EXISTS idx_batteries_details;
DROP INDEX IF EXISTS idx_batteries_title;
DROP INDEX IF EXISTS idx_batteries_price;

CREATE INDEX IF NOT EXISTS idx_batteries_brand ON batteries(brand);
CREATE INDEX IF NOT EXISTS idx_batteries_voltage ON batteries(voltage);
CREATE INDEX IF NOT EXISTS idx_batteries_capacity ON batteries(capacity);
CREATE INDEX IF NOT EXISTS idx_batteries_technology ON batteries(technology);