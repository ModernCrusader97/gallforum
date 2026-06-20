package handler

import (
	"database/sql"
	"net/http"
	"regexp"

	"arcalive/internal/model"

	"github.com/gin-gonic/gin"
)

type ChannelHandler struct{ db *sql.DB }

func NewChannelHandler(db *sql.DB) *ChannelHandler { return &ChannelHandler{db} }

var slugRe = regexp.MustCompile(`^[a-z0-9\-]{2,30}$`)

func (h *ChannelHandler) List(c *gin.Context) {
	rows, err := h.db.Query(`SELECT id, slug, name, description, owner_id, created_at FROM channels ORDER BY name`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var channels []model.Channel
	for rows.Next() {
		var ch model.Channel
		rows.Scan(&ch.ID, &ch.Slug, &ch.Name, &ch.Description, &ch.OwnerID, &ch.CreatedAt)
		channels = append(channels, ch)
	}
	if channels == nil {
		channels = []model.Channel{}
	}
	c.JSON(http.StatusOK, channels)
}

func (h *ChannelHandler) Get(c *gin.Context) {
	var ch model.Channel
	err := h.db.QueryRow(`SELECT id, slug, name, description, owner_id, created_at FROM channels WHERE slug=?`, c.Param("slug")).
		Scan(&ch.ID, &ch.Slug, &ch.Name, &ch.Description, &ch.OwnerID, &ch.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "채널을 찾을 수 없습니다"})
		return
	}
	c.JSON(http.StatusOK, ch)
}

func (h *ChannelHandler) Create(c *gin.Context) {
	var req struct {
		Slug        string `json:"slug" binding:"required"`
		Name        string `json:"name" binding:"required,min=1,max=30"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !slugRe.MatchString(req.Slug) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "슬러그는 소문자, 숫자, 하이픈만 사용 가능합니다"})
		return
	}
	uid, _ := c.Get("user_id")
	var id int64
	err := h.db.QueryRow(
		`INSERT INTO channels (slug, name, description, owner_id) VALUES (?,?,?,?) RETURNING id`,
		req.Slug, req.Name, req.Description, uid,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "이미 존재하는 채널 슬러그입니다"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "slug": req.Slug})
}
