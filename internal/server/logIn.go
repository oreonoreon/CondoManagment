package server

import (
	"awesomeProject/internal/erro"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (h *Handle) CreateUser(c *gin.Context) {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad payload"})
		return
	}

	user, err := h.ServiceUsers.PrepareToCreateUser(u.Username, u.Password, u.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.ServiceUsers.CreateUser(c.Request.Context(), *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	// ИСПРАВЛЕНИЕ: используем те же настройки, что и при создании
	cookieOptions := sessions.Options{
		Path:     "/",
		MaxAge:   -1, // удаляем cookie
		HttpOnly: true,
		Secure:   h.cfg.IsProduction,
		SameSite: http.SameSiteLaxMode,
	}

	if h.cfg.IsProduction {
		cookieOptions.SameSite = http.SameSiteNoneMode
	}

	session.Options(cookieOptions)

	if err := session.Save(); err != nil {
		zap.L().Error("LogoutHandler: failed to save session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("userID")

		if uid == nil {
			zap.L().Warn("Unauthorized access attempt", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		c.Set("userID", uid)
		c.Next()
	}
}
