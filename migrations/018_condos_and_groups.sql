-- +migrate Up

-- Апартаменты делятся на кондо, кондо — на группы. Названия задаются в админке.
-- Апартамент может быть без кондо, а в кондо — без группы.
-- Миграция только добавляет, старый код её не замечает: применять ДО деплоя бэкенда.

CREATE TABLE condos (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE apartment_groups (
    id       SERIAL PRIMARY KEY,
    condo_id INTEGER NOT NULL REFERENCES condos(id) ON DELETE CASCADE,
    name     VARCHAR(100) NOT NULL,
    UNIQUE (condo_id, name),
    UNIQUE (condo_id, id) -- цель составного внешнего ключа из apartments
);

ALTER TABLE apartments
    ADD COLUMN condo_id INTEGER REFERENCES condos(id),
    ADD COLUMN group_id INTEGER,
    -- группа должна принадлежать кондо апартамента
    ADD CONSTRAINT apartments_group_in_condo FOREIGN KEY (condo_id, group_id) REFERENCES apartment_groups(condo_id, id),
    ADD CONSTRAINT apartments_group_needs_condo CHECK (group_id IS NULL OR condo_id IS NOT NULL);

CREATE INDEX apartments_condo_id_idx ON apartments(condo_id);
CREATE INDEX apartments_group_id_idx ON apartments(group_id);

-- Serenity раньше определялся по префиксу номера и шёл в календаре отдельным блоком
WITH serenity AS (
    INSERT INTO condos (name)
    SELECT 'Serenity' WHERE EXISTS (SELECT 1 FROM apartments WHERE room_number LIKE 'Serenity %')
    RETURNING id
)
UPDATE apartments SET condo_id = serenity.id
FROM serenity
WHERE room_number LIKE 'Serenity %';

-- +migrate Down

ALTER TABLE apartments
    DROP CONSTRAINT apartments_group_needs_condo,
    DROP CONSTRAINT apartments_group_in_condo,
    DROP COLUMN group_id,
    DROP COLUMN condo_id;

DROP TABLE apartment_groups;
DROP TABLE condos;
