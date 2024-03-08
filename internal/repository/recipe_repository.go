package repository

import (
	"context"
	"github.com/nseptio/food-recipes/internal/dto"
	"github.com/nseptio/food-recipes/internal/entity"
	"gorm.io/gorm"
)

type RecipeRepository struct {
	Repository[entity.Recipe]
}

func NewRecipeRepository() *RecipeRepository {
	return &RecipeRepository{}
}

func (r *RecipeRepository) GetAll(ctx context.Context, db *gorm.DB, recipes *[]entity.Recipe) error {
	return db.WithContext(ctx).Model(&entity.Recipe{}).Limit(20).
		Preload("Ingredients").
		Preload("Steps").
		Find(&recipes).Error
}

func (r *RecipeRepository) GetByName(ctx context.Context, db *gorm.DB, recipes *[]entity.Recipe, idList []int, name string) error {
	query := db.WithContext(ctx).Model(&entity.Recipe{}).Limit(20).
		Preload("Ingredients").
		Preload("Steps")

	if len(idList) > 0 {
		query = query.Where("id In (?)", idList)
	}

	query = query.Where("name LIKE ?", "%"+name+"%")

	return query.Find(&recipes).Error
}

func (r *RecipeRepository) GetAllById(ctx context.Context, db *gorm.DB, recipes *[]entity.Recipe, idList []int) error {
	return db.WithContext(ctx).Model(&entity.Recipe{}).
		Preload("Ingredients").
		Preload("Steps").
		Where("id IN (?)", idList).Find(&recipes).Error
}

func (r *RecipeRepository) FindIdByIngredients(ctx context.Context, db *gorm.DB, idList *[]int, included []string, excluded []string) error {
	return db.WithContext(ctx).Model(&entity.Ingredient{}).
		Select("recipe_id").
		Group("recipe_id").
		Having("SUM(CASE WHEN name IN (?) THEN 1 ELSE 0 END) > 0", included).
		Having("SUM(CASE WHEN name IN (?) THEN 1 ELSE 0 END) = 0", excluded).
		Find(&idList).Error
}

func (r *RecipeRepository) GetMostLiked(ctx context.Context, db *gorm.DB, recipesLikes *[]dto.RecipeLikes) error {
	return db.WithContext(ctx).Model(&entity.Recipe{}).
		Joins("JOIN user_like_recipe ULR ON recipe.id = ULR.recipe_id").
		Select("recipe.id, recipe.name, COUNT(ULR.recipe_id) AS like_count").
		Group("recipe.id").
		Order("like_count DESC").
		Limit(20).
		Scan(&recipesLikes).Error
}
