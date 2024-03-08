package repository

import (
	"context"
	"github.com/nseptio/food-recipes/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByEmail(ctx context.Context, db *gorm.DB, user *entity.User, email string) error {
	return db.WithContext(ctx).Where("email = ?", email).Take(user).Error
}

func (r *UserRepository) FindByUsername(ctx context.Context, db *gorm.DB, user *entity.User, name string) error {
	return db.WithContext(ctx).Where("username = ?", name).Take(user).Error
}

func (r *UserRepository) CountByUsernameAndEmail(ctx context.Context, db *gorm.DB, username string, email string) (int64, error) {
	var total int64
	err := db.WithContext(ctx).Model(&entity.User{}).
		Select("username", "email").
		Where("username = ?", username).
		Or("email = ?", email).
		Count(&total).Error
	return total, err
}

func (r *UserRepository) LikeRecipe(ctx context.Context, db *gorm.DB, user *entity.User, recipe *entity.Recipe) error {
	return db.WithContext(ctx).Model(user).Association("LikeRecipes").Append(recipe)
}

func (r *UserRepository) UnlikeRecipe(ctx context.Context, db *gorm.DB, user *entity.User, recipe *entity.Recipe) error {
	return db.WithContext(ctx).Model(user).Association("LikeRecipes").Delete(recipe)
}

func (r *UserRepository) GetLikedRecipes(ctx context.Context, db *gorm.DB, user *entity.User) ([]entity.Recipe, error) {
	var recipes []entity.Recipe
	err := db.WithContext(ctx).Model(user).
		Preload("Ingredients").
		Preload("Steps").
		Association("LikeRecipes").
		Find(&recipes)
	return recipes, err
}
