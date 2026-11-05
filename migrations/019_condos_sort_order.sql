-- +migrate Up

-- Порядок кондо в календаре задаётся в админке. Сначала — по алфавиту, как было.
-- Требует миграцию 018. Применять ДО деплоя бэкенда, который читает sort_order.

ALTER TABLE condos ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;

UPDATE condos c SET sort_order = o.pos
FROM (SELECT id, row_number() OVER (ORDER BY name) AS pos FROM condos) o
WHERE c.id = o.id;

-- +migrate Down

ALTER TABLE condos DROP COLUMN sort_order;
