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
	ErrFullyMatchOtherBooking = &Err{1000, "Букинг полностью совпадает с:"}
	ErrMatchWithOtherBooking  = &Err{1001, "Букинг пересекается со следуюшим бронированиями:"}
)
