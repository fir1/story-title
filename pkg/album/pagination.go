package album

type Pagination struct {
	Limit  uint32
	Offset uint32
}

type PaginationResponse struct {
	TotalItems uint32
	PageSize   uint32
}
