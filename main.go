package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kckecheng/k8s-svclist/query"
)

func main() {
	// counter := 0
	// for {
	//   counter++
	//   fmt.Printf("Round of queries: %d\n", counter)
	//   nodeSvc := query.NewNodeSVC()
	//   pretty.Println(nodeSvc.Namespaces)
	//   pretty.Println(nodeSvc.Nodes)
	//   pretty.Println(nodeSvc.Services)
	//   time.Sleep(5 * time.Minute)
	// }
	nodeSvc := query.NewNodeSVC()

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.tmpl")
	router.GET("/", func(c *gin.Context) {
		nodeSvc.Lock.Lock()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"nodes":    nodeSvc.Nodes,
			"services": nodeSvc.Services,
		})
		nodeSvc.Lock.Unlock()
	})
	router.Run(":8080")
}
