-- +migrate Up

CREATE TABLE IF NOT EXISTS reservation_statuses (
    id              SERIAL PRIMARY KEY,
    reservation_id  INT NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    status_type_id  INT NOT NULL REFERENCES status_types(id),
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    set_by          UUID REFERENCES users(id),
    set_at          TIMESTAMP(0) NOT NULL DEFAULT now(),
    removed_at      TIMESTAMP(0)
);

CREATE INDEX idx_reservation_statuses_reservation_id ON reservation_statuses(reservation_id);
CREATE INDEX idx_reservation_statuses_active ON reservation_statuses(reservation_id, status_type_id) WHERE is_active;

-- +migrate Down

DROP TABLE IF EXISTS reservation_statuses;
