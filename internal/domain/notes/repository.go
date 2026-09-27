package notes

type Repository interface {
	GetByID(id int64) (*Note, error)
	Create(n *Note) error
	List(page, limit int) ([]*NoteListRow, int64, error)
	Delete(id int64) error
}
