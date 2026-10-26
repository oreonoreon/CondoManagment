package server

import (
	"awesomeProject/internal/erro"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// CreateUser POST /calendar/createUser — старая ручка создания пользователя (новая — POST /calendar/admin/users).
func (h *Handle) CreateUser(c *gin.Context) {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Role     string `json:"role"` // необязательно, по умолчанию manager
	}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad payload"})
		return
	}

	createdUser, err := h.admin.CreateUser(c.Request.Context(), u.Username, u.Password, u.Phone, u.Role)
	if err != nil {
		adminError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"userID":  createdUser.ID,
	})
}

func (h *Handle) LoginHandler(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad payload"})
		return
	}

	user, err := h.GetUser(c.Request.Context(), creds.Username)
	if err != nil {
		if errors.Is(err, erro.ErrWrongCreds) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Роль в сессии не храним: SessionAuthMiddleware на каждый запрос читает актуальные роль и права из БД.
	session := sessions.Default(c)
	session.Set("userID", user.ID.String())

	if err := session.Save(); err != nil {
		zap.L().Error("LoginHandler: failed to save session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save session"})
		return
	}

	zap.L().Info("Session saved successfully", zap.String("username", creds.Username))
	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (h *Handle) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()

	// Настройки должны совпадать с теми, что использовались при создании
	// куки в Gin() — иначе браузер не сможет корректно её удалить.
	cookieOptions := sessions.Options{
		Path:     "/",
		MaxAge:   -1, // удаляем cookie
		HttpOnly: true,
		Secure:   h.cfg.IsProduction,
		SameSite: http.SameSiteLaxMode,
	}

	session.Options(cookieOptions)

	if err := session.Save(); err != nil {
		zap.L().Error("LogoutHandler: failed to save session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// SessionAuthMiddleware пускает только вошедших пользователей и кладёт в контекст
// актуальные из БД роль и права (смена роли действует сразу, удалённый пользователь теряет доступ).
func (h *Handle) SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid, _ := session.Get("userID").(string)

		id, err := uuid.Parse(uid)
		if err != nil {
			zap.L().Warn("Unauthorized access attempt", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		user, err := h.ServiceUsers.GetCurrentUser(c.Request.Context(), id)
		if err != nil {
			zap.L().Error("SessionAuthMiddleware: GetCurrentUser", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if user == nil {
			zap.L().Warn("Session of deleted user", zap.String("userID", uid))
			session.Clear()
			if err := session.Save(); err != nil {
				zap.L().Error("SessionAuthMiddleware: failed to clear session", zap.Error(err))
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		c.Set(currentUserContextKey, *user)
		c.Set("userID", uid)
		c.Next()
	}
}

// Me GET /calendar/me — текущий пользователь, его роль и права (фронт по ним строит интерфейс).
func (h *Handle) Me(c *gin.Context) {
	user, ok := mustCurrentUser(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"name":        user.Name,
		"role":        user.Role,
		"permissions": user.Permissions.List(),
	})
}
