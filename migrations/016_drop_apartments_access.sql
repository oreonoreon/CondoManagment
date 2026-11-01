-- +migrate Up

-- Удаляет старый механизм видимости апартаментов, заменённый apartment_access (миграция 015).
-- Применять ТОЛЬКО после деплоя бэкенда, работающего с apartment_access:
-- предыдущая версия читает apartments.access и apartment_owners.

ALTER TABLE apartments DROP COLUMN IF EXISTS access;

DROP TABLE IF EXISTS apartment_owners;

DELETE FROM role_permissions WHERE permission = 'apartments.owned_only';

-- +migrate Down

INSERT INTO role_permissions (role, permission) VALUES ('apartmentOwner', 'apartments.owned_only')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS apartment_owners (
    room_number VARCHAR(255) NOT NULL REFERENCES apartments(room_number) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (room_number, user_id)
);
CREATE INDEX IF NOT EXISTS idx_apartment_owners_user_id ON apartment_owners(user_id);

INSERT INTO apartment_owners (room_number, user_id)
SELECT aa.room_number, aa.user_id
FROM apartment_access aa
JOIN users u ON u.id = aa.user_id AND u.role = 'apartmentOwner'
ON CONFLICT DO NOTHING;

ALTER TABLE apartments ADD COLUMN IF NOT EXISTS access VARCHAR[] DEFAULT ARRAY['admin', 'manager', 'agent']::varchar[];

UPDATE apartments a
SET access = ARRAY['admin']::varchar[] || COALESCE(
    (SELECT array_agg(aa.role ORDER BY aa.role) FROM apartment_access aa
     WHERE aa.room_number = a.room_number AND aa.role IS NOT NULL AND a.is_active),
    '{}'::varchar[]);
