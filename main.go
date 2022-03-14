// Package swis-api is RESTful API core backend for sakalWeb Information System v5.
package main

import (
	"fmt"
	"swis-api/users"

	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()

	// testing proxy setting
	router.SetTrustedProxies([]string{
		"10.4.5.0/25",
	})

	// root URI
	router.GET("/", func(c *gin.Context){
		c.String(200, "yes papi")
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})

	// users CRUD
	router.GET("/users", users.GetUsers)
	router.GET("/users/:id", users.GetUserByID)
	router.POST("/users", users.PostUser)
	//router.PUT("/users/:id", users.PutUserByID)
	//router.DELETE("/users/:id", users.DeleteUserByID)

	// attach router to http.Server and start it
	router.Run(":8080")
}

