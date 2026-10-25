package server

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/services"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var ownerID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type fakeAccessStorage struct {
	forbidden map[string][]string    // role -> закрытые апартаменты
	owned     map[uuid.UUID][]string // userID -> свои апартаменты
	rooms     map[int]string         // reservationID -> апартамент
}

func (f fakeAccessStorage) ReadForbiddenRooms(_ context.Context, role string) ([]string, error) {
	return f.forbidden[role], nil
}

func (f fakeAccessStorage) ReadOwnedRooms(_ context.Context, userID uuid.UUID) ([]string, error) {
	return f.owned[userID], nil
}

func (f fakeAccessStorage) ReadRoomByReservationID(_ context.Context, id int) (string, bool, error) {
	room, ok := f.rooms[id]
	return room, ok, nil
}

func (f fakeAccessStorage) ReadRoomByReservationInfoID(_ context.Context, id int) (string, bool, error) {
	return f.ReadRoomByReservationID(context.Background(), id)
}

var testUsers = map[string]entities.CurrentUser{
	"admin":   {Role: "admin", Permissions: entities.NewPermissions(entities.PermCalendarRead, entities.PermUsersManage)},
	"manager": {Role: "manager", Permissions: entities.NewPermissions(entities.PermCalendarRead)},
	"owner": {ID: ownerID, Role: "apartmentOwner",
		Permissions: entities.NewPermissions(entities.PermCalendarRead, entities.PermApartmentsOwnedOnly)},
}

// newTestRouter — роутер с пользователем user (пустая строка — без входа).
func newTestRouter(user string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := &Handle{access: services.NewServiceAccess(fakeAccessStorage{
		forbidden: map[string][]string{"manager": {"ADMIN1"}},
		owned:     map[uuid.UUID][]string{ownerID: {"OWN1"}},
		rooms:     map[int]string{1: "ADMIN1", 2: "OPEN1", 3: "OWN1"},
	})}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		if u, ok := testUsers[user]; ok {
			c.Set(currentUserContextKey, u)
		}
	})
	r.POST("/createUser", RequirePermission(entities.PermUsersManage), func(c *gin.Context) { c.Status(http.StatusCreated) })
	r.GET("/both", RequirePermission(entities.PermCalendarRead, entities.PermUsersManage), func(c *gin.Context) { c.Status(http.StatusOK) })
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
			c.JSON(http.StatusOK, filterByRoom(access, []string{"ADMIN1", "OPEN1", "OWN1", "UNKNOWN"}, func(s string) string { return s }))
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

func TestRequireRooms(t *testing.T) {
	manager := newTestRouter("manager")
	assert.Equal(t, http.StatusForbidden, do(manager, http.MethodGet, "/room/ADMIN1").Code)
	assert.Equal(t, http.StatusOK, do(manager, http.MethodGet, "/room/OPEN1").Code)
	assert.Equal(t, http.StatusOK, do(manager, http.MethodGet, "/room/UNKNOWN").Code)

	assert.Equal(t, http.StatusOK, do(newTestRouter("admin"), http.MethodGet, "/room/ADMIN1").Code)

	// владельцу доступны только свои апартаменты
	owner := newTestRouter("owner")
	assert.Equal(t, http.StatusOK, do(owner, http.MethodGet, "/room/OWN1").Code)
	assert.Equal(t, http.StatusForbidden, do(owner, http.MethodGet, "/room/OPEN1").Code)
	assert.Equal(t, http.StatusForbidden, do(owner, http.MethodGet, "/room/UNKNOWN").Code)
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
	assert.JSONEq(t, `["OPEN1","OWN1","UNKNOWN"]`, do(newTestRouter("manager"), http.MethodGet, "/filter").Body.String())
	assert.JSONEq(t, `["ADMIN1","OPEN1","OWN1","UNKNOWN"]`, do(newTestRouter("admin"), http.MethodGet, "/filter").Body.String())
	assert.JSONEq(t, `["OWN1"]`, do(newTestRouter("owner"), http.MethodGet, "/filter").Body.String())
}
