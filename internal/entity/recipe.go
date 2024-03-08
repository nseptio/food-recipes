package entity

type Recipe struct {
	Id           int          `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	Name         string       `gorm:"column:name" json:"name"`
	CookingTime  int          `gorm:"column:cooking_time" json:"cookingTime"`
	Serving      int          `gorm:"column:serving" json:"serving"`
	Ingredients  []Ingredient `gorm:"foreignKey:recipe_id;references:id" json:"ingredients"`
	Steps        []Step       `gorm:"foreignKey:recipe_id;references:id" json:"steps"`
	LikedByUsers []User       `gorm:"many2many:user_like_recipe;foreignKey:id;joinForeignKey:recipe_id;references:id;joinReferences:user_id" json:"-"`
}

func (r *Recipe) TableName() string {
	return "recipe"
}
