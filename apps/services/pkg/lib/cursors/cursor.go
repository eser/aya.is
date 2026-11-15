package cursors

const (
	defaultLimit = 20
	// maxLimit bounds one page: the web client asks for at most 24 rows, and an unbounded value both wraps when
	// it is narrowed to int32 for SQL and lets a single request read the whole table.
	maxLimit = 100

	defaultSortBy  = "created_at"
	defaultSortDir = "asc"
)

type Cursor struct {
	Filters map[string]string

	Offset *string
	Seed   string

	SortBy  string
	SortDir string

	Limit int
}

func NewCursor(limit int, offset *string) *Cursor {
	if limit <= 0 {
		limit = defaultLimit
	}

	limit = min(limit, maxLimit)

	return &Cursor{
		Limit:  limit,
		Offset: offset,
		Seed:   "",

		SortBy:  defaultSortBy,
		SortDir: defaultSortDir,

		Filters: make(map[string]string),
	}
}
