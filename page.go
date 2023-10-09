package query

var _ Pagination = (*Page)(nil)

type Pagination interface {
	GetCurr() int
	GetSize() int
	SetTotal(total int64)
}

type Page struct {
	Curr  int   `json:"curr"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

func NewPage(curr, size int) *Page {
	return &Page{
		Curr: curr,
		Size: size,
	}
}

func (p *Page) GetCurr() int {
	return p.Curr
}

func (p *Page) GetSize() int {
	return p.Size
}

func (p *Page) SetTotal(total int64) {
	p.Total = total
}
