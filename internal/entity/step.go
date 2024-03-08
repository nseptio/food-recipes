package entity

type Step struct {
	Id          int    `gorm:"primaryKey;column:id;autoIncrement" json:"-"`
	RecipeId    int    `gorm:"recipe_id" json:"-"`
	StepNumber  int    `gorm:"column:step_number" json:"stepNumber"`
	Description string `gorm:"column:description" json:"description"`
}

func (s *Step) TableName() string {
	return "step"
}
