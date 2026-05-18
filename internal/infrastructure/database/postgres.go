package database

import (
	"decent/internal/config"
	"decent/internal/domain/entity"
	"decent/internal/domain/repository"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDatabase struct {
	db *gorm.DB
}

func NewPostgresDatabase(cfg config.DatabaseConfig) (*PostgresDatabase, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			LogLevel: logger.Info,
		}),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&entity.Book{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &PostgresDatabase{db: db}, nil
}

func (p *PostgresDatabase) BookRepository() repository.BookRepository {
	return NewBookRepository(p.db)
}

func (p *PostgresDatabase) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
