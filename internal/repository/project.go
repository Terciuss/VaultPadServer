package repository

import (
	"database/sql"
	"fmt"

	"github.com/user/access-storage-server/internal/model"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) ListByUser(userID int64) ([]model.Project, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, encrypted_name, encrypted_content,
		        sort_order, created_at, updated_at
		 FROM projects WHERE user_id = ? ORDER BY sort_order ASC, created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.EncryptedName, &p.EncryptedContent,
			&p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) GetByID(id, userID int64) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(
		`SELECT id, user_id, encrypted_name, encrypted_content,
		        sort_order, created_at, updated_at
		 FROM projects WHERE id = ? AND user_id = ?`,
		id, userID,
	).Scan(
		&p.ID, &p.UserID, &p.EncryptedName, &p.EncryptedContent,
		&p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return &p, nil
}

func (r *ProjectRepository) Create(p *model.Project) (*model.Project, error) {
	result, err := r.db.Exec(
		`INSERT INTO projects (user_id, encrypted_name, encrypted_content, sort_order)
		 VALUES (?, ?, ?, ?)`,
		p.UserID, p.EncryptedName, p.EncryptedContent, p.SortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	id, _ := result.LastInsertId()
	return r.GetByID(id, p.UserID)
}

func (r *ProjectRepository) Update(p *model.Project) error {
	_, err := r.db.Exec(
		`UPDATE projects SET encrypted_name = ?, encrypted_content = ?,
		        sort_order = ?
		 WHERE id = ? AND user_id = ?`,
		p.EncryptedName, p.EncryptedContent,
		p.SortOrder, p.ID, p.UserID,
	)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) ListMetaByUser(userID int64) ([]model.ProjectMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, updated_at FROM projects WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list project meta: %w", err)
	}
	defer rows.Close()

	var metas []model.ProjectMeta
	for rows.Next() {
		var m model.ProjectMeta
		if err := rows.Scan(&m.ID, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project meta: %w", err)
		}
		metas = append(metas, m)
	}
	return metas, rows.Err()
}

func (r *ProjectRepository) Delete(id, userID int64) error {
	result, err := r.db.Exec(
		"DELETE FROM projects WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

func (r *ProjectRepository) ListAll() ([]model.Project, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, encrypted_name, encrypted_content,
		        sort_order, created_at, updated_at
		 FROM projects ORDER BY sort_order ASC, created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list all projects: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.EncryptedName, &p.EncryptedContent,
			&p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) ListAllMeta() ([]model.ProjectMeta, error) {
	rows, err := r.db.Query(`SELECT id, updated_at FROM projects`)
	if err != nil {
		return nil, fmt.Errorf("list all project meta: %w", err)
	}
	defer rows.Close()

	var metas []model.ProjectMeta
	for rows.Next() {
		var m model.ProjectMeta
		if err := rows.Scan(&m.ID, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project meta: %w", err)
		}
		metas = append(metas, m)
	}
	return metas, rows.Err()
}

func (r *ProjectRepository) GetByIDRaw(id int64) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(
		`SELECT id, user_id, encrypted_name, encrypted_content,
		        sort_order, created_at, updated_at
		 FROM projects WHERE id = ?`,
		id,
	).Scan(
		&p.ID, &p.UserID, &p.EncryptedName, &p.EncryptedContent,
		&p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return &p, nil
}

func (r *ProjectRepository) UpdateRaw(p *model.Project) error {
	_, err := r.db.Exec(
		`UPDATE projects SET encrypted_name = ?, encrypted_content = ?,
		        sort_order = ?
		 WHERE id = ?`,
		p.EncryptedName, p.EncryptedContent,
		p.SortOrder, p.ID,
	)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) DeleteRaw(id int64) error {
	result, err := r.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}
