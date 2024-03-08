package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/nseptio/food-recipes/internal/config"
	"github.com/nseptio/food-recipes/internal/controller"
	"github.com/nseptio/food-recipes/internal/repository"
	"github.com/nseptio/food-recipes/internal/service"
)

func main() {
	viper := config.NewViper()
	app := config.NewFiber(viper)

	recipeRepository := repository.NewRecipeRepository()
	recipeService := service.NewRecipeService(
		config.ConnectToDatabase(viper),
		*validator.New(),
		recipeRepository,
	)

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(
		config.ConnectToDatabase(viper),
		userRepository,
		*validator.New(),
	)

	recipeController := controller.NewRecipeController(recipeService)
	userController := controller.NewUserController(userService, recipeService)

	recipe := app.Group("recipes")
	user := app.Group("users")

	recipe.Use(userController.RequireAuth, recipeController.FilterByIngredients)

	recipe.Get("/", recipeController.GetAll)
	recipe.Get("/name/:name", recipeController.GetByName)
	recipe.Get("/id/:id", recipeController.GetById)
	recipe.Get("/top-liked", recipeController.GetMostLiked)

	user.Post("/register", userController.Register)
	user.Post("/login", userController.Login)

	userToRecipe := user.Group("/recipes")
	userToRecipe.Use(userController.RequireAuth)
	userToRecipe.Post("/:recipeId/like", userController.LikeRecipe)
	userToRecipe.Delete("/:recipeId/like", userController.UnlikeRecipe)
	userToRecipe.Get("/liked", userController.GetLikedRecipes)

	err := app.Listen("localhost:8080")
	if err != nil {
		panic(err)
	}
}
