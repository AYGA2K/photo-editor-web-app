package models

import guuid "github.com/google/uuid"

type Image struct {
	Name      string     `json:"name"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	Imageid   guuid.UUID `gorm:"primaryKey" json:"imageid"`
	Category  string     `json:"category"`
	UserRefer guuid.UUID `json:"-"`
}
