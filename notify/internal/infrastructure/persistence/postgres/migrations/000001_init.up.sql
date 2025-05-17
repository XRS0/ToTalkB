CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    payload BYTEA NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
); 