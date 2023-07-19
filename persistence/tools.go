package persistence

import "strings"

func CheckErr(err error) error {

	if err != nil {
		if !strings.Contains(err.Error(), "LastInsertId") {
			return err
		}
	}
	return nil
}

type Pager struct {
	PageNo     int // start with 1
	PageSize   int // > 0
	TotalPage  int
	TotalCount int
}

func (p *Pager) IsFirstPage() bool {
	return p.PageNo == 1
}

func (p *Pager) IsLastPage() bool {
	return p.PageNo == p.TotalPage
}

func (p *Pager) GetLimit() (limit int, start int) {
	limit = p.PageSize
	start = (p.PageNo - 1) * p.PageSize
	return limit, start
}

func (p *Pager) SetPageNo(no int) {
	p.PageNo = no
	p.update()
}

func (p *Pager) Next() {
	p.SetPageNo(p.PageNo + 1)
}

func (p *Pager) SetTotal(total int) {
	p.TotalCount = total
	p.update()
}

func (p *Pager) update() {
	if p.TotalCount > 0 && p.PageSize > 0 {
		tp := p.TotalCount / p.PageSize
		if p.TotalCount%p.PageSize > 0 {
			tp += 1
		}
		p.TotalPage = tp
	}
}

func BuildPager(pageNo int, pageSize int) Pager {
	p := Pager{PageNo: pageNo, PageSize: pageSize}
	return p
}
