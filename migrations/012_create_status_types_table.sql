-- +migrate Up

CREATE TABLE IF NOT EXISTS status_types (
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(100) NOT NULL,
    color       VARCHAR(20)  DEFAULT '',
    sort_order  INT          NOT NULL DEFAULT 0,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE
);

INSERT INTO status_types (code, name, sort_order) VALUES
    ('paid_to_owner',    'Оплачено собственнику',      10),
    ('guest_checked_in', 'Гость заехал',                20),
    ('guest_paid_full',  'Гость оплатил полностью',     30);

-- +migrate Down

DROP TABLE IF EXISTS status_types;
