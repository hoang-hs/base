package base

type Page struct {
	Page  uint32 `form:"page" json:"page"`
	Limit uint32 `form:"limit" json:"limit"`
}

var LimitDefault uint32 = 10

func SetLimitDefault(limit uint32) {
	LimitDefault = limit
}

func (p *Page) SetDefault() {
	if p == nil {
		return
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit != 0 {
		p.Limit = LimitDefault
	}
}

func (p *Page) GetOffset() int {
	return int((p.Page - 1) * p.Limit)
}

func (p *Page) GetLimit() int {
	return int(p.Limit)
}
