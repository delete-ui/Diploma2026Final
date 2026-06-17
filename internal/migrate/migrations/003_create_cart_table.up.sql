CREATE TABLE IF NOT EXISTS cart (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    battery_id UUID NOT NULL REFERENCES batteries(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, battery_id)
    );

CREATE INDEX idx_cart_user_id ON cart(user_id);
CREATE INDEX idx_cart_battery_id ON cart(battery_id);