package data

import (
	"strings"

	"github.com/mohammadmokh/Movie-API/internal/validator"
)

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type MetaData struct {
	CorrentPage int
	LastPage    int
	Total       int
}

func ValidateFilters(f Filter, v *validator.Validator) {

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page < 10_000_000, "page", "maximum is 10 milion")
	v.Check(f.PageSize > 0, "pageSize", "must be greater than zero")
	v.Check(f.PageSize < 100, "pageSize", "maximum is 100")
	v.Check(validator.Contains(f.SortSafeList, f.Sort), "sort", "invalid")
}

func (f Filter) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
