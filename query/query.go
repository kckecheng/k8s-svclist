package query

import (
	"fmt"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// NodeInfo node key information
type NodeInfo struct {
	Name string
	// Namespace string
	Addresses []map[string]string
}

// SVCInfo service key information
type SVCInfo struct {
	Name      string
	Namespace string
	Type      string
	Ports     []map[string]string
}

// NodeSVC node and svc information container
type NodeSVC struct {
	Nodes      []NodeInfo
	Services   []SVCInfo
	Namespaces []string
	Lock       sync.Mutex
	client     corev1.CoreV1Interface
}

// NewNodeSVC init
func NewNodeSVC(kubeconfig string) *NodeSVC {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	client := clientset.CoreV1()
	ret := NodeSVC{client: client}

	// Init update
	updateInfo(&ret)
	// Periodical update
	go updateInfo(&ret)

	return &ret
}

// periodicalUpdateInfo update information periodically
func periodicalUpdateInfo(nodeSvc *NodeSVC) {
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ticker.C:
			updateInfo(nodeSvc)
		}
	}
}

// updateInfo update namespaces, nodes, and services information
func updateInfo(nodeSvc *NodeSVC) {
	client := nodeSvc.client

	nodes := listNode(client)
	namespaces := listNamespace(client)
	services := []SVCInfo{}
	for _, ns := range namespaces {
		svcs := listService(client, ns)
		services = append(services, svcs...)
	}

	nodeSvc.Lock.Lock()
	nodeSvc.Namespaces = namespaces
	nodeSvc.Nodes = nodes
	nodeSvc.Services = services
	nodeSvc.Lock.Unlock()
}

// listNamespace list namespaces
func listNamespace(client corev1.CoreV1Interface) []string {
	var ret []string

	namespaces, err := client.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, ns := range namespaces.Items {
		ret = append(ret, ns.Name)
	}
	return ret
}

// listService List services
func listService(client corev1.CoreV1Interface, namespace string) []SVCInfo {
	var ret []SVCInfo

	services, err := client.Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, service := range services.Items {
		svc := SVCInfo{}
		svc.Name = service.Name
		svc.Namespace = service.Namespace

		spec := service.Spec
		svc.Type = string(spec.Type)

		var ports []map[string]string
		for _, port := range spec.Ports {
			p := map[string]string{
				"name":     port.Name,
				"port":     fmt.Sprintf("%d", port.Port),
				"protocol": string(port.Protocol),
				"nodePort": fmt.Sprintf("%d", port.NodePort),
			}
			ports = append(ports, p)
		}
		svc.Ports = ports
		ret = append(ret, svc)
	}

	return ret
}

// listNode list nodes
func listNode(client corev1.CoreV1Interface) []NodeInfo {
	var ret []NodeInfo

	nodes, err := client.Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, node := range nodes.Items {
		n := NodeInfo{
			Name: node.Name,
			// Namespace: node.Namespace,
		}

		var addresses []map[string]string
		for _, address := range node.Status.Addresses {
			addr := map[string]string{}
			addr["type"] = string(address.Type)
			addr["address"] = string(address.Address)
			addresses = append(addresses, addr)
		}
		n.Addresses = addresses

		ret = append(ret, n)
	}

	return ret
}
