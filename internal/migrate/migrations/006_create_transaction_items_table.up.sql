
CREATE TABLE IF NOT EXISTS transaction_items (
                                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    battery_id UUID NOT NULL REFERENCES batteries(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    price_at_time DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);
CREATE INDEX idx_transaction_items_battery_id ON transaction_items(battery_id);