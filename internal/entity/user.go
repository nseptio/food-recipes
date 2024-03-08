package entity

type User struct {
	ID          int      `gorm:"column:id;primaryKey"`
	Username    string   `gorm:"column:username;unique;not null"`
	Email       string   `gorm:"column:email;unique;not null"`
	Password    string   `gorm:"column:password;not null"`
	FirstName   string   `gorm:"column:first_name;not null"`
	LastName    string   `gorm:"column:last_name"`
	CreatedAt   int64    `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64    `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	LikeRecipes []Recipe `gorm:"many2many:user_like_recipe;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:recipe_id"`
}

func (u *User) TableName() string {
	return "user"
}
