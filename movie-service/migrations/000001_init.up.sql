-- 000001_init.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE seat_type_enum AS ENUM ('standard', 'premium', 'vip');
CREATE TYPE booking_status_enum AS ENUM ('confirmed', 'cancelled');

CREATE TABLE movies (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    duration_min INT NOT NULL,
    genre        VARCHAR(100),
    poster_url   TEXT,
    release_date DATE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE halls (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    total_seats INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE seats (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hall_id     UUID NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
    row_number  INT NOT NULL,
    seat_number INT NOT NULL,
    seat_type   seat_type_enum NOT NULL DEFAULT 'standard',
    UNIQUE(hall_id, row_number, seat_number)
);

CREATE TABLE showtimes (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_id   UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    hall_id    UUID NOT NULL REFERENCES halls(id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time   TIMESTAMPTZ NOT NULL,
    price      DECIMAL(10,2) NOT NULL
);

CREATE TABLE bookings (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL,
    showtime_id  UUID NOT NULL REFERENCES showtimes(id),
    status       booking_status_enum NOT NULL DEFAULT 'confirmed',
    total_price  DECIMAL(10,2) NOT NULL,
    booked_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ
);

CREATE TABLE booking_seats (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    seat_id    UUID NOT NULL REFERENCES seats(id),
    UNIQUE(booking_id, seat_id)
);

-- Indexes
CREATE INDEX idx_movies_genre ON movies(genre);
CREATE INDEX idx_showtimes_movie_id ON showtimes(movie_id);
CREATE INDEX idx_showtimes_start_time ON showtimes(start_time);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_showtime_id ON bookings(showtime_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_booking_seats_booking_id ON booking_seats(booking_id);
CREATE INDEX idx_booking_seats_seat_id ON booking_seats(seat_id);
