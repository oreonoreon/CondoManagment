-- +migrate Up

-- Базовая схема БД: снимок локальной и продовой (Railway) базы на 2026-10-22,
-- до миграции 014_roles_and_permissions. Включает всё, что дали миграции 002–013,
-- и таблицы, которые раньше нигде в репозитории не описывались.
--
-- Новая пустая база: применить этот файл, затем миграции начиная с 014.
-- Существующие базы (локальная, Railway): этот файл НЕ применять — всё уже есть.
--
-- Отличие от файлов 002–013: триггера из 009 (reservation_info_set_actual_times)
-- в рабочих базах нет, поэтому здесь его тоже нет.

CREATE EXTENSION IF NOT EXISTS btree_gist WITH SCHEMA public; -- EXCLUDE-ограничение reservations_no_overlap
CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;   -- gen_random_uuid()

-- ─── Апартаменты, гости, брони ───

CREATE TABLE apartments (
    id              SERIAL PRIMARY KEY,
    room_number     VARCHAR(255) CONSTRAINT apartments_room_number_key UNIQUE,
    description     VARCHAR(255),
    airbnb_calendar VARCHAR(255),
    access          VARCHAR[] DEFAULT ARRAY['admin'::text, 'manager'::text] -- роли, которым виден апартамент
);

CREATE TABLE guests (
    guest_id    UUID DEFAULT public.gen_random_uuid() PRIMARY KEY,
    name        VARCHAR(255),
    phone       VARCHAR(255) DEFAULT NULL CONSTRAINT guests_pk UNIQUE,
    description VARCHAR(255)
);

CREATE TABLE reservations (
    id                            SERIAL PRIMARY KEY,
    room_number                   VARCHAR(255) NOT NULL REFERENCES apartments(room_number),
    guest_id                      UUID REFERENCES guests(guest_id),
    check_in                      TIMESTAMP(0),
    check_out                     TIMESTAMP(0),
    price                         INTEGER,
    cleaning_price                INTEGER,
    electricity_and_water_payment VARCHAR(255),
    adult                         INTEGER,
    children                      INTEGER,
    description                   TEXT,
    days                          INTEGER,
    price_for_night               INTEGER,
    CONSTRAINT reservations_check CHECK (check_in < check_out),
    -- брони одного апартамента не пересекаются
    CONSTRAINT reservations_no_overlap EXCLUDE USING gist (
        room_number WITH =,
        tsrange(check_in::timestamp, check_out::timestamp, '[)') WITH &&
    )
);

CREATE TABLE reservation_info (
    id                 SERIAL PRIMARY KEY,
    reservation_id     INTEGER NOT NULL UNIQUE REFERENCES reservations(id) ON DELETE CASCADE,
    deposit            INTEGER     NOT NULL DEFAULT 0,
    deposit_currency   VARCHAR(10)          DEFAULT 'USD',
    prepayment         INTEGER     NOT NULL DEFAULT 0,
    payment_on_checkin INTEGER     NOT NULL DEFAULT 0,
    actual_check_in    TIMESTAMP(0) NOT NULL,
    actual_check_out   TIMESTAMP(0) NOT NULL
);

CREATE TABLE cleaning (
    id             SERIAL PRIMARY KEY,
    reservation_id INTEGER REFERENCES reservations(id) ON DELETE CASCADE,
    cleaning_time  TIMESTAMP(0),
    room           VARCHAR(50), -- не ссылается на apartments: бывают уборки сторонних комнат
    cleaning_price INTEGER      DEFAULT 0,
    laundry_price  INTEGER      DEFAULT 0,
    agent_name     VARCHAR(255) DEFAULT '',
    description    TEXT         DEFAULT '',
    paid           BOOLEAN NOT NULL DEFAULT FALSE
);

-- ─── Пользователи и статусы броней ───

CREATE TABLE users (
    id           UUID PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    passwordhash VARCHAR(255) NOT NULL,
    phone        VARCHAR(255) NOT NULL UNIQUE,
    role         VARCHAR(255) NOT NULL
);

