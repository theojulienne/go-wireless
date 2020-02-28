package api

import "github.com/gin-gonic/gin"

func SetupRoutes(r gin.IRouter) {
	r.GET("/interfaces", listInterfaces)
	r.GET("/interfaces/:iface", notImplemented)
	r.PUT("/interfaces/:iface", notImplemented)

	r.GET("/interfaces/:iface/aps", listAccesPoints)

	r.GET("/interfaces/:iface/networks", listNetworks)
	r.POST("/interfaces/:iface/networks", notImplemented)
	r.PUT("/interfaces/:iface/networks/:id_or_idstr", notImplemented)
	r.GET("/interfaces/:iface/networks/:id_or_idstr", notImplemented)
	r.DELETE("/interfaces/:iface/networks/:id_or_idstr", notImplemented)

}
