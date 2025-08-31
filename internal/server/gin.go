package server

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/services"
	"database/sql"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Gin(h Handle) {
	router := gin.Default()

	allowOrigin, ok := os.LookupEnv("FRONT_URL") //todo вынести все env в подобающее место
	if !ok {
		allowOrigin = "*"
	}

	//---------------------------
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowOrigin}, // или "*" для всех
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//---------------------------

	router.LoadHTMLGlob("html/*.html") // шаблоны

	db, err := sql.Open("postgres", repo.DataSourceName) //todo не хорошо что тут используем пакет repo
	if err != nil {
		panic(err)
	}

	store, err := postgres.NewStore(db, []byte("везде жопа лысого ПОПУГАЯ!"))
	if err != nil {
		panic(err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24, // день
		HttpOnly: true,
		Secure:   true, // в проде true при https
	})
	router.Use(sessions.Sessions("sess", store))

	router.POST("/login", h.LoginHandler)

	// Защищённые маршруты: сначала сессия, потом authorizer
	api := router.Group("/calendar")
	api.Use(SessionAuthMiddleware())
	{
		api.POST("/createUser", h.CreateUser)

		api.POST("/logout", h.LogoutHandler)

		api.GET("/sync", h.SynchroniseBookings)
		api.POST("/sync", h.SynchroniseBookings)

		api.GET("/ExcelRes", h.ExcelBookings)
		api.POST("/ExcelRes", h.ExcelBookings)

		api.GET("/middleprice", h.MiddlePriceForPeriod)
		api.POST("/middleprice", h.MiddlePriceForPeriod)

		api.GET("/middlepriceReport", h.MiddlePriceForPeriodReport)
		api.POST("/middlepriceReport", h.MiddlePriceForPeriodReport)

		api.GET("/totalpriceReport", h.TotalPriceForPeriodReport)
		api.POST("/totalpriceReport", h.TotalPriceForPeriodReport)

		api.GET("/report", h.Report)
		api.POST("/report", h.Report)

		api.POST("/r", h.BookingsPost)

		api.GET("/r", h.ApartmentsGet)

		api.PATCH("/updateBooking/:id", h.UpdateBooking)

		api.POST("/createBooking", h.CreateBookingPost)

		api.DELETE("/deleteBooking/:id", h.DeleteBookingByID)

		api.POST("/BnB", h.ScrapBnBPost)

		api.POST("/BnB/locationName", h.ScrapBnBLocationNameUpdate)

		api.POST("/BnB/room", h.ScrapBnbRoomUnderstandableTypePatch)
	}
	//router.POST("/BnB", h.ScrapBnBPost)
	//
	//router.POST("/BnB/locationName", h.ScrapBnBLocationNameUpdate)
	//
	//router.POST("/BnB/room", h.ScrapBnbRoomUnderstandableTypePatch)

	router.Run(":8080") // gin сам управляет тайм-аутами, но можно кастомизировать
}

type Handle struct {
	*services.Service
	*services.ServiceSettings
	*services.ServiceExcel
	*services.ServiceApartment
	*services.ServiceBnB
	*services.ServiceUsers
}

func NewHandle(
	serviceReservation *services.Service,
	serviceSettings *services.ServiceSettings,
	serviceExcel *services.ServiceExcel,
	serviceApartment *services.ServiceApartment,
	servicesBnB *services.ServiceBnB,
	servicesUsers *services.ServiceUsers) Handle {
	return Handle{
		serviceReservation,
		serviceSettings,
		serviceExcel,
		serviceApartment,
		servicesBnB,
		servicesUsers,
	}
}

func (h *Handle) UpdateBooking(c *gin.Context) {
	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	zap.L().Error("UpdateBooking", zap.Error(err))
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}
	//
	//request := new(entities.Booking)
	//
	//err := c.BindJSON(request)
	//if err != nil {
	//	zap.L().Error("CreateBookingPost", zap.Error(err))
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}
	//if request == nil {
	//	erro := errors.New("request contain error")
	//	zap.L().Error("CreateBookingPost", zap.Error(erro))
	//	c.String(http.StatusBadRequest, erro.Error())
	//}

}

func (h *Handle) DeleteBookingByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.L().Error("DeleteBookingByID", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	reservation, err := h.Service.DeleteReservation(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, reservation)
}

func (h *Handle) CreateBookingPost(c *gin.Context) {
	request := new(entities.Booking)

	err := c.BindJSON(request)
	if err != nil {
		zap.L().Error("CreateBookingPost", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if request == nil {
		erro := errors.New("request contain error")
		zap.L().Error("CreateBookingPost", zap.Error(erro))
		c.String(http.StatusBadRequest, erro.Error())
	}

	booking, err := h.Service.CreateBooking(c.Request.Context(), *request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, booking)
}

// ApartmentsGet request to get names of all appartment
func (h *Handle) ApartmentsGet(c *gin.Context) {
	apartments, err := h.ServiceApartment.GetAllApartment(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"apartments": apartments,
	})
}

type BookingsGetRequest struct {
	RoomNumber string `json:"room_number"`
}

func (h *Handle) BookingsPost(c *gin.Context) {
	request := new(BookingsGetRequest)
	err := c.BindJSON(request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	bookings, err := h.Service.GetBookingALLForApartment(c.Request.Context(), request.RoomNumber)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"bookings": bookings,
	})
}

func (h *Handle) MiddlePriceForPeriod(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "middleprice.html", nil)
		return
	}

	roomNumber := c.PostForm("room_number")
	start := c.PostForm("start")
	end := c.PostForm("end")

	price, err := h.Service.FindMiddlePriceForPeriod(c.Request.Context(), roomNumber, start, end)
	if err != nil {
		zap.L().Error("FindMiddlePriceForPeriod", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, strconv.Itoa(price))
}

func (h *Handle) SynchroniseBookings(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "sync.html", nil)
		return
	}

	sheetName := c.PostForm("sheet_name")
	roomNumber := c.PostForm("room_number")
	start := c.PostForm("start")
	end := c.PostForm("end")

	bookings, bookingsFromDB, err := h.ServiceExcel.Sync(c.Request.Context(), sheetName, roomNumber, start, end)
	if err != nil {
		zap.L().Error("Sync", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "sync.html", gin.H{
		"Bookings":       bookings,
		"BookingsFromDB": bookingsFromDB,
	})
}

func (h *Handle) ExcelBookings(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "ExcelRes.html", nil)
		return
	}

	sheetName := c.PostForm("sheet_name")
	roomNumber := c.PostForm("room_number")
	start := c.PostForm("start")
	end := c.PostForm("end")

	if roomNumber == "All" || roomNumber == "all" {
		bookings, err := h.ServiceExcel.GetAllBookingInfoForPeriod(c.Request.Context(), sheetName, start, end)
		if err != nil {
			zap.L().Error("ExcelBookings", zap.Error(err))
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.HTML(http.StatusOK, "ExcelRes.html", gin.H{
			"Bookings": bookings,
		})
	} else {
		bookings, err := h.ServiceExcel.GetBookingInfoForPeriod(c.Request.Context(), sheetName, roomNumber, start, end)
		if err != nil {
			zap.L().Error("ExcelBookings", zap.Error(err))
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.HTML(http.StatusOK, "ExcelRes.html", gin.H{
			"Bookings": bookings,
		})
	}

}
