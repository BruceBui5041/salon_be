package models

type Tag struct {
	TagID  uint    `gorm:"primaryKey;autoIncrement"`
	Name   string  `gorm:"uniqueIndex;not null;size:50"`
	Videos []Video `gorm:"many2many:video_tags;"`
}
