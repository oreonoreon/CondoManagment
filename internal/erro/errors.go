package erro

import "fmt"

type ErrCode int32

type Err struct {
	Code ErrCode `json:"code"`
	Desc string  `json:"message"`
}

func (e *Err) Error() string {
	return fmt.Sprintf("{%d %s}", e.Code, e.Desc)
}

func (e *Err) ErrCode() int {
	return int(e.Code)
}

func (e *Err) Message() string {
	return e.Desc
}

var (
	ErrFullyMatchOtherBooking                        = &Err{1000, "Букинг полностью совпадает с:"}
	ErrMatchWithOtherBooking                         = &Err{1001, "Букинг пересекается с другими бронированиями"}
	ErrReservationHasGuestUUIDbutGuestNotFound       = &Err{1002, "Reservations have the guest uuid but in Guests table guest was not found"}
	ErrGuestAlreadyExist                             = &Err{1003, "Guest already exist"}
	ErrBookingExcelModelHaveNotName                  = &Err{1004, "While parsing booking excel model the Name of Guest was not found"}
	ErrBookingExcelModelHaveNotPhone                 = &Err{1005, "While parsing booking excel model the Phone was not found"}
	ErrBookingExcelModelDbCovertCheckInAfterCheckOut = &Err{1006, "The check_in > or = check_out was found"}

	ErrSliceOfBookingIsEmpty = &Err{1007, "There is no any booking for given period or apartment"}

	ErrNoFoundBookingsFromExcel = &Err{1008, "The bookings was not found in excel file"}

	ErrStartDateIsNotBeforeEndDate = &Err{1009, "The start date is not before end date"}

	ErrEmptyResultFromDB          = &Err{1010, "Result from DB returned empty"}
	ErrEmptyResultFromReservation = &Err{1011, "Reservation was not found by id"}

	ErrWrongCreds = &Err{1012, "Invalid credentials"}
)