CREATE TABLE status_types (
    id         SERIAL PRIMARY KEY,
    code       VARCHAR(50)  NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    color      VARCHAR(20)           DEFAULT '',
    sort_order INTEGER      NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE TABLE reservation_statuses (
    id             SERIAL PRIMARY KEY,
    reservation_id INTEGER NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    status_type_id INTEGER NOT NULL REFERENCES status_types(id),
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    set_by         UUID REFERENCES users(id),
    set_at         TIMESTAMP(0) NOT NULL DEFAULT now(),
    removed_at     TIMESTAMP(0)
);

CREATE INDEX idx_reservation_statuses_reservation_id ON reservation_statuses(reservation_id);
CREATE INDEX idx_reservation_statuses_active ON reservation_statuses(reservation_id, status_type_id) WHERE is_active;

-- Справочник статусов в том виде, в каком он сейчас в рабочих базах
INSERT INTO status_types (id, code, name, color, sort_order, is_active) VALUES
    (1, 'paid_to_owner',    'Оплачено собственнику',   '', 2, TRUE),
    (2, 'guest_checked_in', 'Гость заехал',            '', 1, TRUE),
    (3, 'guest_paid_full',  'Гость оплатил полностью', '', 3, FALSE);
SELECT setval('status_types_id_seq', (SELECT max(id) FROM status_types));

-- ─── Служебные таблицы ───

-- Сессии (github.com/antonlindstrom/pgstore; создаёт таблицу и сама, если её нет)
CREATE TABLE http_sessions (
    id          BIGSERIAL PRIMARY KEY,
    key         BYTEA,
    data        BYTEA,
    created_on  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_on TIMESTAMPTZ,
    expires_on  TIMESTAMPTZ
);

CREATE INDEX http_sessions_expiry_idx ON http_sessions(expires_on);
CREATE INDEX http_sessions_key_idx ON http_sessions(key);

-- Настройки чтения Excel-календаря (пишутся приложением при старте из etc/config.yaml)
CREATE TABLE settings (
    excel_path               VARCHAR(255),
    sheet_name               VARCHAR(255) UNIQUE,
    sheet_startcolumnum      INTEGER,
    sheet_yearandmonthrownum INTEGER,
    sheet_daysrownum         INTEGER,
    sheet_startrow           INTEGER,
    sheet_apartmentcolumnum  INTEGER
);

-- ─── Данные парсера Airbnb ───

CREATE TABLE airbnb_host (
    id        VARCHAR(255) PRIMARY KEY,
    name      VARCHAR(255),
    trustable BOOLEAN
);

CREATE TABLE airbnb_room (
    roomid             BIGINT PRIMARY KEY,
    badges             VARCHAR(255)[],
    name               VARCHAR(255),
    title              VARCHAR(255),
    type               VARCHAR(255),
    kind               VARCHAR(255),
    category           VARCHAR(255),
    images             VARCHAR(255)[],
    understandabletype VARCHAR(255),
    hostid             VARCHAR(255) REFERENCES airbnb_host(id)
);

CREATE TABLE airbnb_price (
    roomid        BIGINT REFERENCES airbnb_room(roomid),
    price         DOUBLE PRECISION,
    check_in      DATE,
    check_out     DATE,
    scraping_date DATE
);

CREATE TABLE coordinates_airbnb_room (
    roomid       BIGINT REFERENCES airbnb_room(roomid),
    latitude     DOUBLE PRECISION,
    longitude    DOUBLE PRECISION,
    locationname VARCHAR(255)
);

CREATE TABLE rating_airbnb_room (
    roomid      BIGINT REFERENCES airbnb_room(roomid),
    value       DOUBLE PRECISION,
    reviewcount INTEGER
);

-- +migrate Down

DROP TABLE IF EXISTS rating_airbnb_room;
DROP TABLE IF EXISTS coordinates_airbnb_room;
DROP TABLE IF EXISTS airbnb_price;
DROP TABLE IF EXISTS airbnb_room;
DROP TABLE IF EXISTS airbnb_host;
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS http_sessions;
DROP TABLE IF EXISTS reservation_statuses;
DROP TABLE IF EXISTS status_types;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS cleaning;
DROP TABLE IF EXISTS reservation_info;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS guests;
DROP TABLE IF EXISTS apartments;
