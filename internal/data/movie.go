package data

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mohammadmokh/Movie-API/internal/validator"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"-"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {

	v.Check(movie.Title != "", "title", "title cannot be empty")
	v.Check(len(movie.Title) <= 500, "title", "title cannot be more than 500 character")

	v.Check(movie.Year != 0, "year", "year cannot be empty")
	v.Check(movie.Year >= 1888, "year", "invalid year")

	v.Check(movie.Runtime != 0, "runtime", "runtime cannot be empty")
	v.Check(movie.Runtime >= 1, "runtime", "invalid runtime")

	v.Check(len(movie.Genres) <= 5 && len(movie.Genres) >= 1, "genres", "movie should have between 1 and 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "genres should be unique")
}

type MovieModel struct {
	DB *sql.DB
}

func (m *MovieModel) Create(movie *Movie) error {

	query := `INSERT INTO movies (title, year, runtime, genres) VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`
	err := m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)).
		Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	return err
}

func (m *MovieModel) Get(id int64) (*Movie, error) {

	movie := &Movie{}
	query := `SELECT id, title, year, runtime, genres, version FROM movies WHERE id = $1`
	err := m.DB.QueryRow(query, id).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return movie, nil
}

func (m *MovieModel) Delete(id int64) error {

	query := `DELETE FROM movies WHERE id = $1`
	res, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *MovieModel) Update(id int64, movie *Movie) error {

	query := `UPDATE movies	SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1 WHERE id = $5 AND version = $6
	 RETURNING version`
	err := m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), id, movie.Version).Scan(&movie.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (m *MovieModel) GetAll(title string, genres []string, f Filter) ([]*Movie, MetaData, error) {

	var rows *sql.Rows
	var err error

	fmt.Println(f.Page, f.PageSize)

	if title == "" && (len(genres) == 1 && genres[0] == "") {
		query := fmt.Sprintf(
			`SELECT count(*) OVER(), id, title, year, runtime, genres, version FROM movies ORDER BY %s %s, id ASC
			 OFFSET %d LIMIT %d`, strings.TrimPrefix(f.Sort, "-"), f.sortDirection(), (f.Page-1)*f.PageSize, f.PageSize)
		rows, err = m.DB.Query(query)

	} else if title == "" && (genres[0] != "") {
		query := fmt.Sprintf(
			`SELECT count(*) OVER(), id, title, year, runtime, genres, version FROM movies WHERE (genres @> $1) ORDER BY %s %s, id ASC 
			OFFSET %d LIMIT %d`, strings.TrimPrefix(f.Sort, "-"), f.sortDirection(), (f.PageSize-1)*f.Page, f.PageSize)
		rows, err = m.DB.Query(query, pq.Array(genres))

	} else if title != "" && (len(genres) == 1 && genres[0] == "") {
		query := fmt.Sprintf(`SELECT count(*) OVER(), id, title, year, runtime, genres, version FROM movies WHERE LOWER(title) LIKE 
		'%%' || $1 || '%%' ORDER BY %s %s, id ASC OFFSET %d LIMIT %d`, strings.TrimPrefix(f.Sort, "-"),
			f.sortDirection(), (f.PageSize-1)*f.Page, f.PageSize)
		rows, err = m.DB.Query(query, title)

	} else {
		query := fmt.Sprintf(`SELECT count(*) OVER(), id, title, year, runtime, genres, version FROM movies WHERE (genres @> $1) AND 
		(LOWER(title) LIKE '%%' || $2 || '%%') ORDER BY %s %s, id ASC  OFFSET %d LIMIT %d`,
			strings.TrimPrefix(f.Sort, "-"), f.sortDirection(), (f.PageSize-1)*f.Page, f.PageSize)
		rows, err = m.DB.Query(query, pq.Array(genres), title)
	}

	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var movies []*Movie
	var total int
	for rows.Next() {

		movie := Movie{}

		err = rows.Scan(
			&total, &movie.ID, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
		if err != nil {
			return nil, MetaData{}, err
		}

		movies = append(movies, &movie)
	}

	metadata := MetaData{
		CorrentPage: f.Page,
		LastPage:    int(math.Ceil(float64(total) / float64(f.PageSize))),
		Total:       total,
	}
	return movies, metadata, nil
}
