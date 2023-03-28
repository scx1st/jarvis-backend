package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io"
	"jarvis-backend/config"
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

//GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(client *kubernetes.Clientset, podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取pod详情失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取pod详情失败,%v\n", err))
	}
	return pod, nil
}

//DeletePod 删除pod
func (p *pod) DeletePod(client *kubernetes.Clientset, podName, namespace string) (err error) {
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除pod失败,%v\n", err))
		return errors.New(fmt.Sprintf("删除pod失败,%v\n", err))
	}
	return nil
}

//UpdatePod 更新pod
func (p *pod) UpdatePod(client *kubernetes.Clientset, namespace, content string) (err error) {
	//1.content转成pod结构体，反序列化为pod对象,content就是pod的整个json体
	var pod = &corev1.Pod{}
	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化失败,%v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败,%v\n", err))
	}
	//2.更新pod
	_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新pod失败,%v\n", err))
		return errors.New(fmt.Sprintf("更新pod失败,%v\n", err))
	}
	return nil
}

//GetPodContainer 获取pod中的容器名
func (p *pod) GetPodContainer(client *kubernetes.Clientset, podName, namespace string) (containers []string, err error) {
	//1.获取pod详情
	pod, err := p.GetPodDetail(client, podName, namespace)
	if err != nil {
		return nil, err
	}
	//2.从pod对象中拿到容器名
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

//GetPodLog 获取pod中的容器日志
func (p *pod) GetPodLog(client *kubernetes.Clientset, containerName, podName, namespace string) (log string, err error) {
	//1.设置日志的配置，容器名以及tail的行数
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	//2.获取request实例
	req := client.CoreV1().Pods(namespace).GetLogs(podName, option)
	//3.发起request请求，返回一个ioReadCloser类型（等同于response。body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(fmt.Sprintf("获取podLog失败,%v\n", err))
		return "", errors.New(fmt.Sprintf("获取podLog失败,%v\n", err))
	}
	defer podLogs.Close()
	//4.将response body写入缓冲区， 目的是为了转成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(fmt.Sprintf("复制podLog失败,%v\n", err))
		return "", errors.New(fmt.Sprintf("复制podLog失败,%v\n", err))
	}
	return buf.String(), nil
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
