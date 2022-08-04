package main

import (
	"net/http"
	"strings"

	"github.com/mohammadmokh/Movie-API/internal/data"
	"github.com/mohammadmokh/Movie-API/internal/validator"
)

func (app *application) showMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := getIDParams(r)
	if err != nil || id < 0 {
		http.NotFound(rw, r)
		return
	}
	movie, err := app.models.Movie.Get(id)
	if err != nil {
		if err == data.ErrNoRecord {
			http.NotFound(rw, r)
			return
		}
		app.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = writeJsonResponse(rw, http.StatusOK, map[string]data.Movie{"movie": *movie})
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) createMovie(rw http.ResponseWriter, r *http.Request) {

	var input struct {
		Genres  []string `json:"genres"`
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
	}

	err := readJson(rw, r, &input)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	movie := &data.Movie{
		Genres:  input.Genres,
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
	}

	validator := validator.New()
	data.ValidateMovie(validator, movie)
	if !validator.Valid() {

		err := writeJsonResponse(rw, http.StatusBadRequest, map[string]interface{}{"errors": validator.Errors})
		if err != nil {
			app.logger.Println(err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = app.models.Movie.Create(movie)
	if err != nil {

		if err == data.ErrEditConflict {
			err = writeJsonResponse(rw, http.StatusConflict, map[string]string{"error": "a conflict happend. try again"})
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				app.logger.Println(err)
			}
			return
		}

		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logger.Println(err)
		return
	}

	err = writeJsonResponse(rw, http.StatusOK, map[string]data.Movie{"movie": *movie})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logger.Println(err)
	}

}

func (app *application) deleteMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := getIDParams(r)
	if err != nil || id < 0 {
		http.NotFound(rw, r)
		return
	}

	err = app.models.Movie.Delete(id)
	if err != nil {

		if err == data.ErrNoRecord {
			http.NotFound(rw, r)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logger.Println(err)
		return
	}

	err = writeJsonResponse(rw, http.StatusOK, map[string]string{"message": "deleted succesfully"})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logger.Println(err)
	}
}

func (app *application) updateMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := getIDParams(r)
	if err != nil || id < 0 {
		http.NotFound(rw, r)
		return
	}
	movie, err := app.models.Movie.Get(id)
	if err != nil {
		if err == data.ErrNoRecord {
			http.NotFound(rw, r)
			return
		}
		app.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var input struct {
		Genres  []string `json:"genres"`
		Title   *string  `json:"title"`
		Year    *int32   `json:"year"`
		Runtime *int32   `json:"runtime"`
	}

	err = readJson(rw, r, &input)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	validator := validator.New()
	data.ValidateMovie(validator, movie)
	if !validator.Valid() {

		err := writeJsonResponse(rw, http.StatusBadRequest, map[string]interface{}{"errors": validator.Errors})
		if err != nil {
			app.logger.Println(err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = app.models.Movie.Update(id, movie)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logger.Println(err)
		return
	}

	err = writeJsonResponse(rw, http.StatusOK, map[string]data.Movie{"movie": *movie})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logger.Println(err)
	}
}

func (app *application) allMovies(rw http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string
		Genres  []string
		Filters data.Filter
	}

	input.Filters.SortSafeList = []string{"id", "-id", "title", "-title", "year", "-year", "runtime", "-runtime"}

	var err error
	validator := validator.New()

	input.Title = r.URL.Query().Get("title")
	input.Genres = strings.Split(r.URL.Query().Get("genres"), ",")
	input.Filters.Page = getIntqs(*r.URL, "page", validator, 1)
	input.Filters.PageSize = getIntqs(*r.URL, "pageSize", validator, 20)
	input.Filters.Sort = r.URL.Query().Get("sort")
	if input.Filters.Sort == "" {
		input.Filters.Sort = "id"
	}
	data.ValidateFilters(input.Filters, validator)
	if !validator.Valid() {
		err := writeJsonResponse(rw, http.StatusBadRequest, map[string]interface{}{"errors": validator.Errors})
		if err != nil {
			app.logger.Println(err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	movies, metadata, err := app.models.Movie.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = writeJsonResponse(rw, http.StatusOK, map[string]interface{}{"metadate": metadata, "movie": movies})
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
