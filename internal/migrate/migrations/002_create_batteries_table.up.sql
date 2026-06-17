CREATE TABLE IF NOT EXISTS batteries (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    details JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_batteries_title ON batteries(title);
CREATE INDEX idx_batteries_price ON batteries(price);
CREATE INDEX idx_batteries_details ON batteries USING GIN (details);

COMMENT ON COLUMN batteries.details IS 'JSON с характеристиками: {capacity_ah: 60, polarity: "прямая", terminal_type: "евро", brand: "Varta", ...}';