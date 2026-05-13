-- +migrate Up

-- Создаём записи в cleaning для всех резерваций, у которых ещё нет связанной записи.
-- Правила маппинга cleaning_price / laundry_price:
--   reservations.cleaning_price = 1500 → cleaning_price=1000, laundry_price=500
--   reservations.cleaning_price = 2500 → cleaning_price=1500, laundry_price=1000
--   иначе                               → cleaning_price=reservations.cleaning_price, laundry_price=0

INSERT INTO cleaning (reservation_id, cleaning_time, room, cleaning_price, laundry_price, agent_name, description, paid)
SELECT
    r.id,
    r.check_out,
    r.room_number,
    CASE
        WHEN r.cleaning_price = 1500 THEN 1000
        WHEN r.cleaning_price = 2500 THEN 1500
        ELSE r.cleaning_price
    END,
    CASE
        WHEN r.cleaning_price = 1500 THEN 500
        WHEN r.cleaning_price = 2500 THEN 1000
        ELSE 0
    END,
    'Our Apartment',
    '',
    TRUE
FROM reservations r
WHERE NOT EXISTS (
    SELECT 1 FROM cleaning c WHERE c.reservation_id = r.id
);

-- +migrate Down

DELETE FROM cleaning
WHERE reservation_id IN (SELECT id FROM reservations);
