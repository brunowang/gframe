package main

import (
	"github.com/brunowang/gframe/gfk8s/gfinformer/example/handler"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	listWatcher := cache.NewListWatchFromClient(client.CoreV1().RESTClient(),
		"configmaps", "default", fields.Everything())
	_, informer := cache.NewInformer(listWatcher, &v1.ConfigMap{}, 0, &handler.CmdHandler{})
	informer.Run(wait.NeverStop)

	select {}
}
