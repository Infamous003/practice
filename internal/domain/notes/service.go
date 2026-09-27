package notes

import "practice/internal/dto"

type Service struct {
	notesRepo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		notesRepo: repo,
	}
}

func (s *Service) GetNoteByID(id int64) (*dto.Note, error) {
	note, err := s.notesRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.Note{
		ID:          int64(note.ID),
		Title:       note.Title,
		Description: note.Description,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}, nil
}

func (s *Service) Create(req *dto.NoteRequest) (*dto.Note, error) {
	newNote := &Note{
		Title:       req.Title,
		Description: req.Description,
	}

	err := s.notesRepo.Create(newNote)
	if err != nil {
		return nil, err
	}

	return &dto.Note{
		ID:          int64(newNote.ID),
		Title:       newNote.Title,
		Description: newNote.Description,
		CreatedAt:   newNote.CreatedAt,
		UpdatedAt:   newNote.UpdatedAt,
	}, nil
}

func (s *Service) List(page, limit int) ([]*NoteListRow, int64, error) {
	notesList, count, err := s.notesRepo.List(page, limit)
	if err != nil {
		return nil, 0, err
	}

	notesResponse := []*dto.NoteListRow{}

	for _, n := range notesList {
		notesResponse = append(notesResponse, &dto.NoteListRow{
			ID:    n.ID,
			Title: n.Title,
		})
	}

	return notesList, count, nil
}

func (s *Service) Delete(id int64) error {
	return s.notesRepo.Delete(id)
}
