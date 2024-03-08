package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nseptio/food-recipes/internal/config"
	"github.com/nseptio/food-recipes/internal/dto"
	"github.com/nseptio/food-recipes/internal/entity"
	"github.com/nseptio/food-recipes/internal/service"
	"time"
)

type UserController struct {
	UserService   *service.UserService
	RecipeService *service.RecipeService
}

func NewUserController(userService *service.UserService, recipeRepository *service.RecipeService) *UserController {
	return &UserController{UserService: userService, RecipeService: recipeRepository}
}

// RequireAuth Middleware
func (c *UserController) RequireAuth(ctx *fiber.Ctx) error {
	viper := config.NewViper()
	hmacSampleSecret := viper.GetString("app.jwt.key")

	tokenString := ctx.Cookies("Authorization", "")

	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
		}

		id := claims["sub"].(float64)
		user, err := c.UserService.FindUserById(ctx.UserContext(), int(id))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		// attach user to context
		ctx.Locals("user", user)

		return ctx.Next()
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(dto.UserRegisterRequest)
	err := ctx.BodyParser(&request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse request body")
	}

	response, err := c.UserService.Register(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	return ctx.Status(201).JSON(response)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(dto.LoginRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse request body")
	}

	response, err := c.UserService.Login(ctx.UserContext(), request)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	cookie := fiber.Cookie{
		Name:     "Authorization",
		Value:    response.Token,
		MaxAge:   3600 * 24,
		Expires:  time.Time{},
		Secure:   false,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	}
	ctx.Cookie(&cookie)

	return ctx.Status(200).JSON(response)
}

func (c *UserController) LikeRecipe(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	idRecipe := ctx.Params("recipeId")
	recipe, err := c.RecipeService.GetRecipeById(ctx.UserContext(), idRecipe)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Recipe not found")
	}

	err = c.UserService.LikeRecipe(ctx.UserContext(), user.(*entity.User), recipe)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to like recipe")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success like recipe",
	})
}

func (c *UserController) UnlikeRecipe(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	idRecipe := ctx.Params("recipeId")
	recipe, err := c.RecipeService.GetRecipeById(ctx.UserContext(), idRecipe)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Recipe not found")
	}

	err = c.UserService.UnlikeRecipe(ctx.UserContext(), user.(*entity.User), recipe)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to unlike recipe")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success unlike recipe",
	})
}

func (c *UserController) GetLikedRecipes(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	recipes, err := c.UserService.GetLikedRecipes(ctx.UserContext(), user.(*entity.User))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get liked recipes")
	}

	return ctx.Status(200).JSON(recipes)
}
