package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes will put handlers to routes on the given router
func SetupRoutes(r gin.IRouter) {
	r.GET("/interfaces", listInterfaces)
	r.GET("/interfaces/:iface", getInterface)
	r.PUT("/interfaces/:iface", notImplemented)

	r.GET("/interfaces/:iface/aps", listAccesPoints)

	r.GET("/interfaces/:iface/networks", listNetworks)
	r.POST("/interfaces/:iface/networks", addNetwork)
	r.PUT("/interfaces/:iface/networks/:id_or_idstr", notImplemented)
	r.GET("/interfaces/:iface/networks/:id_or_idstr", notImplemented)
	r.DELETE("/interfaces/:iface/networks/:id_or_idstr", notImplemented)
}
