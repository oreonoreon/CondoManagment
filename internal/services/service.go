package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"awesomeProject/internal/excelCalendarScraper/models"
	"awesomeProject/internal/excelCalendarScraper/report"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	storageReservation     StorageReservation
	storageGuest           StorageGuest
	storageCleaning        StorageCleaning
	serviceReservationInfo ServiceReservationInfo
}

type StorageReservation interface {
	UpdateReservation(ctx context.Context, r entities.Reservation) (*entities.Reservation, error)
	Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error)
	ReadALLByRoomNumber(ctx context.Context, roomNumber string) ([]entities.Reservation, error)
	ReadWithRoomNumber(ctx context.Context, roomNumber string, checkin, checkout time.Time) ([]entities.Reservation, error)
	FindBookingByGuestUUID(ctx context.Context, uuid uuid.UUID) ([]entities.Reservation, error)
	Delete(ctx context.Context, id int) (*entities.Reservation, error)
	GetReservationByID(ctx context.Context, id int) (*entities.Reservation, error)
	GetReservationsByCheckIn(ctx context.Context, date time.Time) ([]entities.Reservation, error)
	GetReservationsByCheckOut(ctx context.Context, date time.Time) ([]entities.Reservation, error)
}

type StorageGuest interface {
	UpdateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error)
	CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error)
	FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error)
	ReadGuest(ctx context.Context, guestID uuid.UUID) (*entities.Guest, error)
}

type StorageCleaning interface {
	CreateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error)
	UpdateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error)
	GetCleaningByReservationID(ctx context.Context, reservationID int) (*entities.Cleaning, error)
}

func NewService(storage StorageReservation, storageGuest StorageGuest, storageCleaning StorageCleaning, svcReservationInfo ServiceReservationInfo) *Service {
	return &Service{
		storageReservation:     storage,
		storageGuest:           storageGuest,
		storageCleaning:        storageCleaning,
		serviceReservationInfo: svcReservationInfo,
	}
}

func (s *Service) UpdateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error) {
	r, err := s.storageReservation.GetReservationByID(ctx, booking.Reservation.Oid)
	if err != nil {
		zap.L().Error("UpdateBooking", zap.Error(err))
		return nil, err
	}
	if r == nil {
		zap.L().Error("UpdateBooking", zap.Error(erro.ErrEmptyResultFromReservation))
		return nil, erro.ErrEmptyResultFromReservation
	}

	g, err := s.storageGuest.ReadGuest(ctx, r.GuestID)
	if err != nil {
		zap.L().Error("UpdateBooking", zap.Error(err))
		return nil, err
	}
	if g == nil {
		zap.L().Error("UpdateBooking", zap.Error(erro.ErrReservationHasGuestUUIDbutGuestNotFound))
		return nil, erro.ErrReservationHasGuestUUIDbutGuestNotFound
	}

	updateGuest := new(entities.Guest)
	if booking.Guest.Phone != g.Phone {
		updateGuest, err = s.CreateGuest(ctx, booking.Guest)
		if err != nil {
			zap.L().Error("UpdateBooking", zap.Error(err))
			return nil, err
		}
	} else {
		booking.Guest.GuestID = g.GuestID
		updateGuest, err = s.storageGuest.UpdateGuest(ctx, booking.Guest)
		if err != nil {
			zap.L().Error("UpdateBooking", zap.Error(err))
			return nil, err
		}
	}

	booking.Reservation.GuestID = updateGuest.GuestID

	booking.Reservation = applyDefaultTimes(booking.Reservation)
	booking.Reservation = prepareDaysAndPriceForNight(booking.Reservation)

	updateReservation, err := s.storageReservation.UpdateReservation(ctx, booking.Reservation)
	if err != nil {
		zap.L().Error("UpdateBooking", zap.Error(err))
		return nil, err
	}

	//update reservation_info
	booking.ReservationInfo.PaymentOnCheckin = updateReservation.Price + updateReservation.CleaningPrice - booking.ReservationInfo.Prepayment
	booking.ReservationInfo.ReservationID = updateReservation.Oid
	booking.ReservationInfo.ActualCheckIn = updateReservation.CheckIn
	booking.ReservationInfo.ActualCheckOut = updateReservation.CheckOut

	updateReservationInfo, err := s.serviceReservationInfo.UpdateReservationInfoByReservationID(ctx, booking.ReservationInfo)
	if err != nil {
		zap.L().Error("UpdateBooking: failed to update reservation_info", zap.Error(err))
		return nil, err
	}

	// Синхронизируем связанную запись уборки: обновляем время (check_out), комнату и цену уборки.
	// Поля, которые менеджер заполняет вручную (agent_name, laundry_price, paid и др.), не трогаем.
	if err = s.syncCleaningAfterReservationUpdate(ctx, updateReservation); err != nil {
		return nil, err
	}

	b := entities.Booking{
		Guest:           *updateGuest,
		Reservation:     *updateReservation,
		ReservationInfo: *updateReservationInfo,
	}

	return &b, nil
}

