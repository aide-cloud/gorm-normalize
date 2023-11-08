package query

var _ Pagination = (*Page)(nil)

var (
	defaultCurr = 1
	defaultSize = 10
)

type Pagination interface {
	GetCurr() int
	GetSize() int
	SetTotal(total int64)
	GetTotal() int64
}

type Page struct {
	Curr  int   `json:"curr"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

// WithDefaultCurr is used to set default curr
func WithDefaultCurr(curr int) {
	defaultCurr = curr
}

// WithDefaultSize is used to set default size
func WithDefaultSize(size int) {
	defaultSize = size
}

func NewPage(curr, size int) *Page {
	return &Page{
		Curr: curr,
		Size: size,
	}
}

func (p *Page) GetCurr() int {
	if p == nil || p.Curr <= 0 {
		return defaultCurr
	}
	return p.Curr
}

func (p *Page) GetSize() int {
	if p == nil || p.Size <= 0 {
		return defaultSize
	}
	return p.Size
}

func (p *Page) SetTotal(total int64) {
	if p == nil {
		_p := NewPage(defaultCurr, defaultSize)
		*p = *_p
	}
	p.Total = total
}

func (p *Page) GetTotal() int64 {
	if p == nil {
		return 0
	}
	return p.Total
}
