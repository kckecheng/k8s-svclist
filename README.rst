Kubernetes Service List
=========================

We have a large num. of services running on Kubernetes. Most users only need to cosume services running on Kubernetes but has no need to access Kubernetes itself. Our Kubernetes cluster is on primise and there is no plan to add load balance setup - it is tough for end users to remember Kubernetes node addresses and service ports which are not well known.

This simple web based application lists all services (from all namespaces) running on the Kubernetes cluster and refreshes services every half an hour to  find new services.

Usage
------

::

  go build
  export KUBECONFIG=path/to/kubeconfig/file
  export LISTEN=8888
  ./k8s-svclist
