package services

import "github.com/morkid/paginate"

var PG *paginate.Pagination

// For the pagination
func InitPagination() {
	PG = paginate.New()
}
