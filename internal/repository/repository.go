package repository

import (
	"context"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Select("id").Where("id = ?", id).Take(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(ctx context.Context, db *gorm.DB, entity *T, id any) error {
	return db.WithContext(ctx).Where("id = ?", id).Take(entity).Error
}
