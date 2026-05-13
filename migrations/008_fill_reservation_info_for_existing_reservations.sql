-- +migrate Up

INSERT INTO reservation_info (reservation_id, deposit, deposit_currency, prepayment, payment_on_checkin, actual_check_in, actual_check_out)
SELECT
    r.id,
    300,
    'USD',
    0,
    r.price,
    r.check_in,
    r.check_out
FROM reservations r
WHERE NOT EXISTS (
    SELECT 1 FROM reservation_info ri WHERE ri.reservation_id = r.id
);

-- +migrate Down

DELETE FROM reservation_info
WHERE reservation_id IN (SELECT id FROM reservations);
