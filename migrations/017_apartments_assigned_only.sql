-- +migrate Up

-- apartments.assigned_only — роль видит только апартаменты, выданные пользователю лично;
-- выдача апартамента самой роли для неё не действует и админкой запрещена.
-- Нужно владельцам: каждому только свои квартиры, а не все квартиры, выданные роли apartmentOwner.
-- Требует миграцию 015 (apartment_access). Старый код право игнорирует, применять можно до деплоя.

INSERT INTO role_permissions (role, permission)
SELECT name, 'apartments.assigned_only' FROM roles WHERE name = 'apartmentOwner'
ON CONFLICT DO NOTHING;

-- Выдачи апартаментов ролям с этим правом больше не действуют — убираем, чтобы не путали в админке
DELETE FROM apartment_access aa
USING role_permissions rp
WHERE rp.role = aa.role AND rp.permission = 'apartments.assigned_only';

-- +migrate Down

DELETE FROM role_permissions WHERE permission = 'apartments.assigned_only';
