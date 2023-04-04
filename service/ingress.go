package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Ingress ingress

type ingress struct{}

//IngressCreate 定义ingressCreate结构体
type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
	Cluster   string                 `json:"cluster"`
}

//HttpPath 定义ingress的path结构体
type HttpPath struct {
	Path        string        `json:"path"`
	PathType    nwv1.PathType `json:"path_type"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
}

func (i *ingress) CreateIngress(client *kubernetes.Clientset, data *IngressCreate) (err error) {
	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面用于数据组装
	//ingressRule代表的是Hosts
	var ingressRules = make([]nwv1.IngressRule, 0)
	//httpIngressPaths代表的是Paths
	var httpIngressPaths = make([]nwv1.HTTPIngressPath, 0)
	//将data中的数据组装成nwv1.Ingress对象
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: nwv1.IngressStatus{},
	}
	//第一层for循环是将host组装成nwv1.IngressRule类型的对象
	//一个host对应一个ingressrule，每隔ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		//先把host放进去
		ir := nwv1.IngressRule{
			Host: key,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{
					Paths: nil,
				},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPaths = append(httpIngressPaths, hip)
		}
		//给Paths赋值，前面置空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		//将每个ir组装成数组
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules放到ingress中
	ingress.Spec.Rules = ingressRules
	//创建ingress
	_, err = client.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("创建Ingress失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建Ingress失败, %v\n", err))
	}
	return nil
}
