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

func (r *ProjectRepository) GetByID(id int64) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(
		`SELECT id, name, encrypted_content, key_check,
		        sort_order, created_at, updated_at
		 FROM projects WHERE id = ?`,
		id,
	).Scan(
		&p.ID, &p.Name, &p.EncryptedContent, &p.KeyCheck,
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
		`INSERT INTO projects (name, encrypted_content, key_check, sort_order)
		 VALUES (?, ?, ?, ?)`,
		p.Name, p.EncryptedContent, p.KeyCheck, p.SortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	id, _ := result.LastInsertId()
	return r.GetByID(id)
}

func (r *ProjectRepository) Update(p *model.Project) error {
	_, err := r.db.Exec(
		`UPDATE projects SET name = ?, encrypted_content = ?,
		        key_check = ?, sort_order = ?
		 WHERE id = ?`,
		p.Name, p.EncryptedContent,
		p.KeyCheck, p.SortOrder, p.ID,
	)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) Delete(id int64) error {
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

func (r *ProjectRepository) ListAll() ([]model.Project, error) {
	rows, err := r.db.Query(
		`SELECT id, name, encrypted_content, key_check,
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
			&p.ID, &p.Name, &p.EncryptedContent, &p.KeyCheck,
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
