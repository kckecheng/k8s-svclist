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

	router.StaticFile("favicon.ico", "static/image/favicon.ico")
	router.Static("/js", "static/js")
	router.Static("/css", "static/css")
	router.Static("/images", "static/img")
	router.LoadHTMLGlob("templates/*")

	indexFunc := func(c *gin.Context) {
		nodeSvc.Lock.Lock()
		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"title":    "Service",
			"services": nodeSvc.Services,
		})
		nodeSvc.Lock.Unlock()
	}

	nodeFunc := func(c *gin.Context) {
		nodeSvc.Lock.Lock()
		c.HTML(http.StatusOK, "nodes.tmpl", gin.H{
			"title": "Node",
			"nodes": nodeSvc.Nodes,
		})
		nodeSvc.Lock.Unlock()
	}

	router.GET("/", indexFunc)
	router.GET("/services", indexFunc)
	router.GET("/nodes", nodeFunc)

	router.Run(":8080")
}
