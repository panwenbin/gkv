package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	store = make(map[string]string)
	mu    sync.RWMutex
)

func main() {
	r := gin.Default()

	r.GET("/:key", func(c *gin.Context) {
		key := c.Param("key")
		mu.RLock()
		val, ok := store[key]
		mu.RUnlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.String(http.StatusOK, val)
	})

	r.PUT("/:key", func(c *gin.Context) {
		key := c.Param("key")
		val, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		mu.Lock()
		store[key] = string(val)
		mu.Unlock()
		c.Status(http.StatusCreated)
	})

	r.DELETE("/:key", func(c *gin.Context) {
		key := c.Param("key")
		mu.Lock()
		delete(store, key)
		mu.Unlock()
		c.Status(http.StatusNoContent)
	})

	r.Run(":80")
}
