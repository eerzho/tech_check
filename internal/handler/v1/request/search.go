package request

type (
	Search struct {
		Pagination Pagination
		Filters    map[string]string
		Sorts      map[string]string
	}

	Pagination struct {
		Page  int
		Count int
	}
)
