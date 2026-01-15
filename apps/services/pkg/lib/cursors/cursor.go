package cursors

const (
	defaultLimit = 20

	defaultSortBy  = "created_at"
	defaultSortDir = "asc"
)

type Cursor struct {
	Filters map[string]string

	Offset *string

	SortBy  string
	SortDir string

	Limit int
}

func NewCursor(limit int, offset *string) *Cursor {
	if limit <= 0 {
		limit = defaultLimit
	}

	return &Cursor{
		Limit:  limit,
		Offset: offset,

		SortBy:  defaultSortBy,
		SortDir: defaultSortDir,

		Filters: make(map[string]string),
	}
}
