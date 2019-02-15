package main

import (
	"log"
	// etcd v3 registry
	_ "github.com/micro/go-plugins/registry/etcdv3"

	"github.com/gin-gonic/gin"

	"github.com/micro/go-web"
)

type Say struct{}

func (s *Say) Anything(c *gin.Context) {
	log.Print("Received Say.Anything API request")
	c.JSON(200, map[string]string{
		"message": "Hi, this is the Greeter API",
	})
}

func mainttt() {
	// Create service
	service := web.NewService(
		web.Name("com.7yes.api.greeter"),
	)

	service.Init()

	// Create RESTful handler (using Gin)
	say := new(Say)
	router := gin.Default()
	router.Any("/greeter", say.Anything)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
