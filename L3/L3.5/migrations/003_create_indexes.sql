CREATE INDEX idx_bookings_event_id   ON bookings(event_id);
CREATE INDEX idx_bookings_user_id    ON bookings(user_id);

CREATE INDEX idx_bookings_status_booked
    ON bookings(status)
    WHERE status = 'booked';

CREATE INDEX idx_bookings_expires_at ON bookings(expires_at);
