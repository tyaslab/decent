package service

import (
	"decent/internal/domain/entity"
)

type BookService interface {
	CreateBook(title, author, description string, publishedYear int) (*entity.Book, error)
	GetBook(id uint) (*entity.Book, error)
	GetAllBooks() ([]entity.Book, error)
	GetBooksByAuthor(author string) ([]entity.Book, error)
	UpdateBook(id uint, title, author string, description *string, publishedYear *int) (*entity.Book, error)
	DeleteBook(id uint) error
	GetBooksPaginated(page, limit int) ([]entity.Book, int64, error)
}
