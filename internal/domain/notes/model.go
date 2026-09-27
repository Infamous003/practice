package notes

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
}

type NoteListRow struct {
	ID    int64
	Title string
}
