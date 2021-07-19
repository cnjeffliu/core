package k8s

import (
	"context"
	"serv/core/logx"
	"testing"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestK8s(t *testing.T) {
	logx.InitLog("./debug.log")
	InitK8s(WithKubeConfigPatterm("../conf.d/kube/cluster*"))
	clientset := ClientSet[1].Client

	list, err := clientset.CoreV1().Pods(coreV1.NamespaceDefault).List(context.TODO(), metaV1.ListOptions{})
	if err != err {
		logx.Fatal(err)
	}

	for _, pod := range list.Items {
		logx.Debugf("NameSpace:%v \t Name:%v\n", pod.Name, pod.Namespace)
	}

	//列出pod
	// podList, err := clientset.CoreV1().Pods(coreV1.NamespaceDefault).List(&meta_v1.ListOptions{})

	//查询pod
	// pod, err := ClientSet[1].Client.CoreV1().Pods(coreV1.NamespaceDefault).Get(<podName>, meta_v1.GetOptions{})
	// pod, err = ClientSet[1].Client.PodList.Pods(coreV1.NamespaceDefault).Get(podName)

	//创建pod
	// pod, err := clientset.CoreV1().Pods(coreV1.NamespaceDefault).Create(web)

	//更新pod
	// pod, err := clientset.CoreV1().Pods(coreV1.NamespaceDefault).Update(web)

	//删除pod
	// err := clientset.CoreV1().Pods(coreV1.NamespaceDefault).Delete(<podName>, &meta_v1.DeleteOptions{})
}
