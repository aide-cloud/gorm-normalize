package query

var _ Pagination = (*Page)(nil)

var (
	defaultCurr int32 = 1
	defaultSize int32 = 10
)

type Pagination interface {
	GetCurr() int32
	GetSize() int32
	SetTotal(total int64)
	GetTotal() int64
}

type Page struct {
	Curr int32 `json:"curr"`
	Size int32 `json:"size"`
	Total int64 `json:"total"`
}

// WithDefaultCurr is used to set default curr
func WithDefaultCurr(curr int32) {
	defaultCurr = curr
}

// WithDefaultSize is used to set default size
func WithDefaultSize(size int32) {
	defaultSize = size
}

func NewPage(curr, size int32) *Page {
	return &Page{
		Curr: curr,
		Size: size,
	}
}

func (p *Page) GetCurr() int32 {
	if p == nil || p.Curr <= 0 {
		return defaultCurr
	}
	return p.Curr
}

func (p *Page) GetSize() int32 {
	if p == nil || p.Size <= 0 {
		return defaultSize
	}
	return p.Size
}

func (p *Page) SetTotal(total int64) {
	if p == nil {
		return
	}
	p.Total = total
}

func (p *Page) GetTotal() int64 {
	if p == nil {
		return 0
	}
	return p.Total
}
