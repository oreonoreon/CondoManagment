-- +migrate Up


-- Удаляем старые ограничения, завязанные на тип date
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_no_overlap;
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_check;

-- Меняем тип столбцов с date на timestamp(0) (дата + время, точность до секунды)
ALTER TABLE reservations
    ALTER COLUMN check_in  TYPE timestamp(0) USING check_in::timestamp(0),
    ALTER COLUMN check_out TYPE timestamp(0) USING check_out::timestamp(0);

-- Восстанавливаем ограничение check_in < check_out
ALTER TABLE reservations
    ADD CONSTRAINT reservations_check CHECK (check_in < check_out);

-- Восстанавливаем exclusion-ограничение (исключаем пересечение броней),
-- теперь используем tsrange вместо daterange
ALTER TABLE reservations
    ADD CONSTRAINT reservations_no_overlap
        EXCLUDE USING gist (
            room_number WITH =,
            tsrange(check_in::timestamp, check_out::timestamp, '[)') WITH &&
        );

-- +migrate Down

-- Откатываем exclusion-ограничение и check
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_no_overlap;
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_check;

-- Возвращаем тип столбцов обратно в date (время будет отброшено)
ALTER TABLE reservations
    ALTER COLUMN check_in  TYPE date USING check_in::date,
    ALTER COLUMN check_out TYPE date USING check_out::date;

-- Восстанавливаем ограничения под тип date
ALTER TABLE reservations
    ADD CONSTRAINT reservations_check CHECK (check_in < check_out);

ALTER TABLE reservations
    ADD CONSTRAINT reservations_no_overlap
        EXCLUDE USING gist (
            room_number WITH =,
            daterange(check_in, check_out, '[)'::text) WITH &&
        );


