package service

import "time"

//dataSelector 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDatalist   []DataCell
	dataSelectorQuery *DataSelectorQuery
}

//DataCell 用于各种资源list的类型转换，转换后可以使用dataSelector的排序、过滤和分页
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

//DataSelectorQuery 定义过滤和分页属性，过滤：Name，分页：limit和page
type DataSelectorQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

//Len 用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDatalist)
}

//Swap 用于数组中的元素在比较大小后怎么交换位置， 可定义升降序
//i, j是切片的下标
func (d *dataSelector) Swap(i, j int) {
	d.GenericDatalist[i], d.GenericDatalist[j] = d.GenericDatalist[j], d.GenericDatalist[i]
}

//Less 用于定义数组中的元素排序的“大小”的比较方式
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDatalist[i].GetCreation()
	b := d.GenericDatalist[j].GetCreation()
	return b.Before(a)
}
