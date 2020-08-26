package main

import (
	"fmt"
	"time"

	"github.com/kckecheng/k8s-svclist/query"
	"github.com/kr/pretty"
)

func main() {
	counter := 0
	for {
		counter++
		fmt.Printf("Round of queries: %d\n", counter)
		nodeSvc := query.NewNodeSVC()
		pretty.Println(nodeSvc.Namespaces)
		pretty.Println(nodeSvc.Nodes)
		pretty.Println(nodeSvc.Services)
		time.Sleep(5 * time.Minute)
	}
}
