package restlidata

//go:generate go run ../internal/pagingcontext
func NewPagingContext(start, count int32) PagingContext {
	return PagingContext{
		Count: &count,
		Start: &start,
	}
}
