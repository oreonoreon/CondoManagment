package server

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/entities"
	"awesomeProject/internal/services"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SpecialStore struct {
	*pgstore.PGStore
}

// sameSiteString переводит http.SameSite в читаемую строку для логов.
func sameSiteString(s http.SameSite) string {
	switch s {
	case http.SameSiteNoneMode:
		return "None"
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteStrictMode:
		return "Strict"
	default:
		return "Default"
	}
}

func (s SpecialStore) Options(options sessions.Options) {
	s.PGStore.Options = options.ToGorillaOptions()
}

func Gin(h Handle) {
	router := gin.Default()

	// CORS конфигурация
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{h.cfg.FrontURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.LoadHTMLGlob("html/*.html")

	//------------------ Всё это неплохо бы вынести --------------
	db, err := sql.Open("postgres", h.cfg.DatabaseDSN) //todo не хорошо что здесь идёт повторное открытие соединения с бд
	if err != nil {
		panic(err)
	}

	notSpesialStore, err := pgstore.NewPGStoreFromPool(db, []byte(h.cfg.SessionSecret))
	if err != nil {
		panic(err)
	}

	defer notSpesialStore.StopCleanup(notSpesialStore.Cleanup(time.Hour * 24))

	store := SpecialStore{notSpesialStore}
	//-----------------------------------------------------------------------------

	// Настройки cookies. Frontend теперь всегда обращается к бэкенду через
	// собственный reverse-proxy (тот же origin), поэтому SameSite=Lax
	// достаточно и в dev, и в production — это same-site (и чаще даже
	// same-origin) запрос, а Lax не блокируется ITP/Safari в отличие от None.
	// Secure остаётся завязан на production, т.к. там всегда HTTPS.
	cookieOptions := sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24, // 24 часа
		HttpOnly: true,
		Secure:   h.cfg.IsProduction, // true только в production с HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	store.Options(cookieOptions)
	router.Use(sessions.Sessions("sess", store))

	// Явно логируем реально применённые настройки, влияющие на куки/CORS,
	// чтобы на проде (Railway) сразу было видно, что подхватилось из ENV.
	zap.L().Info("Session/CORS runtime settings",
		zap.Bool("IsProduction", h.cfg.IsProduction),
		zap.String("FrontURL", h.cfg.FrontURL),
		zap.Bool("CookieSecure", cookieOptions.Secure),
		zap.String("CookieSameSite", sameSiteString(cookieOptions.SameSite)),
	)

	router.POST("/login", h.LoginHandler)

	// Защищённые маршруты
	api := router.Group("/calendar")
	// Права проверяются по role_permissions (см. migrations/014_roles_and_permissions.sql).
	api.Use(h.SessionAuthMiddleware())
	{
		read := RequirePermission(entities.PermCalendarRead)
		write := RequirePermission(entities.PermCalendarWrite)
		reports := RequirePermission(entities.PermReportsRead)
		tools := RequirePermission(entities.PermToolsUse)
		today := RequirePermission(entities.PermTodayRead)

		api.GET("/me", h.Me)
		api.POST("/logout", h.LogoutHandler)
		api.POST("/createUser", RequirePermission(entities.PermUsersManage), h.CreateUser)

		api.GET("/sync", tools, h.SynchroniseBookings)
		api.POST("/sync", tools, h.SynchroniseBookings)
		api.GET("/ExcelRes", tools, h.ExcelBookings)
		api.POST("/ExcelRes", tools, h.ExcelBookings)
		api.POST("/BnB", tools, h.ScrapBnBPost)
		api.POST("/BnB/locationName", tools, h.ScrapBnBLocationNameUpdate)
		api.POST("/BnB/room", tools, h.ScrapBnbRoomUnderstandableTypePatch)
		//api.GET("/middleprice", h.MiddlePriceForPeriod)
		//api.POST("/middleprice", h.MiddlePriceForPeriod)
		//api.GET("/middlepriceReport", h.MiddlePriceForPeriodReport)
		//api.POST("/middlepriceReport", h.MiddlePriceForPeriodReport)

		api.POST("/totalpriceReport", reports, h.TotalPriceForPeriodReport)
		api.POST("/totalpriceReportXlsx", reports, h.TotalPriceForPeriodReportXlsx)
		api.POST("/report", reports, h.Report)

		//depricated
		//api.POST("/r", h.BookingsPost)
		//api.POST("/rall", h.AllBookingsPost)

		api.POST("/rall", read, h.AllBookingsPostNEW)
		api.GET("/r", read, h.ApartmentsGet)
		api.PATCH("/updateBooking", write, h.UpdateBooking)
		api.POST("/createBooking", write, h.CreateBookingPost)
		api.DELETE("/deleteBooking/:id", write, h.DeleteBookingByID)

		// Bookings by date (вкладка Today)
		api.GET("/bookings/check-in/:date", today, h.GetBookingsByCheckIn)
		api.GET("/bookings/check-out/:date", today, h.GetBookingsByCheckOut)

		// Cleaning CRUD — без ограничения по апартаментам
		cleaning := api.Group("/cleaning", RequirePermission(entities.PermCleaningManage))
		cleaning.GET("", h.GetAllCleaning)
		cleaning.GET("/date/:date", h.GetCleaningByDate)
		cleaning.GET("/:id", h.GetCleaningByID)
		cleaning.POST("", h.CreateCleaning)
		cleaning.PATCH("/:id", h.UpdateCleaning)
		cleaning.DELETE("/:id", h.DeleteCleaning)

		// ReservationInfo
		api.PATCH("/reservation-info/:id", write, h.UpdateReservationInfo)

		// Statuses
		api.GET("/status-types", read, h.GetStatusTypes)
		api.GET("/reservations/:id/statuses", read, h.GetReservationStatuses)
		api.POST("/reservations/:id/statuses/:statusTypeId", write, h.ToggleReservationStatus)

		// Админка
		admin := api.Group("/admin")
		usersManage := RequirePermission(entities.PermUsersManage)
		apartmentsManage := RequirePermission(entities.PermApartmentsManage)
		rolesManage := RequirePermission(entities.PermRolesManage)

		// список пользователей нужен и для выдачи доступа к апартаментам
		admin.GET("/users", RequireAnyPermission(entities.PermUsersManage, entities.PermApartmentsManage), h.AdminListUsers)
		admin.POST("/users", usersManage, h.AdminCreateUser)
		admin.PATCH("/users/:id", usersManage, h.AdminUpdateUser)
		admin.DELETE("/users/:id", usersManage, h.AdminDeleteUser)

		admin.GET("/apartments", apartmentsManage, h.AdminListApartments)
		admin.POST("/apartments", apartmentsManage, h.AdminCreateApartment)
		admin.PATCH("/apartments/:id", apartmentsManage, h.AdminUpdateApartment)
		admin.PUT("/apartments/:id/access", apartmentsManage, h.AdminSetApartmentAccess)

		// список ролей нужен и для выбора роли пользователю, и для выдачи доступа к апартаментам
		admin.GET("/roles", RequireAnyPermission(entities.PermUsersManage, entities.PermApartmentsManage, entities.PermRolesManage), h.AdminListRoles)
		admin.GET("/permissions", rolesManage, h.AdminPermissionCatalog)
		admin.POST("/roles", rolesManage, h.AdminCreateRole)
		admin.PATCH("/roles/:name", rolesManage, h.AdminUpdateRole)
		admin.DELETE("/roles/:name", rolesManage, h.AdminDeleteRole)
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
	access *services.ServiceAccess // не встраиваем, чтобы не конфликтовали имена методов
	admin  *services.ServiceAdmin
	cfg    *config.ConfigEnv // добавляем поле конфига
}

func NewHandle(
	serviceReservation *services.Service,
	servicesInterface *services.TransactionalService,
	serviceSettings *services.ServiceSettings,
	serviceExcel *services.ServiceExcel,
	serviceApartment *services.ServiceApartment,
	servicesBnB *services.ServiceBnB,
	servicesUsers *services.ServiceUsers,
	serviceAccess *services.ServiceAccess,
	serviceAdmin *services.ServiceAdmin,
	cfg *config.ConfigEnv) Handle {
	return Handle{
		serviceReservation,
		servicesInterface,
		serviceSettings,
		serviceExcel,
		serviceApartment,
		servicesBnB,
		servicesUsers,
		serviceAccess,
		serviceAdmin,
		cfg,
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
	// проверяем и текущий апартамент брони, и тот, на который её переносят
	if !h.requireReservation(c, request.Reservation.Oid) || !h.requireRooms(c, request.Reservation.RoomNumber) {
		return
	}
	booking, err := h.TransactionalService.UpdateBooking(c.Request.Context(), *request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	view, ok := bookingViewFor(c, *booking)
	if !ok {
		return
	}
	c.IndentedJSON(http.StatusOK, view)
}

func (h *Handle) DeleteBookingByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.L().Error("DeleteBookingByID", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if !h.requireReservation(c, id) {
		return
	}

	reservation, err := h.TransactionalService.DeleteReservation(c.Request.Context(), id)
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
	if !h.requireRooms(c, request.Reservation.RoomNumber) {
		return
	}

	booking, err := h.TransactionalService.CreateBooking(c.Request.Context(), *request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	view, ok := bookingViewFor(c, *booking)
	if !ok {
		return
	}
	c.IndentedJSON(http.StatusCreated, view)
}

// ApartmentsGet request to get names of all appartment
func (h *Handle) ApartmentsGet(c *gin.Context) {
	user, ok := mustCurrentUser(c)
	if !ok {
		return
	}

	apartments, err := h.ServiceApartment.GetApartmentsFor(c.Request.Context(), user)
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

//depricated
//func (h *Handle) BookingsPost(c *gin.Context) {
//	request := new(BookingsGetRequest)
//	err := c.BindJSON(request)
//	if err != nil {
//		c.String(http.StatusBadRequest, err.Error())
//		return
//	}
//
//	bookings, err := h.TransactionalService.GetBookingALLForApartment(c.Request.Context(), request.RoomNumber)
//	if err != nil {
//		c.String(http.StatusInternalServerError, err.Error())
//		return
//	}
//	c.IndentedJSON(http.StatusOK, gin.H{
//		"bookings": bookings,
//	})
//}

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

	bookings, err := h.TransactionalService.GetBookingALLForApartmentALL(c.Request.Context(), request.RoomNumbers)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"bookings": bookings,
	})
}

func (h *Handle) AllBookingsPostNEW(c *gin.Context) {
	request := new(AllBookingsGetRequest)
	if err := c.BindJSON(request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	access, ok := h.roomAccess(c)
	if !ok {
		return
	}
	roomNumbers := filterByRoom(access, request.RoomNumbers, func(room string) string { return room })

	bookings, err := h.TransactionalService.GetBookingsForRooms(c.Request.Context(), roomNumbers)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	view, ok := bookingsViewFor(c, bookings)
	if !ok {
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"bookings": view,
	})
}

//func (h *Handle) MiddlePriceForPeriod(c *gin.Context) {
//	if c.Request.Method == http.MethodGet {
//		c.HTML(http.StatusOK, "middleprice.html", nil)
//		return
//	}
//
//	roomNumber := c.PostForm("room_number")
//	start := c.PostForm("start")
//	end := c.PostForm("end")
//
//	price, err := h.TransactionalService.FindMiddlePriceForPeriod(c.Request.Context(), roomNumber, start, end)
//	if err != nil {
//		zap.L().Error("FindMiddlePriceForPeriod", zap.Error(err))
//		c.String(http.StatusInternalServerError, err.Error())
//		return
//	}
//	c.String(http.StatusOK, strconv.Itoa(price))
//}

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
