package domain

import (
	"errors"
	"strconv"
)

const (
	ErrInvalidPaginationPage  = "el número de página debe ser mayor a 0"
	ErrInvalidPaginationLimit = "el límite debe estar entre 10 y 100"
	DefaultPage               = 1
	DefaultLimit              = 10
	MaxLimit                  = 100
	MinLimit                  = 10
)

type Pagination struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=10,max=100"`
	Search string `json:"search,omitempty"`
	Offset int    `json:"-"`
}

type PaginatedResult[T any] struct {
	Data       []T             `json:"data"`
	Pagination *PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"total_pages"`
	HasPrevious  bool  `json:"has_previous"`
	HasNext      bool  `json:"has_next"`
	PreviousPage *int  `json:"previous_page,omitempty"`
	NextPage     *int  `json:"next_page,omitempty"`
}

func NewPagination(page, limit int, search string) *Pagination {
	p := &Pagination{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	p.ApplyDefaults()
	p.CalculateOffset()

	return p
}

func (p *Pagination) ApplyDefaults() {
	if p.Page < DefaultPage {
		p.Page = DefaultPage
	}

	if p.Limit < DefaultLimit {
		p.Limit = DefaultLimit
	}

	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
}

func (p *Pagination) CalculateOffset() {
	p.Offset = (p.Page - 1) * p.Limit
}

func (p *Pagination) Validate() error {
	if p.Page < 1 {
		return errors.New(ErrInvalidPaginationPage)
	}

	if p.Limit < MinLimit || p.Limit > MaxLimit {
		return errors.New(ErrInvalidPaginationLimit)
	}

	return nil
}

func NewPaginatedResult[T any](data []T, pagination *Pagination, total int64) *PaginatedResult[T] {
	totalPages := int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	info := &PaginationInfo{
		Page:        pagination.Page,
		Limit:       pagination.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasPrevious: pagination.Page > 1,
		HasNext:     pagination.Page < totalPages,
	}

	if info.HasPrevious {
		prev := pagination.Page - 1
		info.PreviousPage = &prev
	}

	if info.HasNext {
		next := pagination.Page + 1
		info.NextPage = &next
	}

	return &PaginatedResult[T]{
		Data:       data,
		Pagination: info,
	}
}

func ParsePaginationFromQuery(pageStr, limitStr, searchStr string) (*Pagination, error) {
	var page, limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, errors.New(ErrInvalidPaginationPage)
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, errors.New(ErrInvalidPaginationLimit)
		}
	}

	pagination := NewPagination(page, limit, searchStr)

	if err := pagination.Validate(); err != nil {
		return nil, err
	}

	return pagination, nil
}
