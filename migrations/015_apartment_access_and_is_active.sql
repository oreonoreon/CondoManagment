-- +migrate Up

-- Видимость апартаментов переезжает из массива apartments.access в таблицу apartment_access
-- (выдача доступа роли или конкретному пользователю), неактивность — в apartments.is_active.
-- Миграция только добавляет: старые apartments.access и apartment_owners удаляет 016,
-- которую нужно применять после деплоя бэкенда, работающего с apartment_access.

-- ─── Неактивные апартаменты ───

ALTER TABLE apartments ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Неактивными считались апартаменты, видимые только admin
UPDATE apartments SET is_active = FALSE WHERE NOT ('manager' = ANY(COALESCE(access, '{}')));

-- Апартаменты когда-то загружались с явными id, и счётчик отстал от max(id):
-- без этого создание апартамента упирается в занятый id
SELECT setval('apartments_id_seq', (SELECT COALESCE(max(id), 1) FROM apartments));

-- ─── Доступ к апартаментам ───

-- Пользователь видит апартамент, если доступ выдан его роли или лично ему
-- (неактивные — только при праве apartments.inactive; право apartments.all — все без выдачи).
CREATE TABLE IF NOT EXISTS apartment_access (
    id          SERIAL PRIMARY KEY,
    room_number VARCHAR(255) NOT NULL REFERENCES apartments(room_number) ON UPDATE CASCADE ON DELETE CASCADE,
    role        VARCHAR(50)  REFERENCES roles(name) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id     UUID         REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT apartment_access_role_xor_user CHECK ((role IS NULL) <> (user_id IS NULL))
);

CREATE UNIQUE INDEX IF NOT EXISTS apartment_access_role_key ON apartment_access(room_number, role) WHERE role IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS apartment_access_user_key ON apartment_access(room_number, user_id) WHERE user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS apartment_access_user_id_idx ON apartment_access(user_id) WHERE user_id IS NOT NULL;

-- Роли из apartments.access (admin не нужен — у него право apartments.all)
INSERT INTO apartment_access (room_number, role)
SELECT a.room_number, r.role
FROM apartments a
CROSS JOIN LATERAL unnest(a.access) AS r(role)
JOIN roles ON roles.name = r.role
WHERE a.room_number IS NOT NULL AND r.role <> 'admin'
ON CONFLICT DO NOTHING;

-- Неактивным апартаментам — те же роли, что у активных: после повторной активации
-- они снова будут видны, скрывает их теперь только is_active
INSERT INTO apartment_access (room_number, role)
SELECT a.room_number, r.role
FROM apartments a
CROSS JOIN (SELECT DISTINCT unnest(access) AS role FROM apartments WHERE is_active) AS r
JOIN roles ON roles.name = r.role
WHERE a.room_number IS NOT NULL AND NOT a.is_active AND r.role <> 'admin'
ON CONFLICT DO NOTHING;

-- Владельцы — личная выдача доступа
INSERT INTO apartment_access (room_number, user_id)
SELECT room_number, user_id FROM apartment_owners
ON CONFLICT DO NOTHING;

-- ─── Права ───
--   apartments.all      — все апартаменты без выдачи доступа
--   apartments.inactive — видеть и неактивные апартаменты
--   apartments.manage   — админка: апартаменты и доступ к ним
--   roles.manage        — админка: роли и их права
INSERT INTO role_permissions (role, permission) VALUES
    ('admin', 'apartments.all'),
    ('admin', 'apartments.inactive'),
    ('admin', 'apartments.manage'),
    ('admin', 'roles.manage')
ON CONFLICT DO NOTHING;

-- ─── Пользователи ───

-- Вход по имени: имя должно быть уникальным
CREATE UNIQUE INDEX IF NOT EXISTS users_name_key ON users(name);

-- Удаление пользователя не должно ломаться об историю статусов броней
ALTER TABLE reservation_statuses DROP CONSTRAINT IF EXISTS reservation_statuses_set_by_fkey;
ALTER TABLE reservation_statuses
    ADD CONSTRAINT reservation_statuses_set_by_fkey FOREIGN KEY (set_by) REFERENCES users(id) ON DELETE SET NULL;

-- +migrate Down

ALTER TABLE reservation_statuses DROP CONSTRAINT IF EXISTS reservation_statuses_set_by_fkey;
ALTER TABLE reservation_statuses
    ADD CONSTRAINT reservation_statuses_set_by_fkey FOREIGN KEY (set_by) REFERENCES users(id);

DROP INDEX IF EXISTS users_name_key;

DELETE FROM role_permissions
WHERE role = 'admin' AND permission IN ('apartments.all', 'apartments.inactive', 'apartments.manage', 'roles.manage');

DROP TABLE IF EXISTS apartment_access;

ALTER TABLE apartments DROP COLUMN IF EXISTS is_active;
