-- +migrate Up

-- Триггер: если actual_check_in / actual_check_out не переданы при INSERT,
-- автоматически подставляет check_in / check_out из связанной записи reservations.

CREATE OR REPLACE FUNCTION reservation_info_set_actual_times()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.actual_check_in IS NULL THEN
        SELECT check_in  INTO NEW.actual_check_in  FROM reservations WHERE id = NEW.reservation_id;
    END IF;
    IF NEW.actual_check_out IS NULL THEN
        SELECT check_out INTO NEW.actual_check_out FROM reservations WHERE id = NEW.reservation_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_reservation_info_actual_times
    BEFORE INSERT ON reservation_info
    FOR EACH ROW EXECUTE FUNCTION reservation_info_set_actual_times();

-- +migrate Down

DROP TRIGGER  IF EXISTS trg_reservation_info_actual_times ON reservation_info;
DROP FUNCTION IF EXISTS reservation_info_set_actual_times();
