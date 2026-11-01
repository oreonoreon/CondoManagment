package server

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"awesomeProject/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Админка: /calendar/admin/... Права на каждую группу ручек — в Gin().

// adminError пишет ответ для ошибки админки: тексты доменных ошибок показываются в интерфейсе.
func adminError(c *gin.Context, err error) {
	var e *erro.Err
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, erro.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, erro.ErrAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, erro.ErrInvalidReference), errors.Is(err, erro.ErrUnknownRole):
		status = http.StatusBadRequest
	case errors.As(err, &e) && e.Code == erro.CodeValidation:
		status = http.StatusBadRequest
	}
	if status == http.StatusInternalServerError {
		zap.L().Error("admin", zap.String("path", c.Request.URL.Path), zap.Error(err))
		c.JSON(status, gin.H{"error": "internal error", "message": "Внутренняя ошибка сервера"})
		return
	}
	message := err.Error()
	if errors.As(err, &e) {
		message = e.Message()
	}
	c.JSON(status, gin.H{"error": message, "message": message})
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad payload", "message": "Некорректный запрос: " + err.Error()})
		return false
	}
	return true
}

func uuidParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id", "message": "Некорректный id"})
		return uuid.Nil, false
	}
	return id, true
}

func intParam(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id", "message": "Некорректный id"})
		return 0, false
	}
	return id, true
}

// ─── Пользователи ───

// AdminListUsers GET /calendar/admin/users
func (h *Handle) AdminListUsers(c *gin.Context) {
	users, err := h.admin.ListUsers(c.Request.Context())
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

type createUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

// AdminCreateUser POST /calendar/admin/users
func (h *Handle) AdminCreateUser(c *gin.Context) {
	var req createUserRequest
	if !bindJSON(c, &req) {
		return
	}
	user, err := h.admin.CreateUser(c.Request.Context(), req.Name, req.Password, req.Phone, req.Role)
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// AdminUpdateUser PATCH /calendar/admin/users/:id — поля name, phone, role, password (любые из них).
func (h *Handle) AdminUpdateUser(c *gin.Context) {
	current, ok := mustCurrentUser(c)
	if !ok {
		return
	}
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var req services.UserChanges
	if !bindJSON(c, &req) {
		return
	}
	user, err := h.admin.UpdateUser(c.Request.Context(), current, id, req)
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// AdminDeleteUser DELETE /calendar/admin/users/:id
func (h *Handle) AdminDeleteUser(c *gin.Context) {
	current, ok := mustCurrentUser(c)
	if !ok {
		return
	}
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := h.admin.DeleteUser(c.Request.Context(), current, id); err != nil {
		adminError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ─── Апартаменты ───

// AdminListApartments GET /calendar/admin/apartments — все, включая неактивные, с выданным доступом.
func (h *Handle) AdminListApartments(c *gin.Context) {
	apartments, err := h.admin.ListApartments(c.Request.Context())
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, apartments)
}

type createApartmentRequest struct {
	RoomNumber     string      `json:"room_number"`
	Description    string      `json:"description"`
	AirbnbCalendar string      `json:"airbnbCalendar"`
	IsActive       *bool       `json:"is_active"` // по умолчанию true
	Roles          []string    `json:"roles"`
	UserIDs        []uuid.UUID `json:"user_ids"`
}

// AdminCreateApartment POST /calendar/admin/apartments
func (h *Handle) AdminCreateApartment(c *gin.Context) {
	var req createApartmentRequest
	if !bindJSON(c, &req) {
		return
	}
	apartment := entities.Apartment{
		RoomNumber:     req.RoomNumber,
		Description:    req.Description,
		AirbnbCalendar: req.AirbnbCalendar,
		IsActive:       req.IsActive == nil || *req.IsActive,
	}
	created, err := h.admin.CreateApartment(c.Request.Context(), apartment,
		entities.ApartmentAccess{Roles: req.Roles, UserIDs: req.UserIDs})
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

// AdminUpdateApartment PATCH /calendar/admin/apartments/:id — поля description, airbnbCalendar, is_active.
func (h *Handle) AdminUpdateApartment(c *gin.Context) {
	id, ok := intParam(c, "id")
	if !ok {
		return
	}
	var req entities.ApartmentUpdate
	if !bindJSON(c, &req) {
		return
	}
	apartment, err := h.admin.UpdateApartment(c.Request.Context(), id, req)
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, apartment)
}

// AdminSetApartmentAccess PUT /calendar/admin/apartments/:id/access — заменяет весь доступ: {roles, user_ids}.
func (h *Handle) AdminSetApartmentAccess(c *gin.Context) {
	id, ok := intParam(c, "id")
	if !ok {
		return
	}
	var req entities.ApartmentAccess
	if !bindJSON(c, &req) {
		return
	}
	apartment, err := h.admin.SetApartmentAccess(c.Request.Context(), id, req)
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, apartment)
}

// ─── Роли ───

// AdminListRoles GET /calendar/admin/roles
func (h *Handle) AdminListRoles(c *gin.Context) {
	roles, err := h.admin.ListRoles(c.Request.Context())
	if err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// AdminPermissionCatalog GET /calendar/admin/permissions — все права, которые можно выдать роли.
func (h *Handle) AdminPermissionCatalog(c *gin.Context) {
	c.JSON(http.StatusOK, h.admin.PermissionCatalog())
}

type createRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// AdminCreateRole POST /calendar/admin/roles
func (h *Handle) AdminCreateRole(c *gin.Context) {
	var req createRoleRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.admin.CreateRole(c.Request.Context(), req.Name, req.Description, req.Permissions); err != nil {
		adminError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

// AdminUpdateRole PATCH /calendar/admin/roles/:name — description и/или permissions (заменяются целиком).
func (h *Handle) AdminUpdateRole(c *gin.Context) {
	current, ok := mustCurrentUser(c)
	if !ok {
		return
	}
	var req services.RoleChanges
	if !bindJSON(c, &req) {
		return
	}
	if err := h.admin.UpdateRole(c.Request.Context(), current, c.Param("name"), req); err != nil {
		adminError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// AdminDeleteRole DELETE /calendar/admin/roles/:name
func (h *Handle) AdminDeleteRole(c *gin.Context) {
	current, ok := mustCurrentUser(c)
	if !ok {
		return
	}
	if err := h.admin.DeleteRole(c.Request.Context(), current, c.Param("name")); err != nil {
		adminError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
