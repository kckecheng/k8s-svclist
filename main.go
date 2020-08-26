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

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static")

	indexFunc := func(c *gin.Context) {
		nodeSvc.Lock.Lock()
		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"services": nodeSvc.Services,
		})
		nodeSvc.Lock.Unlock()
	}

	nodeFunc := func(c *gin.Context) {
		nodeSvc.Lock.Lock()
		c.HTML(http.StatusOK, "nodes.tmpl", gin.H{
			"nodes": nodeSvc.Nodes,
		})
		nodeSvc.Lock.Unlock()
	}

	router.GET("/", indexFunc)
	router.GET("/services", indexFunc)
	router.GET("/nodes", nodeFunc)

	router.Run(":8080")
}
