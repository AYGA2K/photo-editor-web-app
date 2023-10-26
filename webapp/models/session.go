package models

import (
	"time"

	guuid "github.com/google/uuid"
)

type Session struct {
	Expires   time.Time  `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	Sessionid guuid.UUID `gorm:"primaryKey" json:"sessionid"`
	UserRefer guuid.UUID `json:"-"`
}
