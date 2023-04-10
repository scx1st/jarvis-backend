package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
)

var HelmConfig helmConfig

type helmConfig struct {
	//这种方式初始化行不通，比如有10个命名空间，10个actionConfig都是最后一个初始化的明明空间
	//ActionConfigMap map[string]*action.Configuration
}

//GetAc 获取helm config配置
func (*helmConfig) GetAc(cluster, namespace string) (*action.Configuration, error) {
	//获取kubeconfig
	kubeconfig, ok := K8s.KubeConfMap[cluster]
	if !ok {
		logger.Error("actionConfig初始化失败,cluster不存在")
		return nil, errors.New("actionConfig初始化失败,cluster不存在")
	}
	//new一个actionConfig对象
	actionConfig := new(action.Configuration)
	cf := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeconfig,
		Namespace:  &namespace,
	}
	if err := actionConfig.Init(cf, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		logger.Error(fmt.Sprintf("actionConfig初始化失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("actionConfig初始化失败, %v\n", err))
	}
	return actionConfig, nil
}
