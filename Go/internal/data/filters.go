package data

import (
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	LastSeen       uuid.UUID `json:"lastSeen,omitzero"`
	Next           bool      `json:"next"`
	ResponseLength int       `json:"responseLength"`
}

type Filters struct {
	ID            *uuid.UUID `json:"id"`
	Name          *string    `json:"name"`
	Username      *string    `json:"username,omitempty"`
	Email         *string    `json:"email,omitempty"`
	CreatedAtFrom *time.Time `json:"createdAtFrom"`
	CreatedAtTo   *time.Time `json:"createdAtTo"`
	UpdatedAtFrom *time.Time `json:"updatedAtFrom"`
	UpdatedAtTo   *time.Time `json:"updatedAtTo"`

	OrderBy         []string   `json:"order_by,omitempty"`
	OrderBySafeList []string   `json:"order_by_safe_list,omitempty"`
	LastSeen        *uuid.UUID `json:"lastSeen"`
	PageSize        int        `json:"page_size,omitempty"`
}
