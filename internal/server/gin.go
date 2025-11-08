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

	allowOrigin, ok := os.LookupEnv("FRONT_URL")
	if !ok {
		allowOrigin = "http://localhost:5173" // явно указываем для разработки
	}

	// Определяем режим работы
	var isProduction bool
	ginMode, ok := os.LookupEnv("GIN_MODE")
	if !ok {
		isProduction = false
		zap.L().Info("GIN_MODE not set, default to debug")
	} else {
		if ginMode == "release" {
			isProduction = true
		} else {
			isProduction = false
		}
	}

	// CORS конфигурация
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.LoadHTMLGlob("html/*.html")

	db, err := sql.Open("postgres", repo.DataSourceName)
	if err != nil {
		panic(err)
	}

	store, err := postgres.NewStore(db, []byte("везде жопа лысого ПОПУГАЯ!"))
	if err != nil {
		panic(err)
	}

	// ИСПРАВЛЕНИЕ: правильные настройки cookies в зависимости от окружения
	cookieOptions := sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24, // 24 часа
		HttpOnly: true,
		Secure:   isProduction,         // true только в production с HTTPS
		SameSite: http.SameSiteLaxMode, // для разработки Lax, для production можно None
	}

	// Если production и используется cross-origin, нужен SameSite=None
	if isProduction {
		cookieOptions.SameSite = http.SameSiteNoneMode
	}

	store.Options(cookieOptions)
	router.Use(sessions.Sessions("sess", store))

	// Добавим middleware для логирования cookies (для отладки)
	//router.Use(func(c *gin.Context) {
	//	zap.L().Debug("Request cookies", zap.String("cookies", c.Request.Header.Get("Cookie")))
	//	c.Next()
	//	zap.L().Debug("Response Set-Cookie", zap.Any("set-cookie", c.Writer.Header()["Set-Cookie"]))
	//})

	router.POST("/login", h.LoginHandler)

	// Защищённые маршруты
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
		api.POST("/rall", h.AllBookingsPost)
		api.GET("/r", h.ApartmentsGet)
		api.PATCH("/updateBooking", h.UpdateBooking)
		api.POST("/createBooking", h.CreateBookingPost)
		api.DELETE("/deleteBooking/:id", h.DeleteBookingByID)
		api.POST("/BnB", h.ScrapBnBPost)
		api.POST("/BnB/locationName", h.ScrapBnBLocationNameUpdate)
		api.POST("/BnB/room", h.ScrapBnbRoomUnderstandableTypePatch)
	}

	router.Run(":8080")
}

type Handle struct {
	*services.Service
	*services.TransactionalService
	*services.ServiceSettings
	*services.ServiceExcel
	*services.ServiceApartment
	*services.ServiceBnB
	*services.ServiceUsers
}

func NewHandle(
	serviceReservation *services.Service,
	servicesInterface *services.TransactionalService,
	serviceSettings *services.ServiceSettings,
	serviceExcel *services.ServiceExcel,
	serviceApartment *services.ServiceApartment,
	servicesBnB *services.ServiceBnB,
	servicesUsers *services.ServiceUsers) Handle {
	return Handle{
		serviceReservation,
		servicesInterface,
		serviceSettings,
		serviceExcel,
		serviceApartment,
		servicesBnB,
		servicesUsers,
	}
}

func (h *Handle) UpdateBooking(c *gin.Context) {
	request := new(entities.Booking)

	err := c.BindJSON(request)
	if err != nil {
		zap.L().Error("UpdateBooking", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	booking, err := h.TransactionalService.UpdateBooking(c.Request.Context(), *request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, booking)
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

	booking, err := h.TransactionalService.CreateBooking(c.Request.Context(), *request)
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

type AllBookingsGetRequest struct {
	RoomNumbers []string `json:"room_numbers"`
}

func (h *Handle) AllBookingsPost(c *gin.Context) {
	//request := make([]string, 0, 10)
	request := new(AllBookingsGetRequest)
	err := c.BindJSON(request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	bookings, err := h.Service.GetBookingALLForApartmentALL(c.Request.Context(), request.RoomNumbers)
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
