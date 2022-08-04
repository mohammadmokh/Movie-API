package data

import (
	"database/sql"
	"errors"
)

var ErrNoRecord = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")

type Models struct {
	Movie MovieModel
}

func NewModel(DB *sql.DB) Models {

	return Models{
		Movie: MovieModel{DB: DB},
	}
}
