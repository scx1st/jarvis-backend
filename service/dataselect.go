package service

import (
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

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

//Sort 重写以上三个方法用使用sort.Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

//Filter 用于过滤元素，比较元素的Name的属性，若包含，则返回
func (d *dataSelector) Filter() *dataSelector {
	//若Name的传参为空，则返回所有元素
	if d.dataSelectorQuery.FilterQuery.Name == "" {
		return d
	}
	//若Name的传参不为空，则返回元素中包含name的所有元素
	filteredList := []DataCell{}
	for _, value := range d.GenericDatalist {
		matched := true
		objName := value.GetName()
		if !strings.Contains(objName, d.dataSelectorQuery.FilterQuery.Name) {
			matched = false
			continue
		}
		if matched {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDatalist = filteredList
	return d
}

//Paginate 用于数组分页，根据Limit和Page的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectorQuery.PaginateQuery.Limit
	page := d.dataSelectorQuery.PaginateQuery.Page
	//验证参数合法，若不合法，则返回全部数据
	if limit <= 0 || page <= 0 {
		return d
	}
	//定义offset
	//举例：25个元素的切片 limit10
	// page1 start 0 end 10
	// page1 start 10 end 20
	// page1 start 20 end 30
	startIndex := limit * (page - 1)
	endIndex := limit * page
	if len(d.GenericDatalist) < endIndex {
		endIndex = len(d.GenericDatalist)
	}
	d.GenericDatalist = d.GenericDatalist[startIndex:endIndex]
	return d
}

//podCell 定义podCell类型，实现两个方法GetCreation GetName，可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
