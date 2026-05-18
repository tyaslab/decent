package repository

import (
	"decent/internal/domain/entity"
)

type BookRepository interface {
	Create(book *entity.Book) error
	GetByID(id uint) (*entity.Book, error)
	GetAll() ([]entity.Book, error)
	GetByAuthor(author string) ([]entity.Book, error)
	Update(book *entity.Book) error
	Delete(id uint) error
	GetAllPaginated(page, limit int) ([]entity.Book, int64, error)
}
