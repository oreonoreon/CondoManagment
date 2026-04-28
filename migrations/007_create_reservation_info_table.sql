-- +migrate Up

CREATE TABLE IF NOT EXISTS reservation_info (
    id                  SERIAL PRIMARY KEY,
    reservation_id      INT NOT NULL UNIQUE
                            REFERENCES reservations(id) ON DELETE CASCADE,
    deposit             INT         NOT NULL DEFAULT 0,
    deposit_currency    VARCHAR(10)          DEFAULT 'USD',
    prepayment          INT         NOT NULL DEFAULT 0,
    payment_on_checkin  INT         NOT NULL DEFAULT 0,
    notes               TEXT                 DEFAULT ''
);

-- +migrate Down

DROP TABLE IF EXISTS reservation_info;