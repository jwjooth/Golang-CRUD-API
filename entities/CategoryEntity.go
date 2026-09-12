package entities

import "time"

const TableNameCategory = "categories"

type CategoryEntity struct {
	ID uint `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name string `gorm:"column:name;size:255;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (*CategoryEntity) TableName() string {
	return TableNameCategory
}
