package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"arcalive/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CommentHandler struct{ db *sql.DB }

func NewCommentHandler(db *sql.DB) *CommentHandler { return &CommentHandler{db} }

func (h *CommentHandler) List(c *gin.Context) {
	postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	rows, err := h.db.Query(`
		SELECT c.id, c.post_id, c.parent_id, c.user_id, COALESCE(u.username,''), COALESCE(c.guest_name,''), c.content, c.created_at
		FROM comments c
		LEFT JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at`, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	all := map[int64]*model.Comment{}
	var roots []*model.Comment
	for rows.Next() {
		var cm model.Comment
		rows.Scan(&cm.ID, &cm.PostID, &cm.ParentID, &cm.UserID, &cm.Username, &cm.GuestName, &cm.Content, &cm.CreatedAt)
		cm.Replies = []*model.Comment{}
		all[cm.ID] = &cm
		if cm.ParentID == nil {
			roots = append(roots, &cm)
		}
	}
	for _, cm := range all {
		if cm.ParentID != nil {
			if parent, ok := all[*cm.ParentID]; ok {
				parent.Replies = append(parent.Replies, cm)
			}
		}
	}
	if roots == nil {
		roots = []*model.Comment{}
	}
	c.JSON(http.StatusOK, roots)
}

func (h *CommentHandler) Create(c *gin.Context) {
	postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Content       string `json:"content" binding:"required,min=1"`
		ParentID      *int64 `json:"parent_id"`
		GuestName     string `json:"guest_name"`
		GuestPassword string `json:"guest_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, loggedIn := c.Get("user_id")
	if !loggedIn && (req.GuestName == "" || req.GuestPassword == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "비회원은 닉네임과 비밀번호를 입력하세요"})
		return
	}

	var uid interface{}
	var guestHash string
	if loggedIn {
		uid = userID
	} else {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.GuestPassword), 10)
		guestHash = string(hash)
	}

	var id int64
	h.db.QueryRow(
		`INSERT INTO comments (post_id, parent_id, user_id, guest_name, guest_password, content) VALUES (?,?,?,?,?,?) RETURNING id`,
		postID, req.ParentID, uid, req.GuestName, guestHash, req.Content,
	).Scan(&id)

	h.db.Exec(`UPDATE posts SET comment_count = comment_count + 1 WHERE id = ?`, postID)

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
