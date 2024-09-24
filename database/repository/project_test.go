package repository

import (
	"context"
	"fmt"
	"gofi/database/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	p := &entity.Project{
		OwnerID:     uuid.New(),
		Name:        "Test Project",
		Description: "Test Description",
	}

	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *ProjectRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				// Mock the expected query and result
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery(`INSERT INTO "project" (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`).
					WithArgs(p.Name).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				record, err := repo.CreateProject(context.Background(), p)
				require.NoError(t, err)
				require.NotNil(t, record)
				require.Equal(t, expectedID, record.ID)
				require.Equal(t, expectedCreatedAt, record.CreatedAt)
				require.Equal(t, expectedUpdatedAt, record.UpdatedAt)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting project",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO "project" (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`).
					WithArgs(p.Name).
					WillReturnError(fmt.Errorf("error inserting project"))

				_, err := repo.CreateProject(context.Background(), p)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProjectRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestGetProject(t *testing.T) {
	p := &entity.Project{
		OwnerID:     uuid.New(),
		Name:        "Test Project",
		Description: "Test Description",
	}

	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *ProjectRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "owner_id", "name", "description"}).
					AddRow(expectedID, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.OwnerID, p.Name, p.Description)

				mock.ExpectQuery(`SELECT * FROM "project" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnRows(rows)

				record, err := repo.GetProject(context.Background(), expectedID)
				require.NoError(t, err)
				require.Equal(t, expectedID, record.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting project",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT * FROM "project" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnError(fmt.Errorf("error getting project"))

				_, err := repo.GetProject(context.Background(), expectedID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProjectRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestListProjects(t *testing.T) {
	p := &entity.Project{
		OwnerID:     uuid.New(),
		Name:        "Test Project",
		Description: "Test Description",
	}

	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *ProjectRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "owner_id", "name", "description"}).
					AddRow(expectedID, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.OwnerID, p.Name, p.Description)

				mock.ExpectQuery(`SELECT * FROM "project"`).WillReturnRows(rows)

				records, err := repo.ListProjects(context.Background())
				require.NoError(t, err)
				require.Len(t, records, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying project",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT * FROM "project"`).WillReturnError(fmt.Errorf("error querying project"))

				_, err := repo.ListProjects(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProjectRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestUpdateProject(t *testing.T) {
	expectedID := uuid.New()

	p := &entity.Project{
		OwnerID:     expectedID,
		Name:        "Test Project",
		Description: "Test Description",
	}

	np := &entity.Project{
		ID:          expectedID,
		Name:        "Another Project",
		Description: "Another Description",
	}

	tcs := []struct {
		name string
		test func(*testing.T, *ProjectRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery(`INSERT INTO "project" (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`).
					WithArgs(p.Name).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				cp, err := repo.CreateProject(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, expectedID, cp.ID)
				require.Equal(t, expectedCreatedAt, cp.CreatedAt)
				require.Equal(t, expectedUpdatedAt, cp.UpdatedAt)

				mock.ExpectExec(`UPDATE "project" SET owner_id=?, name=?, description=?, updated_at=? WHERE id=?`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := repo.UpdateProject(context.Background(), np)
				require.NoError(t, err)
				require.Equal(t, expectedID, up.ID)
				require.Equal(t, np.Name, up.Name)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating project",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "project" SET owner_id=?, name=?, description=?, updated_at=? WHERE id=?`).
					WillReturnError(fmt.Errorf("error updating project"))

				_, err := repo.UpdateProject(context.Background(), p)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProjectRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestDeleteProject(t *testing.T) {
	expectedID := uuid.New()

	tcs := []struct {
		name string
		test func(*testing.T, *ProjectRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "project" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.DeleteProject(context.Background(), expectedID)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting project",
			test: func(t *testing.T, repo *ProjectRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "project" WHERE id=$1`).
					WithArgs(expectedID).
					WillReturnError(fmt.Errorf("error deleting project"))

				err := repo.DeleteProject(context.Background(), expectedID)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProjectRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}
