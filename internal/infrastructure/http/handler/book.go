package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"decent/internal/application/dto"
	"decent/internal/domain/service"
	"decent/internal/infrastructure/auth"

	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
)

type BookHandler struct {
	bookService service.BookService
	jwtService  *auth.JWTService
	validator   *validator.Validate
}

func NewBookHandler(bookService service.BookService, jwtService *auth.JWTService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		jwtService:  jwtService,
		validator:   validator.New(),
	}
}

func (h *BookHandler) CreateBook(c echo.Context) error {
	var req dto.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	book, err := h.bookService.CreateBook(req.Title, req.Author, req.Description, req.PublishedYear)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBook(c echo.Context) error {
	id := c.Param("id")
	var bookID uint
	if _, err := fmt.Sscanf(id, "%d", &bookID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book ID"})
	}

	book, err := h.bookService.GetBook(bookID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "book not found"})
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetAllBooks(c echo.Context) error {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBooksByAuthor(c echo.Context) error {
	author := c.QueryParam("author")
	if author == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "author parameter is required"})
	}

	books, err := h.bookService.GetBooksByAuthor(author)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBooksPaginated(c echo.Context) error {
	page := 1
	limit := 10

	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	books, total, err := h.bookService.GetBooksPaginated(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"books": books,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	id := c.Param("id")
	var bookID uint
	if _, err := fmt.Sscanf(id, "%d", &bookID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book ID"})
	}

	var req dto.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	title := ""
	author := ""
	description := req.Description
	publishedYear := req.PublishedYear

	if req.Title != nil {
		title = *req.Title
	}
	if req.Author != nil {
		author = *req.Author
	}

	book, err := h.bookService.UpdateBook(bookID, title, author, description, publishedYear)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	id := c.Param("id")
	var bookID uint
	if _, err := fmt.Sscanf(id, "%d", &bookID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book ID"})
	}

	if err := h.bookService.DeleteBook(bookID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}