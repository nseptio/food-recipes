package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nseptio/food-recipes/internal/dto"
	"github.com/nseptio/food-recipes/internal/entity"
	"github.com/nseptio/food-recipes/internal/repository"
	"gorm.io/gorm"
)

type RecipeService struct {
	DB               *gorm.DB
	Validate         validator.Validate
	RecipeRepository *repository.RecipeRepository
}

func NewRecipeService(DB *gorm.DB, validate validator.Validate, recipeRepository *repository.RecipeRepository) *RecipeService {
	return &RecipeService{
		DB:               DB,
		Validate:         validate,
		RecipeRepository: recipeRepository,
	}
}

func (s *RecipeService) GetRecipeById(ctx context.Context, id string) (*entity.Recipe, error) {
	var recipe entity.Recipe
	if err := s.RecipeRepository.FindById(ctx, s.DB, &recipe, id); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Recipe not found")
	}

	return &recipe, nil
}

func (s *RecipeService) FilterByIngredients(ctx context.Context, included []string, excluded []string) (*[]int, error) {
	if len(included) == 0 && len(excluded) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "empty field")
	}

	var idList []int
	if err := s.RecipeRepository.FindIdByIngredients(ctx, s.DB, &idList, included, excluded); err != nil {
		panic(err)
	}

	return &idList, nil
}

func (s *RecipeService) GetByName(ctx context.Context, idList []int, name string) (*[]entity.Recipe, error) {
	if len(name) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "empty field")
	}

	var recipes []entity.Recipe
	if err := s.RecipeRepository.GetByName(ctx, s.DB, &recipes, idList, name); err != nil {
		panic(err) // TODO:
	}

	return &recipes, nil
}

func (s *RecipeService) GetAll(ctx context.Context, idList []int) (*[]entity.Recipe, error) {
	var recipes []entity.Recipe
	if len(idList) > 0 {
		if err := s.RecipeRepository.GetAllById(ctx, s.DB, &recipes, idList); err != nil {
			panic(err) // TODO:
		}
	} else {
		if err := s.RecipeRepository.GetAll(ctx, s.DB, &recipes); err != nil {
			panic(err) // TODO:
		}
	}

	return &recipes, nil
}

func (s *RecipeService) GetMostLiked(ctx context.Context) (*[]dto.RecipeLikes, error) {
	var recipesLikes []dto.RecipeLikes
	if err := s.RecipeRepository.GetMostLiked(ctx, s.DB, &recipesLikes); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get most liked recipes")
	}

	return &recipesLikes, nil
}
