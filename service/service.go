package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var Servicev1 servicev1

type servicev1 struct{}

//ServiceCreate 定义serviceCreate结构体
type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
	Cluster       string            `json:"cluster"`
}

//CreateService 创建service
func (s *servicev1) CreateService(client *kubernetes.Clientset, data *ServiceCreate) (err error) {
	//将data中的数据组装成corev1.Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
		Status: corev1.ServiceStatus{},
	}
	//根据service类型来判断不同的配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}
	//创建service
	_, err = client.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("创建Service失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建Service失败, %v\n", err))
	}
	return nil
}
