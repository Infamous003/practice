package postgres

import (
	"errors"
	"practice/internal/domain/notes"

	"gorm.io/gorm"
)

type notesRepo struct {
	db *gorm.DB
}

func NewNotesRepo(db *gorm.DB) *notesRepo {
	return &notesRepo{
		db: db,
	}
}

func (r *notesRepo) GetByID(id int64) (*notes.Note, error) {
	note := &notes.Note{}

	err := r.db.First(note, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, notes.ErrResourceNotFound
		default:
			return nil, err
		}
	}

	return note, nil
}

func (r *notesRepo) Create(note *notes.Note) error {
	return r.db.Create(note).Error
}

func (r *notesRepo) List(page, limit int) ([]*notes.NoteListRow, int64, error) {
	notes := []*notes.NoteListRow{}
	offset := (page - 1) * limit
	var count int64

	res := r.db.Table("notes n").
		Select("n.id, n.title").
		Where("n.deleted_at IS NULL").
		Count(&count).
		Order("n.created_at DESC").
		Limit(limit).
		Offset(offset).Scan(&notes)

	if res.Error != nil {
		return nil, count, res.Error
	}

	return notes, count, nil
}

func (repo *notesRepo) Delete(id int64) error {
	res := repo.db.Delete(&notes.Note{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return notes.ErrResourceNotFound
	}

	return nil
}
