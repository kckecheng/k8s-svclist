package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kckecheng/k8s-svclist/query"
)

func initConfig() (string, uint64) {
	var (
		kubeconfig string
		port       uint64
		err        error
	)

	// Get kubeconfig from KUBECONFIG
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		log.Panic(`kubeconfig should be specified by running "export KUBECONFIG=path/to/kubeconfig/file"`)
	}

	if kubeconfig == "" {
		log.Panic("Env var KUBECONFIG should not be empty")
	}

	// Make sure the kubeconfig file is readble
	_, err = ioutil.ReadFile(kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}

	// Get port from LISTEN
	sport, ok := os.LookupEnv("LISTEN")
	if !ok || sport == "" {
		port = 8080
		log.Println(`Listen on port 8080. Execute "export LISTEN=xxxxx" to change the listening port`)
	} else {
		port, err = strconv.ParseUint(sport, 10, 64)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	return kubeconfig, port
}

func main() {
	kubeconfig, port := initConfig()

	nodeSvc := query.NewNodeSVC(kubeconfig)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.StaticFile("favicon.ico", "static/img/favicon.ico")
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

	router.Run(fmt.Sprintf(":%d", port))
}
