package entity

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
	OwnerID     uuid.UUID  `json:"owner_id" db:"owner_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
}

type ProjectReq struct {
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type ProjectRes struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	OwnerID     uuid.UUID  `json:"owner_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}
