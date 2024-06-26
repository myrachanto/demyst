package support

type Paginator struct {
	Orderby  string `json:"orderby"`
	Order    int    `json:"order"`
	Search   string `json:"search"`
	Style    string `json:"style"`
	Page     int    `json:"page"`
	Pagesize int    `json:"pagesize"`
}
type Paginator2 struct {
	Location    string `json:"location"`
	Sublocation string `json:"sublocation"`
	Kind        string `json:"kind"`
	Majorcat    string `json:"majorcat"`
	Page        int    `json:"page"`
	Pagesize    int    `json:"pagesize"`
}

func Paginate(x []int, skip int, size int) []int {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}
