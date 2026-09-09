-- +migrate Up

CREATE TABLE IF NOT EXISTS cleaning (
    id              SERIAL PRIMARY KEY,
    reservation_id  INT REFERENCES reservations(id) ON DELETE CASCADE,
    cleaning_time   TIMESTAMP(0),
    room            VARCHAR(50),
    cleaning_price  INT          DEFAULT 0,
    laundry_price   INT          DEFAULT 0,
    agent_name      VARCHAR(255) DEFAULT '',
    description     TEXT         DEFAULT '',
    paid            BOOLEAN      NOT NULL DEFAULT FALSE
);

-- +migrate Down

DROP TABLE IF EXISTS cleaning;