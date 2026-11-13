package main

import (
	"github.com/cbojar/nossl/fetcher"
	"github.com/cbojar/nossl/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
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

	router.Run("0.0.0.0:61212")
}
