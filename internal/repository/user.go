package repository

import (
	"database/sql"
	"fmt"

	"github.com/user/access-storage-server/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateWithAdmin(email, passwordHash string, isAdmin bool) (*model.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (email, password_hash, is_admin) VALUES (?, ?, ?)",
		email, passwordHash, isAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	id, _ := result.LastInsertId()
	return r.GetByID(id)
}

func (r *UserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		"SELECT id, email, password_hash, is_admin, created_at FROM users WHERE email = ?",
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		"SELECT id, email, password_hash, is_admin, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) ListAll() ([]model.User, error) {
	rows, err := r.db.Query(
		"SELECT id, email, password_hash, is_admin, created_at FROM users ORDER BY id ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) UpdateAdmin(id int64, isAdmin bool) error {
	_, err := r.db.Exec("UPDATE users SET is_admin = ? WHERE id = ?", isAdmin, id)
	if err != nil {
		return fmt.Errorf("update user admin: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateProfile(id int64, email, passwordHash string, isAdmin bool) error {
	if passwordHash != "" {
		_, err := r.db.Exec(
			"UPDATE users SET email = ?, password_hash = ?, is_admin = ? WHERE id = ?",
			email, passwordHash, isAdmin, id,
		)
		if err != nil {
			return fmt.Errorf("update user profile: %w", err)
		}
	} else {
		_, err := r.db.Exec(
			"UPDATE users SET email = ?, is_admin = ? WHERE id = ?",
			email, isAdmin, id,
		)
		if err != nil {
			return fmt.Errorf("update user profile: %w", err)
		}
	}
	return nil
}

func (r *UserRepository) CountAdmins() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = TRUE").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return count, nil
}

func (r *UserRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
