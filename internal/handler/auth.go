package handler

import (
	"database/sql"
	"net/http"
	"time"

	"arcalive/internal/middleware"
	"arcalive/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{ db *sql.DB }

func NewAuthHandler(db *sql.DB) *AuthHandler { return &AuthHandler{db} }

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=20"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	var id int64
	err := h.db.QueryRow(
		`INSERT INTO users (username, email, password_hash) VALUES (?,?,?) RETURNING id`,
		req.Username, req.Email, string(hash),
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "이미 사용 중인 아이디 또는 이메일입니다"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var u model.User
	err := h.db.QueryRow(`SELECT id, username, email, password_hash FROM users WHERE email=?`, req.Email).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "이메일 또는 비밀번호가 틀렸습니다"})
		return
	}
	claims := &middleware.Claims{
		UserID:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(middleware.JwtSecret())
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": u.ID, "username": u.Username, "email": u.Email}})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetInt64("user_id")
	var u model.User
	h.db.QueryRow(`SELECT id, username, email, created_at FROM users WHERE id=?`, uid).
		Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	c.JSON(http.StatusOK, u)
}
