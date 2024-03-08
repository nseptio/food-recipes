package entity

type Ingredient struct {
	Id       int    `gorm:"primaryKey;column:id;autoIncrement" json:"-"`
	RecipeId int    `gorm:"column:recipe_id" json:"-"`
	Name     string `gorm:"column:name" json:"name"`
	Quantity int    `gorm:"column:quantity" json:"quantity"`
	Unit     string `gorm:"column:unit" json:"unit"`
}

func (i *Ingredient) TableName() string {
	return "ingredient"
}