func (s *Service) DeleteReservation(ctx context.Context, id int) (*entities.Reservation, error) {
	reserv, err := s.storageReservation.GetReservationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reserv != nil {
		_, err := s.storageReservation.Delete(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	zap.L().Info("DeleteReservation", zap.Any("reserv", reserv))
	return reserv, nil
}

func (s *Service) CreateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error) {
	guest, err := s.CreateGuest(ctx, booking.Guest)
	if err != nil {
		return nil, err
	}

	booking.Reservation.GuestID = guest.GuestID

	reservation, err := s.CreateReservation(ctx, booking.Reservation)
	if err != nil {
		return nil, err
	}

	// Создаём reservation_info с данными, переданными с фронта
	booking.ReservationInfo.ReservationID = reservation.Oid
	booking.ReservationInfo.PaymentOnCheckin = booking.Reservation.Price + booking.Reservation.CleaningPrice - booking.ReservationInfo.Prepayment
	// Фактическое время заезда/выезда по умолчанию совпадает с договорным
	booking.ReservationInfo.ActualCheckIn = reservation.CheckIn
	booking.ReservationInfo.ActualCheckOut = reservation.CheckOut

	info, err := s.serviceReservationInfo.CreateReservationInfo(ctx, booking.ReservationInfo)
	if err != nil {
		zap.L().Error("CreateBooking: failed to create reservation_info", zap.Error(err))
		return nil, err
	}

	b := entities.Booking{Guest: *guest, Reservation: *reservation, ReservationInfo: *info}
	return &b, nil
}

func (s *Service) CreateReservation(ctx context.Context, reservation entities.Reservation) (*entities.Reservation, error) {
	if reservation.GuestID == uuid.Nil {
		return nil, errors.New("uuid is nil")
	}

	reservation = applyDefaultTimes(reservation)
	if reservation.Days == 0 {
		res := prepareDaysAndPriceForNight(reservation)
		reservation = res
	}
	r, err := s.storageReservation.Create(ctx, reservation)
	if err != nil {
		zap.L().Debug("CreateReservation", zap.Error(err), zap.Any("booking", reservation))
		return nil, err
	}
	if r == nil {
		return nil, erro.ErrEmptyResultFromDB
	}

	// Автоматически создаём запись уборки: cleaning_time = время выезда гостя
	reservationID := r.Oid
	cleaning := entities.Cleaning{
		ReservationID: &reservationID,
		CleaningTime:  r.CheckOut,
		Room:          r.RoomNumber,
		CleaningPrice: r.CleaningPrice,
		LaundryPrice:  0,
		AgentName:     "Our Apartment",
		Description:   "",
		Paid:          true,
	}
	_, err = s.storageCleaning.CreateCleaning(ctx, cleaning)
	if err != nil {
		zap.L().Error("CreateReservation: failed to create cleaning record", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// applyDefaultTimes устанавливает время по умолчанию, если оно не указано (00:00:00):
// check_in → 14:00:00, check_out → 11:00:00
func applyDefaultTimes(reservation entities.Reservation) entities.Reservation {
	hIn, mIn, sIn := reservation.CheckIn.Clock()
	hOut, mOut, sOut := reservation.CheckOut.Clock()

	if hIn == 0 && mIn == 0 && sIn == 0 && hOut == 0 && mOut == 0 && sOut == 0 {
		reservation.CheckIn = time.Date(
			reservation.CheckIn.Year(),
			reservation.CheckIn.Month(),
			reservation.CheckIn.Day(),
			14, 0, 0, 0,
			reservation.CheckIn.Location(),
		)
		reservation.CheckOut = time.Date(
			reservation.CheckOut.Year(),
			reservation.CheckOut.Month(),
			reservation.CheckOut.Day(),
			11, 0, 0, 0,
			reservation.CheckOut.Location(),
		)
	}

	return reservation
}

// что бы не считать в коде стоимость ночи и количество дней нужно отдать это на вычеслении бд (раньше это делала бд в вычесляемых столбцах но при удаление контейнера почему всё пропало хотя и потключены volumes)
func prepareDaysAndPriceForNight(reservation entities.Reservation) entities.Reservation {
	reservation = countDays(reservation)

	return countPriceForOneNight(reservation)
}

func countDays(reservation entities.Reservation) entities.Reservation {
	out := reservation.CheckOut.UTC().Truncate(24 * time.Hour)
	in := reservation.CheckIn.UTC().Truncate(24 * time.Hour)

	reservation.Days = int(out.Sub(in).Hours() / 24)

	return reservation
}

func countPriceForOneNight(reservation entities.Reservation) entities.Reservation {
	if reservation.Price == 0 {
		reservation.PriceForOneNight = 0
		return reservation
	}
	reservation.PriceForOneNight = int(reservation.Price / reservation.Days)
	return reservation
}

func (s *Service) CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	guest, err := s.storageGuest.FindGuestByPhoneNumber(ctx, g.Phone)
	if err != nil {
		return nil, err
	}
	if guest == nil {
		if g.GuestID == uuid.Nil {
			g.GuestID = uuid.New()
		}
		guest, err = s.storageGuest.CreateGuest(ctx, g)
		if err != nil {
			return nil, err
		}

		if guest == nil {
			return nil, erro.ErrEmptyResultFromDB
		}
	}

	return guest, nil
}

func (s *Service) CreateReservation1(ctx context.Context, booking entities.Booking) error {
	guest, err := s.storageGuest.FindGuestByPhoneNumber(ctx, booking.Phone)
	if err != nil {
		return err
	}
	if guest == nil {
		if booking.Guest.GuestID == uuid.Nil {
			booking.Guest.GuestID = uuid.New()
		}
		guest, err = s.storageGuest.CreateGuest(ctx, booking.Guest)
		if err != nil {
			return err
		}
		booking.Reservation.GuestID = guest.GuestID
	} else {
		zap.L().Debug("CreateReservation", zap.Error(erro.ErrGuestAlreadyExist), zap.Any("Guest", booking.Guest))
		booking.Guest = *guest
		booking.Reservation.GuestID = guest.GuestID
	}

	//чекнем пересечение по датам на эту комнату
	values, err := s.storageReservation.ReadWithRoomNumber(ctx, booking.RoomNumber, booking.CheckIn, booking.CheckOut)
	if err != nil {
		return err
	}
	if len(values) > 0 {
		if len(values) == 1 {
			if values[0].CheckIn.UTC() == booking.CheckIn &&
				values[0].CheckOut.UTC() == booking.CheckOut &&
				values[0].Price == booking.Price &&
				values[0].RoomNumber == booking.RoomNumber {
				return fmt.Errorf("%w\nзаписываемое значение: %v;\nзначение из БД: %v;", erro.ErrFullyMatchOtherBooking, booking, values[0])
			}
		}
		return errors.New(fmt.Sprintf("Букинг: %v ; пересекается со следуюшим бронированиями: %v ;", booking, values))
	}

	//запишем новое бронирование в бд
	_, err = s.storageReservation.Create(ctx, booking.Reservation)
	if err != nil {
		zap.L().Debug("CreateReservation", zap.Error(err), zap.Any("booking", booking))
		return err
	}
	return nil
}

// CreateReport Создание репорта для собственника
func (s *Service) CreateReport(ctx context.Context, roomNumber string, startPeriod string, endPeriod string) ([]byte, error) {
	bookings, err := s.GetBooking(ctx, roomNumber, startPeriod, endPeriod)
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, erro.ErrSliceOfBookingIsEmpty
	}

	report, err := report.ReportForOwner(bookings)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *Service) GetBookingALLForApartmentALL(ctx context.Context, roomNumbers []string) ([]entities.Booking, error) {
	allBokings := make([]entities.Booking, 0, 500)
	for _, room := range roomNumbers {
		bookings, err := s.GetBookingALLForApartment(ctx, room)
		if err != nil {
			return nil, err
		}
		allBokings = append(allBokings, bookings...)
	}
	return allBokings, nil
}

func (s *Service) GetBookingALLForApartment(ctx context.Context, roomNumber string) ([]entities.Booking, error) {
	bookings := make([]entities.Booking, 0, 50)
	reservation, err := s.GetReservationALLForApartment(ctx, roomNumber)
	if err != nil {
		return nil, err
	}
	for _, r := range reservation {
		guest, err := s.storageGuest.ReadGuest(ctx, r.GuestID)
		if err != nil {
			return nil, err
		}
		if guest == nil {
			zap.L().Error("GetBookingALLForApartment", zap.Error(erro.ErrReservationHasGuestUUIDbutGuestNotFound))
			return nil, erro.ErrReservationHasGuestUUIDbutGuestNotFound
		}

		rInfo, err := s.serviceReservationInfo.GetReservationInfoByReservationID(ctx, r.Oid)
		if err != nil {
			return nil, err
		}
		if rInfo == nil {
			zap.L().Debug("GetBookingALLForApartment: GetReservationInfoByReservationID", zap.Error(erro.ErrEmptyResultFromDB), zap.Int("reservationID", r.Oid))
			rInfo = new(entities.ReservationInfo)
		}

		booking := entities.Booking{
			Guest:           *guest,
			Reservation:     r,
			ReservationInfo: *rInfo,
		}

		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *Service) GetReservationALLForApartment(ctx context.Context, roomNumber string) ([]entities.Reservation, error) {
	reservations, err := s.storageReservation.ReadALLByRoomNumber(ctx, roomNumber)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (s *Service) GetBooking(ctx context.Context, roomNumber string, start string, end string) ([]entities.Booking, error) {
	bookings := make([]entities.Booking, 0, 5)
	reservation, err := s.GetReservationForPeriodByApartment(ctx, roomNumber, start, end)
	if err != nil {
		return nil, err
	}
	for _, r := range reservation {
		guest, err := s.storageGuest.ReadGuest(ctx, r.GuestID)
		if err != nil {
			return nil, err
		}
		if guest == nil {
			zap.L().Error("GetBooking", zap.Error(erro.ErrReservationHasGuestUUIDbutGuestNotFound))
			return nil, erro.ErrReservationHasGuestUUIDbutGuestNotFound
		}
		booking := entities.Booking{
			Guest:       *guest,
			Reservation: r,
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *Service) GetReservationForPeriodByApartment(ctx context.Context, roomNumber string, start string, end string) ([]entities.Reservation, error) {
	checkin, err := models.TimeConvert(start)
	if err != nil {
		return nil, err
	}

	checkout, err := models.TimeConvert(end)
	if err != nil {
		return nil, err
	}

	reservations, err := s.storageReservation.ReadWithRoomNumber(ctx, roomNumber, checkin, checkout)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (s *Service) GetReservationByPhoneNumber(ctx context.Context, phone string) ([]entities.Reservation, error) {
	guest, err := s.storageGuest.FindGuestByPhoneNumber(ctx, phone)
	if err != nil {
		return nil, err
	}
	if guest == nil {
		return nil, nil
	}
	bookings, err := s.storageReservation.FindBookingByGuestUUID(ctx, guest.GuestID)
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, nil
	}
	return bookings, nil
}

func (s *Service) GetBookingByCheckIn(ctx context.Context, date time.Time) ([]entities.Booking, error) {
	infos, err := s.serviceReservationInfo.GetReservationInfosByActualCheckIn(ctx, date)
	if err != nil {
		return nil, err
	}
	return s.buildBookingsFromReservationInfos(ctx, infos)
}

func (s *Service) GetBookingByCheckOut(ctx context.Context, date time.Time) ([]entities.Booking, error) {
	infos, err := s.serviceReservationInfo.GetReservationInfosByActualCheckOut(ctx, date)
	if err != nil {
		return nil, err
	}
	return s.buildBookingsFromReservationInfos(ctx, infos)
}

// buildBookingsFromReservationInfos собирает []Booking по списку reservation_info:
// для каждого достаёт резервацию и гостя.
func (s *Service) buildBookingsFromReservationInfos(ctx context.Context, infos []entities.ReservationInfo) ([]entities.Booking, error) {
	bookings := make([]entities.Booking, 0, len(infos))
	for _, ri := range infos {
		r, err := s.storageReservation.GetReservationByID(ctx, ri.ReservationID)
		if err != nil {
			return nil, err
		}
		if r == nil {
			zap.L().Error("buildBookingsFromReservationInfos: reservation not found", zap.Int("reservationID", ri.ReservationID))
			return nil, erro.ErrEmptyResultFromReservation
		}

		guest, err := s.storageGuest.ReadGuest(ctx, r.GuestID)
		if err != nil {
			return nil, err
		}
		if guest == nil {
			zap.L().Error("buildBookingsFromReservationInfos", zap.Error(erro.ErrReservationHasGuestUUIDbutGuestNotFound))
			return nil, erro.ErrReservationHasGuestUUIDbutGuestNotFound
		}

		bookings = append(bookings, entities.Booking{Guest: *guest, Reservation: *r, ReservationInfo: ri})
	}
	return bookings, nil
}

// Оповещение собственика о предстоящем бронирование
func FutureBooking() {

}

// свободные квартиры на эти даты
func FreeApartmentForDates() {

}

func (s *Service) FindTotalPriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, startPeriod, endPeriod string) (map[string]int, error) {
	m := make(map[string]int)

	for _, apartment := range apartments {
		price, _, err := s.FindTotalPriceForPeriod(ctx, apartment.RoomNumber, startPeriod, endPeriod)
		if err != nil {
			return nil, err
		}
		m[apartment.RoomNumber] = price
	}
	return m, nil
}

func (s *Service) FindMiddlePriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, startPeriod, endPeriod string) (map[string]int, error) {
	m := make(map[string]int)

	for _, apartment := range apartments {
		price, err := s.FindMiddlePriceForPeriod(ctx, apartment.RoomNumber, startPeriod, endPeriod)
		if err != nil {
			return nil, err
		}
		m[apartment.RoomNumber] = price
	}
	return m, nil
}

// цены расчитаные по переуду, к примеру низкий сезон переходящий в высокий
func (s *Service) FindMiddlePriceForPeriod(ctx context.Context, roomNumber string, startPeriod, endPeriod string) (int, error) {
	totalSum, totalDays, err := s.FindTotalPriceForPeriod(ctx, roomNumber, startPeriod, endPeriod)
	if err != nil {
		return 0, err
	}
	if totalSum == 0 || totalDays == 0 {
		return 0, nil
	}
	return totalSum / totalDays, nil
}

func (s *Service) FindTotalPriceForPeriod(ctx context.Context, roomNumber, startPeriod, endPeriod string) (int, int, error) {
	start, err := models.TimeConvert(startPeriod)
	if err != nil {
		return 0, 0, err
	}
	end, err := models.TimeConvert(endPeriod)
	if err != nil {
		return 0, 0, err
	}

	if !start.Before(end) {
		return 0, 0, erro.ErrStartDateIsNotBeforeEndDate
	}
	bookings, err := s.storageReservation.ReadWithRoomNumber(ctx, roomNumber, start, end)
	if err != nil {
		return 0, 0, err
	}
	if len(bookings) == 0 {
		return 0, 0, nil
	}

	var totalDays int
	var totalSum int

	for _, booking := range bookings {
		if booking.Price != 0 {
			switch {
			// кейс когда чекин раньше start а чекаут раньше end
			case (start.After(booking.CheckIn) || start.Equal(booking.CheckIn)) && (end.After(booking.CheckOut) || end.Equal(booking.CheckOut)):
				days := helperForMiddlePrice(booking.CheckOut, start)
				totalSum += days * booking.PriceForOneNight
				totalDays += days
				zap.L().Debug("кейс когда чекин раньше start а чекаут раньше end", zap.Int("days", days), zap.Int("booking.PriceForOneNight", booking.PriceForOneNight))
			// кейс когда чекаут позже end и чекин позже start
			case (start.Before(booking.CheckIn) || start.Equal(booking.CheckIn)) && (end.Before(booking.CheckOut) || end.Equal(booking.CheckOut)):
				days := helperForMiddlePrice(end, booking.CheckIn)
				totalSum += days * booking.PriceForOneNight
				totalDays += days
				zap.L().Debug("кейс когда чекаут позже end и чекин позже start", zap.Int("days", days), zap.Int("booking.PriceForOneNight", booking.PriceForOneNight))
			// кейс когда чекин и чекаут внутри периуда start и end
			case (start.Before(booking.CheckIn) || start.Equal(booking.CheckIn)) && (end.After(booking.CheckOut) || end.Equal(booking.CheckOut)):
				days := helperForMiddlePrice(booking.CheckOut, booking.CheckIn)
				totalSum += days * booking.PriceForOneNight
				totalDays += days
				zap.L().Debug("кейс когда чекин и чекаут внутри периуда start и end", zap.Int("days", days), zap.Int("booking.PriceForOneNight", booking.PriceForOneNight))
			//кейс когда чекин и чекаут за периудом start и end
			case (start.After(booking.CheckIn) || start.Equal(booking.CheckIn)) && (end.Before(booking.CheckOut) || end.Equal(booking.CheckOut)):
				days := helperForMiddlePrice(end, start)
				totalSum += days * booking.PriceForOneNight
				totalDays += days
				zap.L().Debug("кейс когда чекин и чекаут за периудом start и end", zap.Int("days", days), zap.Int("booking.PriceForOneNight", booking.PriceForOneNight))
			}
		}

	}

	return totalSum, totalDays, nil
}

// helperForMiddlePrice result of t-u in days
func helperForMiddlePrice(t, u time.Time) int {
	return int(t.Sub(u).Hours() / 24)
}

// syncCleaningAfterReservationUpdate обновляет связанную запись уборки после изменения резервации.
// Если записи уборки нет — ничего не делаем.
func (s *Service) syncCleaningAfterReservationUpdate(ctx context.Context, r *entities.Reservation) error {
	cleaning, err := s.storageCleaning.GetCleaningByReservationID(ctx, r.Oid)
	if err != nil {
		zap.L().Error("syncCleaningAfterReservationUpdate: GetCleaningByReservationID", zap.Error(err))
		return err
	}
	if cleaning == nil {
		return nil
	}

	cleaning.CleaningTime = r.CheckOut
	cleaning.Room = r.RoomNumber
	cleaning.CleaningPrice = r.CleaningPrice

	_, err = s.storageCleaning.UpdateCleaning(ctx, *cleaning)
	if err != nil {
		zap.L().Error("syncCleaningAfterReservationUpdate: UpdateCleaning", zap.Error(err))
		return err
	}
	return nil
}
