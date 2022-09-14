/*
 * @Author: Jeffrey.Liu
 * @Date: 2021-12-06 14:45:26
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-17 14:58:54
 * @Description: k8s操作处理类
 */
package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	k8srun "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type STATUS int8

const (
	STATUS_UNKNOWN    STATUS = iota
	STATUS_NOT_FOUND         // not found pod
	STATUS_DELETED           // graceful deleted
	STATUS_RUNNING           // pod running, not be in deleting
	STATUS_PENDING           // pending,  kubelet scheduled, deletetimestamp is nil
	STATUS_UNASSIGNED        // not kubelet scheduled
	STATUS_ERROR
)

const (
	DELETE_IMMEDIATLY int64 = iota // delete immediatly
	DELETE_NORMAL                  // graceful delete , 1 second
)

type K8SCli struct {
	clientset *kubernetes.Clientset
	podLister listerv1.PodLister
}

type K8SCliOption func() string

func WithMasterUrl(url string) K8SCliOption {
	return func() string {
		return url
	}
}

func NewK8SCli(kubeconfigPath string, opts ...K8SCliOption) *K8SCli {
	var set *kubernetes.Clientset
	var client = &K8SCli{}
	var masterUrl = ""
	for _, opt := range opts {
		masterUrl = opt()
	}

	config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	set, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client.clientset = set
	stopper := make(chan struct{})
	factory := informers.NewSharedInformerFactory(set, 0)
	nodeInformer := factory.Core().V1().Nodes().Informer()
	podInformer := factory.Core().V1().Pods().Informer()

	go factory.Start(stopper)

	go func() {
		if !cache.WaitForCacheSync(stopper, podInformer.HasSynced, nodeInformer.HasSynced) {
			fmt.Println("wait for cache sync error")
			k8srun.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
			return
		}
	}()

	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			node := obj.(*corev1.Node)
			fmt.Println("delete not implemented", node.Name)
		},
	})

	// podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc: func(obj interface{}) {
	// 		pod := obj.(*corev1.Pod)
	// 		fmt.Printf("New Pod Added : %#v\n", pod.Name)

	// 	},
	// 	UpdateFunc: func(oldObj, newObj interface{}) {
	// 		oldPod := oldObj.(*corev1.Pod)
	// 		newPod := newObj.(*corev1.Pod)
	// 		fmt.Printf("Pod Updated : %v => %v\n", oldPod.Name, newPod.Name)
	// 	},
	// 	DeleteFunc: func(obj interface{}) {
	// 		pod := obj.(*corev1.Pod)
	// 		fmt.Printf("Pod Deleted : %v\n", pod.Name)
	// 	},
	// })
	client.podLister = factory.Core().V1().Pods().Lister()
	return client
}

func WithNamespace(namespace string) K8SCliOption {
	return func() string {
		return namespace
	}
}

func (c *K8SCli) ListPods(opts ...K8SCliOption) []corev1.Pod {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}
	podList, err := c.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != err {
		fmt.Println(err)
		return nil
	}

	return podList.Items
}

func checkError(pod *v1.Pod, err error) (STATUS, *v1.Pod, error) {
	if errors.IsNotFound(err) {
		return STATUS_NOT_FOUND, nil, err
	}

	if err == nil {
		if pod.Status.StartTime == nil {
			return STATUS_UNASSIGNED, pod, nil
		} else if pod.DeletionTimestamp != nil {
			return STATUS_DELETED, pod, nil
		} else if pod.Status.Phase == v1.PodRunning {
			return STATUS_RUNNING, pod, nil
		} else if pod.Status.Phase == v1.PodPending && pod.Status.StartTime != nil {
			return STATUS_PENDING, pod, nil
		}

		return STATUS_UNKNOWN, pod, nil
	}

	return STATUS_ERROR, nil, err
}

func (c *K8SCli) GetPod(name string, opts ...K8SCliOption) (STATUS, *v1.Pod, error) {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}

	return checkError(c.podLister.Pods(namespace).Get(name))
}

func (c *K8SCli) GetPodRT(name string, opts ...K8SCliOption) (STATUS, *v1.Pod, error) {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}

	return checkError(c.clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{}))
}

func GetPodTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}
}

func GetObjTypeMeta(podname string) metav1.ObjectMeta {
	return metav1.ObjectMeta{Name: podname,
		Namespace: corev1.NamespaceDefault,
		Labels:    map[string]string{"name": podname}}
}

func (c *K8SCli) WatchPod(label string, opts ...K8SCliOption) (watch.Interface, error) {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}

	return c.clientset.CoreV1().Pods(namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: label,
	})
}

func (c *K8SCli) CreatePod(pod *corev1.Pod, opts ...K8SCliOption) (*corev1.Pod, error) {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}

	return c.clientset.CoreV1().Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
}

func (c *K8SCli) UpdatePod(pod *corev1.Pod, opts ...K8SCliOption) (*corev1.Pod, error) {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}

	return c.clientset.CoreV1().Pods(namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
}

func (c *K8SCli) DeletePodImmediately(name string, opts ...K8SCliOption) error {
	return c.DeletePod(name, DELETE_IMMEDIATLY, opts...)
}

func (c *K8SCli) DeletePodNormal(name string, opts ...K8SCliOption) error {
	return c.DeletePod(name, DELETE_NORMAL, opts...)
}

func (c *K8SCli) DeletePod(name string, gracePeriodSec int64, opts ...K8SCliOption) error {
	var namespace string = corev1.NamespaceDefault
	for _, opt := range opts {
		namespace = opt()
	}

	return c.clientset.CoreV1().Pods(namespace).Delete(
		context.Background(),
		name,
		metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSec},
	)
}
