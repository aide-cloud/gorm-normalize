package query

var _ Pagination = (*Page)(nil)

type Pagination interface {
	Page() int
	Size() int
	SetTotal(total int64)
}

type Page struct {
	page  int
	size  int
	total int64
}

func NewPage(page, size int) *Page {
	return &Page{
		page: page,
		size: size,
	}
}

func (p *Page) Page() int {
	return p.page
}

func (p *Page) Size() int {
	return p.size
}

func (p *Page) SetTotal(total int64) {
	p.total = total
}
