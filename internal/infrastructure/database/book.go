package database

import (
	"decent/internal/domain/entity"
	"decent/internal/domain/repository"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) repository.BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *entity.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) GetByID(id uint) (*entity.Book, error) {
	var book entity.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) GetAll() ([]entity.Book, error) {
	var books []entity.Book
	err := r.db.Find(&books).Error
	return books, err
}

func (r *BookRepository) GetByAuthor(author string) ([]entity.Book, error) {
	var books []entity.Book
	err := r.db.Where("author = ?", author).Find(&books).Error
	return books, err
}

func (r *BookRepository) Update(book *entity.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Book{}, id).Error
}

func (r *BookRepository) GetAllPaginated(page, limit int) ([]entity.Book, int64, error) {
	var books []entity.Book
	var total int64

	if err := r.db.Model(&entity.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&books).Error
	return books, total, err
}