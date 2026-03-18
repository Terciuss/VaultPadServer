package repository

import (
	"database/sql"
	"fmt"

	"github.com/user/access-storage-server/internal/model"
)

type ShareRepository struct {
	db *sql.DB
}

func NewShareRepository(db *sql.DB) *ShareRepository {
	return &ShareRepository{db: db}
}

func (r *ShareRepository) Share(projectID, userID, sharedBy int64) error {
	_, err := r.db.Exec(
		"INSERT INTO project_shares (project_id, user_id, shared_by) VALUES (?, ?, ?)",
		projectID, userID, sharedBy,
	)
	if err != nil {
		return fmt.Errorf("share project: %w", err)
	}
	return nil
}

func (r *ShareRepository) Unshare(projectID, userID int64) error {
	result, err := r.db.Exec(
		"DELETE FROM project_shares WHERE project_id = ? AND user_id = ?",
		projectID, userID,
	)
	if err != nil {
		return fmt.Errorf("unshare project: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("share not found")
	}
	return nil
}

func (r *ShareRepository) ListSharedProjects(userID int64) ([]model.Project, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.name, p.encrypted_content, p.key_check,
		        p.sort_order, p.created_at, p.updated_at
		 FROM projects p
		 INNER JOIN project_shares ps ON ps.project_id = p.id
		 WHERE ps.user_id = ?
		 ORDER BY p.sort_order ASC, p.created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list shared projects: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.ID, &p.Name, &p.EncryptedContent, &p.KeyCheck,
			&p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan shared project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ShareRepository) ListSharedProjectsMeta(userID int64) ([]model.ProjectMeta, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.updated_at
		 FROM projects p
		 INNER JOIN project_shares ps ON ps.project_id = p.id
		 WHERE ps.user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list shared projects meta: %w", err)
	}
	defer rows.Close()

	var metas []model.ProjectMeta
	for rows.Next() {
		var m model.ProjectMeta
		if err := rows.Scan(&m.ID, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shared project meta: %w", err)
		}
		metas = append(metas, m)
	}
	return metas, rows.Err()
}

func (r *ShareRepository) HasAccess(projectID, userID int64) (bool, error) {
	var exists int
	err := r.db.QueryRow(
		"SELECT 1 FROM project_shares WHERE project_id = ? AND user_id = ? LIMIT 1",
		projectID, userID,
	).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check share access: %w", err)
	}
	return true, nil
}

func (r *ShareRepository) ListSharesByUser(userID int64) ([]model.ProjectShare, error) {
	rows, err := r.db.Query(
		"SELECT id, project_id, user_id, shared_by, created_at FROM project_shares WHERE user_id = ? ORDER BY created_at ASC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list shares by user: %w", err)
	}
	defer rows.Close()

	var shares []model.ProjectShare
	for rows.Next() {
		var s model.ProjectShare
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.UserID, &s.SharedBy, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan share: %w", err)
		}
		shares = append(shares, s)
	}
	return shares, rows.Err()
}

func (r *ShareRepository) ListSharesByProject(projectID int64) ([]model.ProjectShare, error) {
	rows, err := r.db.Query(
		"SELECT id, project_id, user_id, shared_by, created_at FROM project_shares WHERE project_id = ?",
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list shares: %w", err)
	}
	defer rows.Close()

	var shares []model.ProjectShare
	for rows.Next() {
		var s model.ProjectShare
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.UserID, &s.SharedBy, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan share: %w", err)
		}
		shares = append(shares, s)
	}
	return shares, rows.Err()
}
