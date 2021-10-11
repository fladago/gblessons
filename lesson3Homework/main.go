// Package classification Recipes API
//
// This is a sample recipes API.
//
//     Schemes: http, https
//     Host: localhost:8080
//     BasePath: /
//     Version: 1.0.0
//     Contact: Vladimir Romashkin<vladimir@mail.com> https://vladimir.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fladago/recipes-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	authHandler    *handlers.AuthHandler
	recipesHandler *handlers.RecipesHandler
)

func init() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://mongoadmin:secret@172.29.0.1:27017/test?authSource=admin"))
	if err != nil {
		log.Println(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database("demo").Collection("recipes")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "172.29.0.1:6379",
		Password: "",
		DB:       0,
	})

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = &handlers.AuthHandler{}
	status := redisClient.Ping()
	fmt.Println(status)
}

func main() {
	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	err := router.Run()
	if err != nil {
		log.Printf("There is an error: %v", err)
	}
}
