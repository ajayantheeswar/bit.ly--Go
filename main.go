package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ajayantheeswar/bit.ly/controllers"
)

func main() {
	fmt.Print("hii")
	
	router := gin.Default()

	router.Use(cors.Default())
	router.Run()

	controllers.ConnectDatabase()
}