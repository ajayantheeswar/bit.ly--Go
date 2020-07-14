package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ajayantheeswar/bit.ly/database"
	"github.com/ajayantheeswar/bit.ly/AuthHandler"
	"github.com/ajayantheeswar/bit.ly/LinkHandler"
	
)

func main() {
	database.ConnectDatabase()
	
	router := gin.Default()
	//router.Use(cors.Default())/
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"PUT", "PATCH","POST","GET","DELETE"},
        AllowHeaders:     []string{"Origin","Content-Length","Content-Type","Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	router.GET("/:url",linkHandler.RedirectLink)

	router.POST("/signup",authHandler.AuthSignup)
	router.POST("/signin",authHandler.AuthSignIn)

	router.POST("/createlink",authHandler.AuthMiddleware(),linkHandler.CreateLink)
	router.POST("/getalllinks",authHandler.AuthMiddleware(),linkHandler.GetAllLinks)

	router.Run()
	
}