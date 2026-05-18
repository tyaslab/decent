package usecase

import (
	"errors"
	"decent/internal/domain/entity"
	"decent/internal/domain/repository"
	"decent/internal/domain/service"
)

type BookUseCase struct {
	repo repository.BookRepository
}

func NewBookUseCase(repo repository.BookRepository) service.BookService {
	return &BookUseCase{repo: repo}
}

func (uc *BookUseCase) CreateBook(title, author, description string, publishedYear int) (*entity.Book, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if author == "" {
		return nil, errors.New("author is required")
	}
	if publishedYear < 1000 || publishedYear > 9999 {
		return nil, errors.New("published year must be between 1000 and 9999")
	}

	book := &entity.Book{
		Title:         title,
		Author:        author,
		Description:   description,
		PublishedYear: publishedYear,
	}

	if err := uc.repo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (uc *BookUseCase) GetBook(id uint) (*entity.Book, error) {
	return uc.repo.GetByID(id)
}

func (uc *BookUseCase) GetAllBooks() ([]entity.Book, error) {
	return uc.repo.GetAll()
}

func (uc *BookUseCase) GetBooksByAuthor(author string) ([]entity.Book, error) {
	return uc.repo.GetByAuthor(author)
}

func (uc *BookUseCase) UpdateBook(id uint, title, author string, description *string, publishedYear *int) (*entity.Book, error) {
	book, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		book.Title = title
	}
	if author != "" {
		book.Author = author
	}
	if description != nil {
		book.Description = *description
	}
	if publishedYear != nil {
		book.PublishedYear = *publishedYear
	}

	if err := uc.repo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (uc *BookUseCase) DeleteBook(id uint) error {
	return uc.repo.Delete(id)
}

func (uc *BookUseCase) GetBooksPaginated(page, limit int) ([]entity.Book, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return uc.repo.GetAllPaginated(page, limit)
}