package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/PhosFactum/effective-mobile-test-go/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func RunMigrations(db *sql.DB, cfg *config.Config) error {
	// Проверяем существование таблицы migrations
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (version INT PRIMARY KEY)`)
	if err != nil {
		return err
	}

	// Простые миграции - применяем 001_init.up.sql
	var version int
	err = db.QueryRow("SELECT version FROM migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if version == 0 {
		migrationPath := filepath.Join("migrations", "001_init.up.sql")
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO migrations (version) VALUES (1)")
		if err != nil {
			return err
		}

		log.Println("Migrations applied successfully")
	}

	return nil
}
