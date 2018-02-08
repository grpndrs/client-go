// some updates for https://rsmitty.github.io/Kubernetes-Events/
// and http://blog.ctaggart.com/2016/09/accessing-kubernetes-api-on-google.html

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	// "google.golang.org/api/option"
	// "google.golang.org/api/transport"
	// "k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/client/cache"
	// "k8s.io/kubernetes/pkg/client/restclient"
	// "k8s.io/kubernetes/pkg/controller/informers"
	// "k8s.io/kubernetes/pkg/util/wait"

	api        "k8s.io/api/core/v1"
	cache      "k8s.io/client-go/tools/cache"
	clientset  "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	informers  "k8s.io/client-go/informers"
	restclient "k8s.io/client-go/rest"
	v1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	wait       "k8s.io/apimachinery/pkg/util/wait"
)

func newKubernetesClient(clstr *container.Cluster) (*clientset.Clientset, error) {
	cert, err := base64.StdEncoding.DecodeString(clstr.MasterAuth.ClientCertificate)
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(clstr.MasterAuth.ClientKey)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(clstr.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}
	config := &restclient.Config{
		Host:            clstr.Endpoint,
		TLSClientConfig: restclient.TLSClientConfig{CertData: cert, KeyData: key, CAData: ca},
		Username:        clstr.MasterAuth.Username,
		Password:        clstr.MasterAuth.Password,
		// Insecure:        true,
	}
	kbc, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kbc, nil
}

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
	si.Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    podCreated,
			DeleteFunc: podDeleted,
		},
	)
	si.Start(wait.NeverStop)
}
