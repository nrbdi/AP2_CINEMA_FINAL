-- 000001_init.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE payment_status_enum AS ENUM ('pending', 'paid', 'refunded');

CREATE TABLE payments (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL UNIQUE,
    user_id    UUID NOT NULL,
    amount     DECIMAL(10,2) NOT NULL,
    status     payment_status_enum NOT NULL DEFAULT 'pending',
    method     VARCHAR(50) NOT NULL DEFAULT 'system',
    paid_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
