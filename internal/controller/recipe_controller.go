package controller

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nseptio/food-recipes/internal/service"
	"strings"
)

type RecipeController struct {
	RecipeService *service.RecipeService
}

func NewRecipeController(recipeService *service.RecipeService) *RecipeController {
	return &RecipeController{RecipeService: recipeService}
}

// FilterByIngredients Middleware
func (c *RecipeController) FilterByIngredients(ctx *fiber.Ctx) error {
	includeQuery := ctx.Query("include")
	excludeQuery := ctx.Query("exclude")
	if includeQuery == "" && excludeQuery == "" {
		return ctx.Next()
	}

	includedIngredients := strings.Split(includeQuery, ",")
	excludedIngredients := strings.Split(excludeQuery, ",")

	idList, err := c.RecipeService.FilterByIngredients(ctx.UserContext(), includedIngredients, excludedIngredients)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to filter by ingredients")
	}

	ctx.Locals("idList", *idList)
	return ctx.Next()
}

func (c *RecipeController) GetById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	recipe, err := c.RecipeService.GetRecipeById(ctx.UserContext(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Recipe not found")
	}

	return ctx.Status(200).JSON(recipe)
}

func (c *RecipeController) GetByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	var idList []int

	idListInterface := ctx.Locals("idList")
	if idListInterface != nil {
		var ok bool
		idList, ok = idListInterface.([]int)
		if !ok {
			return errors.New("idList is not of type []int")
		}
	}

	filterByName, err := c.RecipeService.GetByName(ctx.UserContext(), idList, name)
	if err != nil {
		panic(err)
	}

	return ctx.Status(200).JSON(filterByName)
}

func (c *RecipeController) GetAll(ctx *fiber.Ctx) error {
	var idList []int

	idListInterface := ctx.Locals("idList")
	if idListInterface != nil {
		var ok bool
		idList, ok = idListInterface.([]int)
		if !ok {
			return errors.New("idList is not of type []int")
		}
	}

	recipes, err := c.RecipeService.GetAll(ctx.UserContext(), idList)
	if err != nil {
		panic(err)
	}

	return ctx.Status(200).JSON(recipes)
}

func (c *RecipeController) GetMostLiked(ctx *fiber.Ctx) error {
	recipesLikes, err := c.RecipeService.GetMostLiked(ctx.UserContext())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get most liked recipes")
	}

	fmt.Println(recipesLikes)
	return ctx.Status(200).JSON(recipesLikes)
}
