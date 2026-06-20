package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Channel struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     *int64    `json:"owner_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Post struct {
	ID           int64     `json:"id"`
	ChannelID    int64     `json:"channel_id"`
	ChannelSlug  string    `json:"channel_slug,omitempty"`
	UserID       *int64    `json:"user_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	GuestName    string    `json:"guest_name,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	ImageURLs    []string  `json:"image_urls"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Comment struct {
	ID        int64      `json:"id"`
	PostID    int64      `json:"post_id"`
	ParentID  *int64     `json:"parent_id,omitempty"`
	UserID    *int64     `json:"user_id,omitempty"`
	Username  string     `json:"username,omitempty"`
	GuestName string     `json:"guest_name,omitempty"`
	Content   string     `json:"content"`
	Replies   []*Comment `json:"replies,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
