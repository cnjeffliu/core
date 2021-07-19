package k8s

import (
	"fmt"
	"serv/core/logx"

	coreV1 "k8s.io/api/core/v1"

	"path/filepath"
	"strconv"
	"strings"

	k8srun "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerV1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SCli struct {
	Client  *kubernetes.Clientset
	PodList listerV1.PodLister
	DevMode int //设备访问模式，hostpath和devplugin
}

const k8sConfNameLen = 4

var ClientSet map[uint16]K8SCli
var kubePattern string = "./conf.d/kube/cluster*"

type KubeOption func()

func WithKubeConfigPatterm(pattern string) KubeOption {
	return func() {
		kubePattern = pattern
	}
}

func InitK8s(opts ...KubeOption) {
	for _, opt := range opts {
		opt()
	}

	ClientSet = make(map[uint16]K8SCli)
	matches, _ := filepath.Glob(kubePattern)
	var clientSet *kubernetes.Clientset
	for _, v := range matches {
		fmt.Println(v)
		var tmp K8SCli
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
			logx.Error(err)
			continue
		}
		mode, err := strconv.Atoi(fields[3])
		if err != nil {
			logx.Error(err)
			continue
		}
		logx.Debug(fields)
		tmp.DevMode = mode

		config, err := clientcmd.BuildConfigFromFlags("", v+"/config")
		if err != nil {
			panic(err.Error())
		}

		// 创建连接
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		tmp.Client = clientSet
		stopper := make(chan struct{})
		factory := informers.NewSharedInformerFactory(clientSet, 0)
		podInformer := factory.Core().V1().Pods().Informer()
		nodeInformer := factory.Core().V1().Nodes().Informer()
		//defer runtime.HandleCrash()

		// 启动 informer，list & watch
		go factory.Start(stopper)

		// 从 apiserver 同步资源，即 list
		//此处改为异步，网络不通，就无法启动服务。一旦启动后，网络再断开，倒是不影响使用
		go func() {
			if !cache.WaitForCacheSync(stopper, podInformer.HasSynced, nodeInformer.HasSynced) {
				logx.Error("wait for cache sync error")
				k8srun.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
				return
			}
		}()

		nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				node := obj.(*coreV1.Node)
				fmt.Println("delete not implemented", node.Name)
			},
		})
		tmp.PodList = factory.Core().V1().Pods().Lister()
		ClientSet[uint16(index)] = tmp
	}
	logx.Debug(len(ClientSet))
}
