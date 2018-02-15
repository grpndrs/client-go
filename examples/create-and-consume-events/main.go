// some updates for https://rsmitty.github.io/Kubernetes-Events/
// and http://blog.ctaggart.com/2016/09/accessing-kubernetes-api-on-google.html
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	api        "k8s.io/api/core/v1"
	cache      "k8s.io/client-go/tools/cache"
	clientset  "k8s.io/client-go/kubernetes"
	informers  "k8s.io/client-go/informers"
	restclient "k8s.io/client-go/rest"
	wait       "k8s.io/apimachinery/pkg/util/wait"
)

func podCreated(obj interface{}) {
	pod := obj.(*api.Pod)
	fmt.Printf("Pod created: %s\n", pod.ObjectMeta.Name)
}

func podDeleted(obj interface{}) {
	pod := obj.(*api.Pod)
	fmt.Printf("Pod deleted: %s\n", pod.ObjectMeta.Name)
}

func watchPods(kbc *clientset.Clientset) {
	resyncPeriod := 30 * time.Minute
	si := informers.NewSharedInformerFactory(kbc, resyncPeriod)
	si.Core().V1().Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    podCreated,
			DeleteFunc: podDeleted,
		},
	)
	si.Start(wait.NeverStop)
}

func main() {

	//Configure cluster info
	config := &restclient.Config{
		Host:     "https://xxx.yyy.zzz:443",
		Username: "kube",
		Password: "supersecretpw",
	}

	//Create a new client to interact with cluster and freak if it doesn't work
	kbc, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalln("Client not created sucessfully:", err)
	}

	//Watch for Pods
	watchPods(kbc)

	//Keep alive
	log.Fatal(http.ListenAndServe(":8080", nil))
}
