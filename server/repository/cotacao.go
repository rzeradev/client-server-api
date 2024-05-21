package repository

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	ID        uint `gorm:"primaryKey"`
	Bid       string
	Timestamp time.Time
}

type Repository struct {
	DB *gorm.DB
}

func InitDB() (*Repository, error) {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		if err := os.Mkdir("./db", os.ModePerm); err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open("./db/cotacoes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Cotacao{})

	return &Repository{DB: db}, nil
}

func (r *Repository) SaveCotacao(ctx context.Context, bid string) error {
	select {
	case <-ctx.Done():
		log.Println("Context Timeout exceeded for database operation (10 ms)")
		logToFile("Context Timeout exceeded for database operation (10 ms)")
		return ctx.Err()
	default:
		cotacao := Cotacao{Bid: bid, Timestamp: time.Now()}
		result := r.DB.WithContext(ctx).Create(&cotacao)
		return result.Error
	}
}

func logToFile(message string) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(message)
}
