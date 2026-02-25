package models

import "github.com/Doris-Mwito5/ginja-ai/internal/null"

type Pagination struct {
	Count    int    `json:"count"`
	NextPage *int64 `json:"next_page"`
	NumPages int    `json:"num_pages"`
	Page     int    `json:"page"`
	Per      int    `json:"per"`
	PrevPage *int64 `json:"prev_page"`
}

func NewPagination(count, page, per int) *Pagination {
	var prevPage, nextPage *int64

	if page > 1 {
		prevPage = null.NullValue(int64(page - 1))
	}

	if per < 1 {
		per = 10
	}

	numPages := count / per
	if count == 0 {
		numPages = 1
	} else if count%per != 0 {
		numPages++
	}

	if page < numPages {
		nextPage = null.NullValue(int64(page + 1))
	}

	return &Pagination{
		Count:    count,
		NextPage: nextPage,
		NumPages: numPages,
		Page:     page,
		Per:      per,
		PrevPage: prevPage,
	}
}
