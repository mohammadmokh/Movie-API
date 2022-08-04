package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mohammadmokh/Movie-API/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func readJson(rw http.ResponseWriter, r *http.Request, input interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed json at character %d", syntaxError.Offset)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains invalid json type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains invalid json type at character %d", unmarshalTypeError.Offset)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed json")

		default:
			return err
		}
	}

	return nil
}

func writeJsonResponse(rw http.ResponseWriter, status int, data interface{}) error {

	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(json)
	return nil
}

func getIDParams(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	return id, err
}

func getIntqs(qs url.URL, key string, v *validator.Validator, defaultval int) int {

	val, err := strconv.Atoi(qs.Query().Get(key))
	if qs.Query().Get(key) == "" {
		return defaultval
	}
	if err != nil {
		v.Add(key, "must be an integer number")
	}
	return val
}

func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
