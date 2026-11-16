package main

import (
	"os"

	"github.com/cbojar/nossl/fetcher"
	"github.com/cbojar/nossl/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	serveAddress := readServeAddress(os.Args[1:])

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.RemoveExtraSlash = true

	fetcher := fetcher.NewFetcher()
	proxyHandler := handlers.NewProxyHandler(fetcher)

	router.GET("/", handlers.RootHandler)
	router.GET("/proxy", proxyHandler.Handle)

	router.Run(serveAddress)
}

func readServeAddress(args []string) string {
	if len(args) < 1 {
		return "0.0.0.0:61212"
	} else {
		return args[0]
	}
}
