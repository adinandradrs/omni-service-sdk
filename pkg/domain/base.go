package domain

import (
	"database/sql"
)

type BaseEntity struct {
	IsDeleted   bool           `json:"is_deleted"`
	CreatedBy   sql.NullString `json:"created_by"`
	CreatedDate sql.NullTime   `json:"created_date"`
	UpdatedBy   sql.NullString `json:"updated_by"`
	UpdatedDate sql.NullTime   `json:"updated_date"`
}

type BussinessError struct {
	ErrorCode    string
	ErrorMessage string
}

type TechnicalError struct {
	Exception string `json:"exception"`
	Occurred  int64  `json:"occurred_time"`
	Ticket    string `json:"ticket"`
}
