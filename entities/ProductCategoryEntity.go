package entities

import "time"

const TableNameProductCategory = "product_categories"

type ProductCategoryEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string    `gorm:"column:name;size:255;not null" json:"name"`
	Code      string    `gorm:"column:code;size:2;not null" json:"code"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM.
func (*ProductCategoryEntity) TableName() string {
	return TableNameProductCategory
}
