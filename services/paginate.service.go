package services

import "github.com/morkid/paginate"

var PG *paginate.Pagination

func InitPagination() {
	PG = paginate.New()
}
