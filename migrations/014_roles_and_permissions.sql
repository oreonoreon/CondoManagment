-- +migrate Up

-- Справочник ролей. users.role ссылается на него (NOT VALID — существующие строки не перепроверяются).
CREATE TABLE IF NOT EXISTS roles (
    name        VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL DEFAULT ''
);

-- Права ролей. Виды прав:
--   calendar.read / calendar.write  — просмотр / изменение броней, статусов, reservation_info
--   today.read                      — вкладка Today (заезды и выезды по дате)
--   cleaning.manage                 — уборки (без ограничения по апартаментам)
--   reports.read                    — отчёты
--   users.manage                    — создание пользователей
--   tools.use                       — Excel-синхронизация и парсер Airbnb
--   apartments.owned_only           — видны только апартаменты из apartment_owners этого пользователя
--                                     (без него видимость определяет apartments.access)
--   booking.field.<поле>            — какие поля брони отдаются роли; booking.field.all — все поля.
--     Поля: room_number, guest_id, guest_name, guest_phone, guest_description, check_in, check_out,
--     price, cleaning_price, electricity_and_water_payment, adult, children, description, days,
--     price_for_night, reservation_info, statuses. id брони отдаётся всегда.
CREATE TABLE IF NOT EXISTS role_permissions (
    role       VARCHAR(50)  NOT NULL REFERENCES roles(name) ON UPDATE CASCADE ON DELETE CASCADE,
    permission VARCHAR(100) NOT NULL,
    PRIMARY KEY (role, permission)
);

-- Владельцы апартаментов (для роли apartmentOwner). У апартамента может быть несколько владельцев.
CREATE TABLE IF NOT EXISTS apartment_owners (
    room_number VARCHAR(255) NOT NULL REFERENCES apartments(room_number) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (room_number, user_id)
);

CREATE INDEX IF NOT EXISTS idx_apartment_owners_user_id ON apartment_owners(user_id);

INSERT INTO roles (name, description) VALUES
    ('admin',          'Полный доступ'),
    ('manager',        'Работа с бронями, уборками и отчётами'),
    ('agent',          'Календарь только на чтение, ограниченный набор полей'),
    ('apartmentOwner', 'Календарь только на чтение по своим апартаментам')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role, permission) VALUES
    ('admin', 'calendar.read'),
    ('admin', 'calendar.write'),
    ('admin', 'today.read'),
    ('admin', 'cleaning.manage'),
    ('admin', 'reports.read'),
    ('admin', 'users.manage'),
    ('admin', 'tools.use'),
    ('admin', 'booking.field.all'),

    ('manager', 'calendar.read'),
    ('manager', 'calendar.write'),
    ('manager', 'today.read'),
    ('manager', 'cleaning.manage'),
    ('manager', 'reports.read'),
    ('manager', 'tools.use'),
    ('manager', 'booking.field.all'),

    ('agent', 'calendar.read'),
    ('agent', 'booking.field.room_number'),
    ('agent', 'booking.field.guest_name'),
    ('agent', 'booking.field.check_in'),
    ('agent', 'booking.field.check_out'),
    ('agent', 'booking.field.adult'),
    ('agent', 'booking.field.children'),
    ('agent', 'booking.field.days'),

    ('apartmentOwner', 'calendar.read'),
    ('apartmentOwner', 'apartments.owned_only'),
    ('apartmentOwner', 'booking.field.guest_name'),
    ('apartmentOwner', 'booking.field.guest_phone'),
    ('apartmentOwner', 'booking.field.room_number'),
    ('apartmentOwner', 'booking.field.check_in'),
    ('apartmentOwner', 'booking.field.check_out'),
    ('apartmentOwner', 'booking.field.price'),
    ('apartmentOwner', 'booking.field.adult'),
    ('apartmentOwner', 'booking.field.children')
ON CONFLICT DO NOTHING;

-- Неизвестные роли у существующих пользователей тоже попадают в справочник (без прав),
-- чтобы внешний ключ не мешал их дальнейшему редактированию.
INSERT INTO roles (name)
SELECT DISTINCT role FROM users
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
    ADD CONSTRAINT users_role_fkey FOREIGN KEY (role) REFERENCES roles(name) ON UPDATE CASCADE NOT VALID;

-- agent видит те же апартаменты, что и manager
UPDATE apartments
SET access = array_append(access, 'agent')
WHERE 'manager' = ANY(access) AND NOT ('agent' = ANY(access));

ALTER TABLE apartments ALTER COLUMN access SET DEFAULT ARRAY['admin', 'manager', 'agent']::varchar[];

-- +migrate Down

ALTER TABLE apartments ALTER COLUMN access SET DEFAULT ARRAY['admin', 'manager']::varchar[];
UPDATE apartments SET access = array_remove(access, 'agent');

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_fkey;

DROP TABLE IF EXISTS apartment_owners;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
