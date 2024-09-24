package repository

import (
	"context"
	"fmt"
	"gofi/database/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (repo *ProjectRepository) CreateProject(ctx context.Context, r *entity.Project) (*entity.Project, error) {
	var (
		lastInsertID uuid.UUID
		createdAt    time.Time
		updatedAt    time.Time
	)

	const query_insert = `
		INSERT INTO "project" (owner_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := repo.db.QueryRowContext(ctx, query_insert, r.Name).
		Scan(&lastInsertID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting project: %w", err)
	}

	r.ID = lastInsertID
	r.CreatedAt = createdAt
	r.UpdatedAt = updatedAt

	return r, nil
}

func (repo *ProjectRepository) GetProject(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	var r entity.Project

	const query_find_one = `
		SELECT * FROM "project" 
		WHERE id=$1
	`

	err := repo.db.GetContext(ctx, &r, query_find_one, id)
	if err != nil {
		return nil, fmt.Errorf("error getting project: %v", err)
	}

	return &r, nil
}

func (repo *ProjectRepository) ListProjects(ctx context.Context) ([]entity.Project, error) {
	var projects []entity.Project

	const query_find_all = `
		SELECT * FROM "project"
	`

	err := repo.db.SelectContext(ctx, &projects, query_find_all)
	if err != nil {
		return nil, fmt.Errorf("error listing projects: %v", err)
	}

	return projects, nil
}

func (repo *ProjectRepository) UpdateProject(ctx context.Context, r *entity.Project) (*entity.Project, error) {
	const query_update = `
		UPDATE "project" SET owner_id=:owner_id, name=:name, description=:description, updated_at=:updated_at 
		WHERE id=:id
	`

	_, err := repo.db.NamedExecContext(ctx, query_update, r)
	if err != nil {
		return nil, fmt.Errorf("error updating project: %v", err)
	}

	return r, nil
}

func (repo *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	const query_delete = `
		DELETE FROM "project" 
		WHERE id=$1
	`

	_, err := repo.db.ExecContext(ctx, query_delete, id)
	if err != nil {
		return fmt.Errorf("error deleting project: %v", err)
	}

	return nil
}
