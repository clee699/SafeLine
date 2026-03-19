package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Upstream struct {
	ID   string `json:"id"`  
	Name string `json:"name"`
}

var upstreams []Upstream

// AddUpstream adds a new upstream service
func AddUpstream(c *gin.Context) {
	var newUpstream Upstream
	if err := c.ShouldBindJSON(&newUpstream); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	upstreams = append(upstreams, newUpstream)
	c.JSON(http.StatusCreated, newUpstream)
}

// EditUpstream edits an existing upstream service
func EditUpstream(c *gin.Context) {
	id := c.Param("id")
	for i, u := range upstreams {
		if u.ID == id {
			var updatedUpstream Upstream
			if err := c.ShouldBindJSON(&updatedUpstream); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			upstreams[i] = updatedUpstream
			upstreams[i].ID = id
			c.JSON(http.StatusOK, upstreams[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Upstream not found"})
}

// DeleteUpstream deletes an upstream service
func DeleteUpstream(c *gin.Context) {
	id := c.Param("id")
	for i, u := range upstreams {
		if u.ID == id {
			upstreams = append(upstreams[:i], upstreams[i+1:]...)
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Upstream not found"})
}

// ListUpstreams lists all upstream services
func ListUpstreams(c *gin.Context) {
	c.JSON(http.StatusOK, upstreams)
}

// HealthCheck checks the health of the service
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func SetupRoutes(r *gin.Engine) {
	r.POST("/upstreams", AddUpstream)
	r.PUT("/upstreams/:id", EditUpstream)
	r.DELETE("/upstreams/:id", DeleteUpstream)
	r.GET("/upstreams", ListUpstreams)
	r.GET("/health", HealthCheck)
}