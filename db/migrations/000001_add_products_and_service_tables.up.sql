CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    tax_rate NUMERIC(5, 2) NOT NULL DEFAULT 0.00,
    metadata JSONB,
    created_at TIMESTAMPZ DEFAULT NOW(),
    updated_at TIMESTAMPZ DEFAULT NOW()
)