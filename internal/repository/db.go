package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name TEXT NOT NULL,
			encrypted_content MEDIUMBLOB NOT NULL,
			key_check BLOB,
			sort_order INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS project_shares (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			project_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			shared_by BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE KEY uq_project_user (project_id, user_id)
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}

	alterQueries := []string{
		"ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE",
		"ALTER TABLE projects DROP COLUMN has_custom_password",
		"ALTER TABLE projects ADD COLUMN name TEXT NOT NULL",
		"ALTER TABLE projects ADD COLUMN key_check BLOB",
		"ALTER TABLE projects DROP COLUMN encrypted_name",
	}
	for _, q := range alterQueries {
		db.Exec(q)
	}

	return nil
}
