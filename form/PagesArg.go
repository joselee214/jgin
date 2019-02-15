package form

import (
	"errors"
)
type PagesArg struct {
	Pagesize int       `form:"pagesize" json:"pagesize"`
	Pagefrom int       `form:"pagefrom" json:"pagefrom"  validate:"gte=0"`
	Desc string        `form:"desc" json:"desc"`
	Asc  string        `form:"asc" json:"asc"`
}

func (p* PagesArg)Validate() (bool,error){
	if p.Pagesize>100 {
		return false,errors.New("一次只能请求100条数据")
	}
	if p.Pagefrom<0 {
		return false,errors.New("分页参数错误")
	}
	return true,nil
}

func (p* PagesArg)GetPageSize() (int){
	return 20
}

func (p* PagesArg)GetPageFrom() (int){
	return p.Pagefrom
}

func (p* PagesArg)GetDesc() (string){
	return p.Desc
}

func (p* PagesArg)GetAsc() (string){
	return p.Asc
}
