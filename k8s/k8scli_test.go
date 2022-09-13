/*
 * @Author: Jeffrey.Liu <zhifeng172@163.com>
 * @Date: 2021-07-19 11:58:51
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-17 14:55:02
 * @Description:
 */
package k8s

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
)

const k8sConfNameLen = 4

var ClientSet map[uint16]*K8SCli

type KubeOption func() string

func WithKubeConfigPatterm(pattern string) KubeOption {
	return func() string {
		return pattern
	}
}

func InitK8s(opts ...KubeOption) {
	var kubePattern string = "./conf.d/kube/cluster*"
	for _, opt := range opts {
		kubePattern = opt()
	}

	ClientSet = make(map[uint16]*K8SCli)
	matches, _ := filepath.Glob(kubePattern)
	for _, v := range matches {
		fmt.Println("field:", v)
		fields := strings.Split(v, "_")
		if len(fields) != k8sConfNameLen {
			//目录的命名规则是cluster_1_cslg_1，因此用_分隔后，有3个string
			//第1个字段，cluster为前缀标识匹配
			// 第2字段，标识clusterid，在tbl_storage表中记录
			// 第3字段cslg标识长沙麓谷机房
			// 第4字段，标识devmode模式，1为devplugin，2为hostpath
			continue
		}
		index, err := strconv.Atoi(fields[1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		mode, err := strconv.Atoi(fields[3])
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		_ = mode

		client := NewK8SCli(v + "/config")

		ClientSet[uint16(index)] = client
	}

	if len(ClientSet) == 0 {
		panic("invalid kuberneters ClientSet")
	}
}

func TestListPods(t *testing.T) {
	InitK8s(WithKubeConfigPatterm("./cluster*"))
	pods := ClientSet[1].ListPods(WithNamespace("server"))

	for _, pod := range pods {
		fmt.Println(pod.Name /*pod.CreationTimestamp,*/, pod.Labels, pod.Namespace, pod.Status.HostIP, pod.Status.PodIP /*,pod.Status.StartTime*/, pod.Status.Phase /*,pod.Status.ContainerStatuses[0].RestartCount,pod.Status.ContainerStatuses[0].Image*/)
	}
}

// go test -v core/k8s/*.go -run TestGetPod -args mobile-542864
func TestGetPod(t *testing.T) {
	InitK8s(WithKubeConfigPatterm("./cluster*"))

	mid := "mobile-678"
	s := os.Args[len(os.Args)-1]
	if len(s) > 0 {
		mid = s
	}

	cnt := 100
	for cnt > 0 {

		fmt.Println("name:", mid)
		fmt.Println("\n********* GetPod *********")
		status, pod, err := ClientSet[1].GetPod(mid)
		if err == nil {
			fmt.Println(status, pod.Name, pod.Status.Phase, pod.Status.Reason, err)
		} else {
			fmt.Printf("status:%v err:%v\n", status, err)
		}

		fmt.Println("\n********* GetPodRT *********")
		status, pod, err = ClientSet[1].GetPodRT(mid)
		if err == nil {
			fmt.Println(status, pod.Name, pod.Status.Phase, pod.Status.Reason, err)
		} else {
			fmt.Printf("status:%v err:%v\n\n", status, err)
		}

		cnt -= 1
		time.Sleep(time.Second)
	}

}

func TestExistedPod(t *testing.T) {
	InitK8s(WithKubeConfigPatterm("./cluster*"))

	time.Sleep(3 * time.Second)

	status, _, _ := ClientSet[1].GetPod("mobile-89023")
	if status == STATUS_RUNNING {
		fmt.Println("found pod")
	} else {
		fmt.Println("not found pod")
	}
}

func TestWatchPod(t *testing.T) {
	InitK8s(WithKubeConfigPatterm("./cluster*"))

	select {}
}

func TestWatchResult(t *testing.T) {
	client := NewK8SCli("./cluster_1_114_2/config")
	label := ""

	status, pod, _ := client.GetPod("name=mobile-386")
	if status != STATUS_RUNNING {
		fmt.Println("found pod")
		return
	}

	for k := range pod.GetLabels() {
		label = k
		break
	}

	watch, err := client.WatchPod(label)
	if err != nil {
		log.Fatal(err.Error())
	}
	go func() {
		for event := range watch.ResultChan() {
			fmt.Printf("Type:%v ", event.Type)
			p, ok := event.Object.(*v1.Pod)
			if !ok {
				log.Fatal("unexpected type")
			}
			fmt.Printf("%v ", p.Name)
			fmt.Println(p.Status.Phase)
			// fmt.Println(p.Status.ContainerStatuses)
		}
	}()

	select {}
}
