package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Pod pod

type pod struct{}

//PodsResp 定义列表的返回类型
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// GetPods 获取pod列表
func (p *pod) GetPods(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取pod列表数据失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取pod列表数据失败,%v\n", err))
	}
	//实例化dataSelector对象
	selectableData := &dataSelector{
		//PodList.Items 返回的是一个切片（slice）
		GenericDatalist: p.toCell(podList.Items),
		dataSelectorQuery: &DataSelectorQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDatalist)
	//再排序和分页
	data := filtered.Sort().Paginate()
	//将[]yDataCell类型的pod列表转为v1.pod列表
	pods := p.fromCell(data.GenericDatalist)

	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}

//toCell 定义pod到DataCell类型转换的方法
func (p *pod) toCell(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

//fromCell 定义DataCell到pod类型转换的方法
func (p *pod) fromCell(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
