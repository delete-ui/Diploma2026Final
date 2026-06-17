
ALTER TABLE batteries
DROP COLUMN IF EXISTS img,
    DROP COLUMN IF EXISTS brand,
    DROP COLUMN IF EXISTS voltage,
    DROP COLUMN IF EXISTS polarity,
    DROP COLUMN IF EXISTS capacity,
    DROP COLUMN IF EXISTS standart,
    DROP COLUMN IF EXISTS technology,
    DROP COLUMN IF EXISTS size_type;

CREATE INDEX IF NOT EXISTS idx_batteries_title ON batteries(title);
CREATE INDEX IF NOT EXISTS idx_batteries_price ON batteries(price);