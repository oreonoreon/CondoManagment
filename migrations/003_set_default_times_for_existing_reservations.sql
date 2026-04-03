-- +migrate Up

-- Проставляем время 13:00:00 для check_in у всех записей, где время равно 00:00:00
UPDATE reservations
SET check_in = date_trunc('day', check_in) + INTERVAL '13 hours'
WHERE check_in::time = '00:00:00';

-- Проставляем время 11:00:00 для check_out у всех записей, где время равно 00:00:00
UPDATE reservations
SET check_out = date_trunc('day', check_out) + INTERVAL '11 hours'
WHERE check_out::time = '00:00:00';

-- +migrate Down

-- Откат: сбрасываем время обратно в 00:00:00 для записей, которые мы изменили
UPDATE reservations
SET check_in = date_trunc('day', check_in)
WHERE check_in::time = '13:00:00';

UPDATE reservations
SET check_out = date_trunc('day', check_out)
WHERE check_out::time = '11:00:00';
