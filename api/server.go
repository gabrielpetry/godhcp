package api

import (
	DhcpdumpMessage "go-dhcpdump/dhcpdump"
	"go-dhcpdump/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartApiServer() {
	log.Info("Starting API server")
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/devices", func(c *gin.Context) {
		device := DhcpdumpMessage.DhcpdumpMessage{}
		c.JSON(http.StatusOK, gin.H{
			"devices": device.GetAll(),
		})
	})

	r.Run()
}
