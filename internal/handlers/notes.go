package handlers

import (
	"errors"
	"log"
	"net/http"
	"practice/internal/context"
	"practice/internal/domain/notes"
	"practice/internal/dto"
	"practice/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct {
	service *notes.Service
	l       *log.Logger
}

func NewNotesHandler(service *notes.Service, l *log.Logger) *NotesHandler {
	return &NotesHandler{
		service: service,
		l:       l,
	}
}

func (h *NotesHandler) GetNoteByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticatedUser := context.GetUser(c)
		if authenticatedUser == nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		note, err := h.service.GetNoteByID(id)
		if err != nil {
			switch {
			case errors.Is(err, notes.ErrResourceNotFound):
				response.ErrorResponse(c, http.StatusNotFound, err.Error())
			default:
				h.l.Printf("get note by id failed: %v\n", err)
				response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		resp := dto.Note{
			ID:          int64(note.ID),
			Title:       note.Title,
			Description: note.Description,
			CreatedAt:   note.CreatedAt,
			UpdatedAt:   note.UpdatedAt,
		}
		response.OK(c, resp)
	}
}

func (h *NotesHandler) CreateNoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.NoteRequest

		err := c.ShouldBindJSON(&req)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newNote, err := h.service.Create(&req)
		if err != nil {
			h.l.Printf("create note failed: %v", err)
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		h.l.Println("resource created")
		response.Created(c, newNote)
	}
}

func (h *NotesHandler) ListNotesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if pageInt < 1 || pageInt > 1000 {
			response.ErrorResponse(c, http.StatusBadRequest, "page must be between 1 and 1000")
			return
		}

		if limitInt < 1 || limitInt > 100 {
			response.ErrorResponse(c, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}

		notes, count, err := h.service.List(pageInt, limitInt)
		if err != nil {
			h.l.Printf("list notes failed: %v", err)
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.OK(c, gin.H{"notes": notes, "total": count})
	}
}

func (h *NotesHandler) DeleteNoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid id")
			return
		}
		err = h.service.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, notes.ErrResourceNotFound):
				response.ErrorResponse(c, http.StatusNotFound, "note not found")
			default:
				h.l.Printf("delete note failed: %v", err)
				response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		response.OK(c, "note successfully deleted")
	}
}
