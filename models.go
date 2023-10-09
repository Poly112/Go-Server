package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID        uuid.UUID `gorm:"size:255;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	ApiKey    string    `gorm:"size:255;not null;unique;" json:"api_key"`
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ApiKey == "" {
		u.ApiKey = generateRandomHexValue()
	}
	return
}

type Feeds struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	Name      string    `gorm:";not null" json:"name"`
	Url       string    `gorm:";not null;unique" json:"url"`
	UserID    uuid.UUID `gorm:"size:255;not null;references:Users(id) ;onUpdate:CASCADE,onDelete:CASCADE" json:"user_id"`
}
type FeedFollows struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	UserID    uuid.UUID `gorm:"size:255;not null;references:Users(id);onUpdate:CASCADE,onDelete:CASCADE;uniqueIndex:user_feed_idx" json:"user_id"`
	FeedID    uuid.UUID `gorm:"size:255;not null;references:Feeds(id);onUpdate:CASCADE,onDelete:CASCADE;uniqueIndex:user_feed_idx" json:"feed_id"`
}

type Posts struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
	Title       string    `gorm:"not null" json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `gorm:"not null" json:"published_at"`
	Url         string    `gorm:"not null; unique" json:"url"`
	FeedID      uuid.UUID `gorm:"not null;references:Feeds(id);onUpdate:CASCADE,onDelete:CASCADE" json:"feed_id"`
}
