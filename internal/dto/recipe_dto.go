package dto

type RecipeLikes struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	TotalLikes int    `gorm:"column:like_count" json:"totalLikes"`
}
