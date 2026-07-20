package pagination

import "strconv"

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Params struct {
	Page     int
	PageSize int
}

func (p Params) Offset() int { return (p.Page - 1) * p.PageSize }
func (p Params) Limit() int  { return p.PageSize }

func Parse(pageStr, sizeStr string, defaultSize, maxSize int) Params {
	if defaultSize <= 0 {
		defaultSize = DefaultPageSize
	}
	if maxSize <= 0 {
		maxSize = MaxPageSize
	}

	page := 1
	if v, err := strconv.Atoi(pageStr); err == nil && v > 1 {
		page = v
	}

	size := defaultSize
	if v, err := strconv.Atoi(sizeStr); err == nil && v > 0 {
		size = v
	}
	if size > maxSize {
		size = maxSize
	}

	return Params{Page: page, PageSize: size}
}

type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func Slice[T any](items []T, p Params) []T {
	start := p.Offset()
	if start < 0 || start >= len(items) {
		return []T{}
	}
	end := start + p.Limit()
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func NewPage[T any](items []T, total int, p Params) Page[T] {
	if items == nil {
		items = []T{}
	}

	totalPages := 0
	if p.PageSize > 0 {
		totalPages = (total + p.PageSize - 1) / p.PageSize
	}

	return Page[T]{
		Items:      items,
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
