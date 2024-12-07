package models

type Request struct {
	Limit    int      `query:"limit"`
	Offset   int      `query:"offset"`
	Preloads []string `query:"preload"`
}

func (r *Request) Validate() {
	if r.Limit <= 0 {
		r.Limit = 10
	}

	if r.Offset <= 0 {
		r.Offset = 0
	}
}

type Meta struct {
	Total  int
	Limit  int
	Offset int
}

type Response struct {
	Message string
	Data    interface{}
	Meta    *Meta
}
