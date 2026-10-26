package server

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/services"
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var ownerID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type fakeApartment struct {
	room   string
	active bool
	roles  []string
	users  []uuid.UUID
}

// fakeAccessStorage повторяет правило видимости из repo.apartmentInScope.
type fakeAccessStorage struct {
	apartments []fakeApartment
	rooms      map[int]string // reservationID -> апартамент
}

func (f fakeAccessStorage) ReadAllowedRooms(_ context.Context, scope entities.ApartmentScope) ([]string, error) {
	var rooms []string
	for _, a := range f.apartments {
		granted := slices.Contains(a.roles, scope.Role) || slices.Contains(a.users, scope.UserID)
		if (scope.Inactive || a.active) && (scope.All || granted) {
			rooms = append(rooms, a.room)
		}
	}
	return rooms, nil
}

func (f fakeAccessStorage) ReadRoomByReservationID(_ context.Context, id int) (string, bool, error) {
	room, ok := f.rooms[id]
	return room, ok, nil
}

func (f fakeAccessStorage) ReadRoomByReservationInfoID(_ context.Context, id int) (string, bool, error) {
	return f.ReadRoomByReservationID(context.Background(), id)
}

var testUsers = map[string]entities.CurrentUser{
	"admin": {Role: "admin", Permissions: entities.NewPermissions(entities.PermCalendarRead, entities.PermUsersManage,
		entities.PermApartmentsAll, entities.PermApartmentsInactive)},
	"manager": {Role: "manager", Permissions: entities.NewPermissions(entities.PermCalendarRead)},
	"owner":   {ID: ownerID, Role: "apartmentOwner", Permissions: entities.NewPermissions(entities.PermCalendarRead)},
}

// newTestRouter — роутер с пользователем user (пустая строка — без входа).
func newTestRouter(user string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := &Handle{access: services.NewServiceAccess(fakeAccessStorage{
		apartments: []fakeApartment{
			{room: "INACT1", active: false, roles: []string{"manager"}},
			{room: "OPEN1", active: true, roles: []string{"manager"}},
			{room: "OWN1", active: true, users: []uuid.UUID{ownerID}},
		},
		rooms: map[int]string{1: "INACT1", 2: "OPEN1", 3: "OWN1"},
	})}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		if u, ok := testUsers[user]; ok {
			c.Set(currentUserContextKey, u)
		}
	})
	r.POST("/createUser", RequirePermission(entities.PermUsersManage), func(c *gin.Context) { c.Status(http.StatusCreated) })
	r.GET("/both", RequirePermission(entities.PermCalendarRead, entities.PermUsersManage), func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/any", RequireAnyPermission(entities.PermUsersManage, entities.PermRolesManage), func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/room/:room", func(c *gin.Context) {
		if h.requireRooms(c, c.Param("room")) {
			c.Status(http.StatusOK)
		}
	})
	r.GET("/reservation/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if h.requireReservation(c, id) {
			c.Status(http.StatusOK)
		}
	})
	r.GET("/filter", func(c *gin.Context) {
		access, ok := h.roomAccess(c)
		if ok {
			c.JSON(http.StatusOK, filterByRoom(access, []string{"INACT1", "OPEN1", "OWN1", "UNKNOWN"}, func(s string) string { return s }))
		}
	})
	return r
}

func do(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}

func TestRequirePermission(t *testing.T) {
	assert.Equal(t, http.StatusCreated, do(newTestRouter("admin"), http.MethodPost, "/createUser").Code)
	assert.Equal(t, http.StatusForbidden, do(newTestRouter("manager"), http.MethodPost, "/createUser").Code)
	assert.Equal(t, http.StatusUnauthorized, do(newTestRouter(""), http.MethodPost, "/createUser").Code)

	// нужны все перечисленные права
	assert.Equal(t, http.StatusOK, do(newTestRouter("admin"), http.MethodGet, "/both").Code)
	assert.Equal(t, http.StatusForbidden, do(newTestRouter("manager"), http.MethodGet, "/both").Code)
}

func TestRequireAnyPermission(t *testing.T) {
	assert.Equal(t, http.StatusOK, do(newTestRouter("admin"), http.MethodGet, "/any").Code)
	assert.Equal(t, http.StatusForbidden, do(newTestRouter("manager"), http.MethodGet, "/any").Code)
	assert.Equal(t, http.StatusUnauthorized, do(newTestRouter(""), http.MethodGet, "/any").Code)
}

func TestRequireRooms(t *testing.T) {
	manager := newTestRouter("manager")
	assert.Equal(t, http.StatusOK, do(manager, http.MethodGet, "/room/OPEN1").Code)
	assert.Equal(t, http.StatusForbidden, do(manager, http.MethodGet, "/room/INACT1").Code, "неактивный")
	assert.Equal(t, http.StatusForbidden, do(manager, http.MethodGet, "/room/OWN1").Code, "не выдан роли")
	assert.Equal(t, http.StatusForbidden, do(manager, http.MethodGet, "/room/UNKNOWN").Code)

	admin := newTestRouter("admin")
	assert.Equal(t, http.StatusOK, do(admin, http.MethodGet, "/room/INACT1").Code)
	assert.Equal(t, http.StatusOK, do(admin, http.MethodGet, "/room/OWN1").Code)

	// владельцу доступны только выданные лично ему
	owner := newTestRouter("owner")
	assert.Equal(t, http.StatusOK, do(owner, http.MethodGet, "/room/OWN1").Code)
	assert.Equal(t, http.StatusForbidden, do(owner, http.MethodGet, "/room/OPEN1").Code)
}

func TestRequireReservation(t *testing.T) {
	manager := newTestRouter("manager")
	assert.Equal(t, http.StatusForbidden, do(manager, http.MethodGet, "/reservation/1").Code)
	assert.Equal(t, http.StatusOK, do(manager, http.MethodGet, "/reservation/2").Code)
	// несуществующая бронь пропускается — хендлер сам ответит "не найдено"
	assert.Equal(t, http.StatusOK, do(manager, http.MethodGet, "/reservation/999").Code)

	owner := newTestRouter("owner")
	assert.Equal(t, http.StatusOK, do(owner, http.MethodGet, "/reservation/3").Code)
	assert.Equal(t, http.StatusForbidden, do(owner, http.MethodGet, "/reservation/2").Code)
}

func TestFilterByRoom(t *testing.T) {
	assert.JSONEq(t, `["OPEN1"]`, do(newTestRouter("manager"), http.MethodGet, "/filter").Body.String())
	assert.JSONEq(t, `["INACT1","OPEN1","OWN1"]`, do(newTestRouter("admin"), http.MethodGet, "/filter").Body.String())
	assert.JSONEq(t, `["OWN1"]`, do(newTestRouter("owner"), http.MethodGet, "/filter").Body.String())
}
